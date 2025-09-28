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

// Client to read data from the iCon device.
type IconClient interface {
	io.Closer
	// Returns the system ID.
	SysId() string
	// Reads data from the device.
	ReadValues() (*model.DataPollResponse, error)
	// Reads data from the device.
	ReadValuesContext(ctx context.Context) (*model.DataPollResponse, error)
	// Experimental!
	// Sets the thermostat settings
	SetThermostatSettings(tab int, thermosSettings model.ThermostatSettings) error
	// Experimental!
	// Sets the thermostat settings
	SetThermostatSettingsContext(ctx context.Context, tab int, thermosSettings model.ThermostatSettings) error
	// Experimental!
	// Sets the general settings
	SetGeneralSettings(tab int, generalSettings *model.GeneralSettings) error
	// Experimental!
	// Sets the general settings
	SetGeneralSettingsContext(ctx context.Context, tab int, generalSettings *model.GeneralSettings) error
}

type iconHttpClient struct {
	authClient   IconAuthClient
	readerClient IconReaderClient
	writerClient IconWriterClient
	sysId        string
	password     string
	session      *IconSession
}

// Creates a new client to fetch data from an iCON device.
func NewIconClient(urlStr string, sysId string, password string) (IconClient, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 1 * time.Second,
			}).Dial,
		},
		Timeout: 10 * time.Second,
	}
	return &iconHttpClient{
		authClient:   NewAuthClient(client, u),
		readerClient: NewIconReaderClient(client, u),
		writerClient: NewIconWriterClient(client, u),
		sysId:        sysId,
		password:     password,
		session:      nil,
	}, nil
}

// Returns the system ID.
func (client *iconHttpClient) SysId() string {
	return client.sysId
}

// Cleans up the client.
func (client *iconHttpClient) Close() error {
	if client.session == nil {
		return nil
	}
	return client.authClient.Logout(context.Background(), client.session)
}

// Reads data from the device.
func (client *iconHttpClient) ReadValues() (*model.DataPollResponse, error) {
	return client.ReadValuesContext(context.Background())
}

// Reads data from the device.
func (client *iconHttpClient) ReadValuesContext(ctx context.Context) (*model.DataPollResponse, error) {
	err := client.ensureSession(ctx)
	if err != nil {
		return nil, err
	}
	return client.readerClient.ReadValues(ctx, client.session)
}

func (client *iconHttpClient) SetThermostatSettings(tab int, thermosSettings model.ThermostatSettings) error {
	return client.SetThermostatSettingsContext(context.Background(), tab, thermosSettings)
}

func (client *iconHttpClient) SetThermostatSettingsContext(ctx context.Context, tab int, thermosSettings model.ThermostatSettings) error {
	err := client.ensureSession(ctx)
	if err != nil {
		return err
	}
	return client.writerClient.SetThermostatSettings(ctx, tab, thermosSettings, client.session)
}

func (client *iconHttpClient) SetGeneralSettings(tab int, generalSettings *model.GeneralSettings) error {
	return client.SetGeneralSettingsContext(context.Background(), tab, generalSettings)
}

func (client *iconHttpClient) SetGeneralSettingsContext(ctx context.Context, tab int, generalSettings *model.GeneralSettings) error {
	err := client.ensureSession(ctx)
	if err != nil {
		return err
	}
	return client.writerClient.SetGeneralSettings(ctx, tab, generalSettings, client.session)
}

func (client *iconHttpClient) ensureSession(ctx context.Context) error {
	if client.session == nil || !client.session.Valid() {
		session, err := client.authClient.Login(ctx, client.sysId, client.password)
		if err != nil {
			return err
		}
		client.session = session
	}
	return nil
}
