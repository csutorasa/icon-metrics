// iCON client scraper
package client

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/csutorasa/icon-metrics/pkg/model"
)

// Client to manage an iCon device.
type IconClient interface {
	io.Closer
	// Returns the system ID.
	SysId() string
	// Reads data from the device.
	ReadValues() (*model.DataPollResponse, error)
	// Reads data from the device.
	ReadValuesContext(ctx context.Context) (*model.DataPollResponse, error)
	// Experimental!
	// Sets the thermostat settings.
	SetThermostatSettings(tab int, thermosSettings model.ThermostatSettings) error
	// Experimental!
	// Sets the thermostat settings.
	SetThermostatSettingsContext(ctx context.Context, tab int, thermosSettings model.ThermostatSettings) error
	// Experimental!
	// Sets the general settings.
	SetGeneralSettings(tab int, generalSettings *model.GeneralSettings) error
	// Experimental!
	// Sets the general settings.
	SetGeneralSettingsContext(ctx context.Context, tab int, generalSettings *model.GeneralSettings) error
}

// Creates a new client to manage an iCON device.
func NewIconClient(url *url.URL, sysId string, password string) (IconClient, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 1 * time.Second,
			}).Dial,
		},
		Timeout: 10 * time.Second,
	}
	return &IconHttpClient{
		authClient:   NewAuthClient(client, url),
		readerClient: NewIconReaderClient(client, url),
		writerClient: NewIconWriterClient(client, url),
		sysId:        sysId,
		password:     password,
		session:      nil,
	}, nil
}

// HTTP Client to manage an iCon device.
type IconHttpClient struct {
	authClient   IconAuthClient
	readerClient IconReaderClient
	writerClient IconWriterClient
	sysId        string
	password     string
	session      *IconSession
}

func (client *IconHttpClient) SysId() string {
	return client.sysId
}

// Cleans up the client.
func (client *IconHttpClient) Close() error {
	if client.session == nil {
		return nil
	}
	return client.authClient.Logout(context.Background(), client.session)
}

func (client *IconHttpClient) ReadValues() (*model.DataPollResponse, error) {
	return client.ReadValuesContext(context.Background())
}

func (client *IconHttpClient) ReadValuesContext(ctx context.Context) (*model.DataPollResponse, error) {
	err := client.ensureSession(ctx)
	if err != nil {
		return nil, err
	}
	return client.readerClient.ReadValues(ctx, client.session)
}

func (client *IconHttpClient) SetThermostatSettings(tab int, thermosSettings model.ThermostatSettings) error {
	return client.SetThermostatSettingsContext(context.Background(), tab, thermosSettings)
}

func (client *IconHttpClient) SetThermostatSettingsContext(ctx context.Context, tab int, thermosSettings model.ThermostatSettings) error {
	err := client.ensureSession(ctx)
	if err != nil {
		return err
	}
	return client.writerClient.SetThermostatSettings(ctx, tab, thermosSettings, client.session)
}

func (client *IconHttpClient) SetGeneralSettings(tab int, generalSettings *model.GeneralSettings) error {
	return client.SetGeneralSettingsContext(context.Background(), tab, generalSettings)
}

func (client *IconHttpClient) SetGeneralSettingsContext(ctx context.Context, tab int, generalSettings *model.GeneralSettings) error {
	err := client.ensureSession(ctx)
	if err != nil {
		return err
	}
	return client.writerClient.SetGeneralSettings(ctx, tab, generalSettings, client.session)
}

func (client *IconHttpClient) ensureSession(ctx context.Context) error {
	if client.session == nil || !client.session.Valid() {
		session, err := client.authClient.Login(ctx, client.sysId, client.password)
		if err != nil {
			return err
		}
		client.session = session
	}
	return nil
}
