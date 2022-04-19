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
		ItemsPerRequest:   10,
		LogSuccessfulPush: true,
		LogErrorPush:      true,
		LogProgress:       false,
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
	fileR, err := os.Open("files/example.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = fileR.Close()
	}()
	data, err := file.ReadAll(fileR, true)
	if err != nil {
		log.Fatal(err)
	}

	//Create file to Write
	fileW, err := os.Create("files/errors.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = fileW.Close()
	}()

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

//Multi mode
func publishMultiMode(ctx context.Context, data [][]string, fileW io.Writer, repo *repository.Repository) {
	var processedLines, errorCounter int
	toPublish := make([][]string, 0, republishConfig.ItemsPerRequest)
	var lastRecords bool
	var total = len(data)
	for {
		if total-processedLines < republishConfig.ItemsPerRequest {
			toPublish = data[processedLines:]
			lastRecords = true
		} else {
			toPublish = data[processedLines : processedLines+republishConfig.ItemsPerRequest]
		}
		if republishConfig.LogProgress {
			log.Printf("Records processed: %v of %v", processedLines, total)
		}
		//MultiPublish items
		response, err := repo.MultiPublish(ctx, toPublish)
		if err != nil {
			if republishConfig.LogErrorPush {
				log.Printf("Error publishing: %s", err.Error())
			}
			errorCounter += len(toPublish)
			err = file.WriteAll(fileW, toPublish)
			if err != nil {
				log.Printf("Error writing to file: %s", err.Error())
				break
			}
			processedLines += len(toPublish)
			if lastRecords {
				break
			}
			continue
		}
		errorCounter += len(response.Errors)
		for _, item := range response.Errors {
			err = file.Write(fileW, item)
			if err != nil {
				log.Printf("Error writing to file: %s", err.Error())
				break
			}
		}
		if republishConfig.LogSuccessfulPush {
			for _, item := range response.Success {
				log.Printf("resource with id:%v processed", item)
			}
		}
		processedLines += len(toPublish)
		if lastRecords {
			break
		}
	}
	log.Printf("Records processed: %v", processedLines)
	log.Printf("Records with error: %v", errorCounter)
}
