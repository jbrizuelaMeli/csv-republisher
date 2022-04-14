package main

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mercadolibre/csv-republisher/config"
	"github.com/mercadolibre/csv-republisher/repository"
	"github.com/mercadolibre/csv-republisher/tools/customcontext"
	"github.com/mercadolibre/csv-republisher/tools/file"
	"github.com/mercadolibre/csv-republisher/tools/restclient"
)

var (
	republishConfig = config.RepublishConfig{
		ItemsPerRequest:  1,
		Goroutines:       6,
		RequestPerSecond: 60,
	}
	restClientConfig = restclient.Config{
		TimeoutMillis: 3000,
		ApiDomain:     "https://production-writer-republish_account-cashbacks-api.furyapps.io",
		ExternalApiCalls: map[string]restclient.ExternalApiCall{
			"cashback-api": {
				Resources: map[string]restclient.Resource{
					"cashback-republish": {
						RequestUri: "/cashback/republish",
					},
				},
			},
		},
	}
)

func main() {
	// open file to Read
	fileR, err := os.Open("files/cashbacks-for-republish.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fileR.Close()
	data, err := file.ReadAll(fileR, true)
	if err != nil {
		log.Fatal(err)
	}

	//Create file to Write
	fileW, err := os.Create("files/errors.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fileW.Close()

	//Build repository
	rc, err := restclient.NewRestClient(restClientConfig)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewRepository(rc)

	//Build context
	ctx := customcontext.WithoutCancel(context.Background())

	if republishConfig.ItemsPerRequest == 1 {
		//Single mode
		publishSingleMode(ctx, data, fileW, repo)
	} else {
		//MultiMode
		publishMultiMode(ctx, data, fileW, repo)
	}

	return
}

func publishSingleMode(ctx context.Context, data [][]string, fileW io.Writer, repo *repository.Repository) {
	// limit concurrency
	semaphore := make(chan struct{}, republishConfig.Goroutines)

	// setting a max rate in req/sec
	rate := make(chan struct{}, republishConfig.RequestPerSecond)
	for i := 0; i < cap(rate); i++ {
		rate <- struct{}{}
	}

	// leaky bucket
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			_, ok := <-rate
			// if this isn't going to run indefinitely, signal
			// this to return by closing the rate channel.
			if !ok {
				return
			}
		}
	}()

	//Run
	var wg sync.WaitGroup
	for _, line := range data {
		wg.Add(1)
		go func(l []string) {
			defer wg.Done()

			// wait for the rate limiter
			rate <- struct{}{}

			// check the concurrency semaphore
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()

			err := repo.Publish(ctx, l)
			if err != nil {
				log.Println(err.Error())
				file.Write(fileW, l)
				return
			}
			log.Printf("resource with ids:%v processed", l)
		}(line)
	}
	wg.Wait()
	close(rate)
}

func publishMultiMode(ctx context.Context, data [][]string, fileW io.Writer, repo *repository.Repository) {
	//Multi mode
	toPublish := make([][]string, 0)
	for idx, line := range data {
		toPublish = append(toPublish, line)
		if idx%(republishConfig.ItemsPerRequest-1) == 0 && idx > 0 {
			//MultiPublish array
			err := repo.MultiPublish(ctx, toPublish)
			if err != nil {
				log.Println(err.Error())
				_ = file.WriteAll(fileW, toPublish)
				continue
			}
			log.Printf("resource with ids:%v processed", toPublish)
		}
	}
}
