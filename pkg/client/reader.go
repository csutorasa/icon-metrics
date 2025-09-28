package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/csutorasa/icon-metrics/pkg/model"
)

// Client to read data from an iCON device.
type IconReaderClient interface {
	// Reads data from the device.
	ReadValues(ctx context.Context, session *IconSession) (*model.DataPollResponse, error)
}

type IconReaderHttpClient struct {
	client *http.Client
	url    *url.URL
}

func NewIconReaderClient(c *http.Client, baseUrl *url.URL) IconReaderClient {
	return &IconReaderHttpClient{
		client: c,
		url:    baseUrl,
	}
}

// Reads data from the device.
func (client *IconReaderHttpClient) ReadValues(ctx context.Context, session *IconSession) (*model.DataPollResponse, error) {
	req, err := postFormRequestContextWithSession(ctx, getUrlWithPath(client.url, "index.php"), readValuesRequest(), session)
	if err != nil {
		return nil, err
	}
	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	updateSession(res, session)
	data := &model.DataPollResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func readValuesRequest() url.Values {
	return url.Values{
		"tab": []string{"datapoll"},
	}
}
