package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/csv-republisher/config"
	"github.com/csv-republisher/repository"
	"github.com/csv-republisher/tools/customcontext"
	"github.com/csv-republisher/tools/file"
	"github.com/csv-republisher/tools/restclient"
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

	//MultiMode
	publishMultiMode(ctx, data, fileW, repo)

	return
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
