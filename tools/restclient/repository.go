package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/csv-republisher/tools/apierrors"
)

type RestClient interface {
	BuildUrl(externalApi string, resource string, params ...interface{}) (string, error)
	HandleError(ctx context.Context, err error, res *http.Response) error
	DoGet(ctx context.Context, url string, result interface{}, additionalHeaders ...Header) error
	DoPost(ctx context.Context, url string, body interface{}, result interface{}, additionalHeaders ...Header) error
}

type restClient struct {
	config Config
	client http.Client
}

type Header struct {
	Key   string
	Value string
}

func NewRestClient(config Config) (RestClient, error) {
	client := http.Client{
		Timeout: config.TimeoutMillis * time.Millisecond,
	}
	return &restClient{
		config: config,
		client: client,
	}, nil
}

func (rc restClient) BuildUrl(externalApi string, resource string, params ...interface{}) (string, error) {
	url := rc.config.ApiDomain
	if val, exist := rc.config.ExternalApiCalls[externalApi]; exist {
		if url == "" {
			url = val.ApiDomain
		}
		url += fmt.Sprintf(val.Resources[resource].RequestUri, params...)
		return url, nil
	}
	return url, errors.New("resource_not_found")
}

func (rc restClient) HandleError(ctx context.Context, err error, res *http.Response) error {
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return apierrors.NewCommunicationError(http.StatusText(res.StatusCode), res.StatusCode)
	}
	return nil
}

func (rc restClient) DoGet(ctx context.Context, url string, result interface{}, additionalHeaders ...Header) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	for _, header := range additionalHeaders {
		req.Header.Add(header.Key, header.Value)
	}
	res, err := rc.client.Do(req)
	if err := rc.HandleError(ctx, err, res); err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = res.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	return nil
}

func (rc restClient) DoPost(ctx context.Context, url string, body interface{}, result interface{}, additionalHeaders ...Header) error {
	var requestBody io.Reader
	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		requestBody = bytes.NewBuffer(reqBody)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, requestBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	for _, header := range additionalHeaders {
		req.Header.Add(header.Key, header.Value)
	}
	res, err := rc.client.Do(req)
	if err := rc.HandleError(ctx, err, res); err != nil {
		return err
	}
	bodyResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err // unable to read whole response body, possible request timeout or TCP error.
	}
	err = res.Body.Close()
	if err != nil {
		return err
	}
	if result != nil {
		if err := json.Unmarshal(bodyResponse, result); err != nil {
			return err
		}
	}
	return nil
}
