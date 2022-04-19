package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/csv-republisher/model"
	"github.com/csv-republisher/tools/restclient"
)

const FuryToken = "{token}"

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

func (r Repository) MultiPublish(ctx context.Context, lines [][]string) (*model.ItemResponseMap, error) {
	url, err := r.restClient.BuildUrl("cashback-api", "cashback-republish")
	if err != nil {
		return nil, err
	}

	itemReqMap := &model.ItemRequestMap{
		Item: make(map[int64][]string),
	}
	bulkResponseMap := &model.ItemResponseMap{
		Success: make(map[int64][]string),
		Errors:  make(map[int64][]string),
	}

	request := &model.MultiRequestNumericIDs{}
	for _, line := range lines {
		id, err := strconv.ParseInt(strings.Join(line, ""), 10, 64)
		if err != nil {
			return nil, err
		}
		item := model.NumericID{ID: id}
		request.IDs = append(request.IDs, item)
		itemReqMap.Item[id] = line
	}

	response := &model.MultiResponseNumericIDs{}
	err = r.restClient.DoPost(ctx, url, request, response, restclient.Header{Key: "X-Auth-Token", Value: FuryToken})
	if err != nil {
		return nil, err
	}

	for _, item := range response.Errors {
		if _, exist := itemReqMap.Item[item.ID]; exist {
			bulkResponseMap.Errors[item.ID] = itemReqMap.Item[item.ID]
		}
	}

	for _, item := range response.IDs {
		if _, exist := itemReqMap.Item[item.ID]; exist {
			bulkResponseMap.Success[item.ID] = itemReqMap.Item[item.ID]
		}
	}

	return bulkResponseMap, nil
}
