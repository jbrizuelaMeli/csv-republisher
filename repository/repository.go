package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/mercadolibre/csv-republisher/model"
	"github.com/mercadolibre/csv-republisher/tools/restclient"
)

const FuryToken = "3712946dfa4ddb21c7ebd4a5b0fa64e266e8e899a45c2ae16f30d7444681541a"

type Repository struct {
	restClient restclient.RestClient
}

func NewRepository(client restclient.RestClient) *Repository {
	return &Repository{
		restClient: client,
	}
}

func (r Repository) Publish(ctx context.Context, line []string) error {
	url, err := r.restClient.BuildUrl("cashback-api", "cashback-republish")
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(strings.Join(line, ""), 10, 64)
	if err != nil {
		return err
	}
	request := &model.NumericID{ID: id}

	err = r.restClient.DoPost(ctx, url, request, nil, restclient.Header{Key: "X-Auth-Token", Value: FuryToken})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) MultiPublish(ctx context.Context, lines [][]string) error {
	url, err := r.restClient.BuildUrl("cashback-api", "cashback-republish")
	if err != nil {
		return err
	}

	request := &model.NumericIDs{}
	for _, line := range lines {
		id, err := strconv.ParseInt(strings.Join(line, ""), 10, 64)
		if err != nil {
			return err
		}
		request.IDs = append(request.IDs, id)
	}

	err = r.restClient.DoPost(ctx, url, request, nil, restclient.Header{Key: "X-Auth-Token", Value: FuryToken})
	if err != nil {
		return err
	}
	return nil
}
