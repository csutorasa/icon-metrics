// iCON client scraper
package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/csutorasa/icon-metrics/pkg/model"
)

// Experimental!
// Client to write settings to an iCON device.
type IconWriterClient interface {
	// Experimental!
	// Sets the thermostat settings.
	SetThermostatSettings(ctx context.Context, tab int, thermosSettings model.ThermostatSettings, session *IconSession) error
	// Experimental!
	// Sets the general settings.
	SetGeneralSettings(ctx context.Context, tab int, generalSettings *model.GeneralSettings, session *IconSession) error
}

// Experimental!
// HTTP client to write settings to an iCON device.
type IconWriterHttpClient struct {
	client *http.Client
	url    *url.URL
}

// Experimental!
// Creates a new IconWriterClient.
func NewIconWriterClient(c *http.Client, baseUrl *url.URL) IconWriterClient {
	return &IconWriterHttpClient{
		client: c,
		url:    baseUrl,
	}
}

func (client *IconWriterHttpClient) SetThermostatSettings(ctx context.Context, tab int, thermosSettings model.ThermostatSettings, session *IconSession) error {
	formData := getValues(thermosSettings.ToValues(tab))
	req, err := postFormRequestContextWithSession(ctx, getUrlWithPath(client.url, "index.php"), formData, session)
	if err != nil {
		return err
	}
	res, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return validateAction(res)
}

func (client *IconWriterHttpClient) SetGeneralSettings(ctx context.Context, tab int, generalSettings *model.GeneralSettings, session *IconSession) error {
	formData := getValues(generalSettings.ToValues(tab))
	req, err := postFormRequestContextWithSession(ctx, getUrlWithPath(client.url, "index.php"), formData, session)
	if err != nil {
		return err
	}
	res, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return validateAction(res)
}

func getValues(parameters map[string][]string) url.Values {
	query := url.Values{}
	for key, values := range parameters {
		if len(values) == 0 {
			continue
		}
		if len(values) == 1 {
			query.Set(key, values[0])
		} else {
			for _, value := range values {
				query.Add(key, value)
			}
		}
	}
	return query
}
