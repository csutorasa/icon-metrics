// iCON client scraper
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/csutorasa/icon-metrics/metrics"
	"github.com/csutorasa/icon-metrics/model"
)

// Client to read data from the iCon device.
type IconClient interface {
	HttpSessionClient
	IconClientReader
	IconClientWriter
	// Returns the system ID.
	SysId() string
}

type iconHttpClient struct {
	client    *http.Client
	url       *url.URL
	sysId     string
	password  string
	sessionId string
	session   metrics.MetricsSession
}

// session cookie name
const phpSessionId = "PHPSESSID"

// Default body size is <1kB and each DP is <1kB.
// Each error should be <1kB.
// 1 MB should be a safe hard limit.
const maxReadBytes = 1024 * 1024

// Creates a new client to fetch data from an iCON device.
func NewIconClient(urlStr string, sysId string, password string, session metrics.MetricsSession) (IconClient, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return &iconHttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 1 * time.Second,
				}).Dial,
			},
			Timeout: 10 * time.Second,
		},
		url:       u,
		sysId:     sysId,
		password:  password,
		sessionId: "",
		session:   session,
	}, nil
}

// Returns the system ID.
func (client *iconHttpClient) SysId() string {
	return client.sysId
}

// Join path to the base URL.
func (client *iconHttpClient) getPath(p string) (*url.URL, error) {
	u, err := url.Parse(client.url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create url from %s: %w", client.url.String(), err)
	}
	u.Path = path.Join(u.Path, p)
	return u, nil
}

// Removes metrics and session data.
func (client *iconHttpClient) removeSession() {
	client.sessionId = ""
}

// Unmarshal JSON content from http response body.
func unmarshalBody(res *http.Response, v any) error {
	body, err := io.ReadAll(io.LimitReader(res.Body, maxReadBytes))
	if err != nil {
		return err
	}
	if len(body) == maxReadBytes {
		return fmt.Errorf("too long response body")
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

// Updates session data from cookies.
func (client *iconHttpClient) updateCookie(cookies []*http.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == phpSessionId {
			client.sessionId = cookie.Value
			return true
		}
	}
	return false
}

// Manages PHP session to an iCON device.
type HttpSessionClient interface {
	io.Closer
	// Logs in and creates a session.
	Login() error
	// Closes a session.
	Logout() error
	// Returns if there is a session.
	IsLoggedIn() bool
}

// Logs in and creates a session.
func (client *iconHttpClient) Login() error {
	timer := metrics.NewTimer()
	formData := url.Values{
		"sysid":    []string{client.sysId},
		"password": []string{client.password},
		"lang":     []string{"hu"},
		"tab":      []string{"login"},
		"form":     []string{"login"},
	}
	req, err := http.NewRequest(http.MethodPost, client.url.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		client.session.HttpClientRequest("login", 0, timer.End())
		return fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		client.session.HttpClientRequest("login", 0, timer.End())
		return fmt.Errorf("failed to execute http call: %w", err)
	}
	defer res.Body.Close()
	client.session.HttpClientRequest("login", res.StatusCode, timer.End())
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login, status code %d", res.StatusCode)
	}
	if !client.updateCookie(res.Cookies()) {
		return errors.New("no session was found")
	}
	data := model.ActionResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		client.removeSession()
		return fmt.Errorf("failed to parse json: %w", err)
	}
	if !data.IsSuccess() {
		client.removeSession()
		return data.CreateError()
	}
	return nil
}

// Closes a session.
func (client *iconHttpClient) Logout() error {
	timer := metrics.NewTimer()
	fomrData := url.Values{
		"logout": []string{"true"},
	}
	url, err := client.getPath("index.php")
	if err != nil {
		client.session.HttpClientRequest("logout", 0, timer.End())
		return err
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		client.session.HttpClientRequest("logout", 0, timer.End())
		return fmt.Errorf("failed to create request: %s", err)
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		client.session.HttpClientRequest("logout", 0, timer.End())
		return fmt.Errorf("failed to execute http call: %w", err)
	}
	defer res.Body.Close()
	client.session.HttpClientRequest("logout", res.StatusCode, timer.End())
	client.removeSession()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to logout, status code %d", res.StatusCode)
	}
	return nil
}

// Returns if there is a session.
func (client *iconHttpClient) IsLoggedIn() bool {
	return client.sessionId != ""
}

// Cleans up the client.
func (client *iconHttpClient) Close() error {
	err := client.Logout()
	return err
}

// Client to read data from an iCON device.
type IconClientReader interface {
	// Reads data from the device.
	ReadValues() (*model.DataPollResponse, error)
}

// Reads data from the device.
func (client *iconHttpClient) ReadValues() (*model.DataPollResponse, error) {
	timer := metrics.NewTimer()
	fomrData := url.Values{
		"tab": []string{"datapoll"},
	}
	url, err := client.getPath("index.php")
	if err != nil {
		client.session.HttpClientRequest("read_values", 0, timer.End())
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		client.session.HttpClientRequest("read_values", 0, timer.End())
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		client.session.HttpClientRequest("read_values", 0, timer.End())
		client.removeSession()
		return nil, fmt.Errorf("failed to execute http call: %w", err)
	}
	defer res.Body.Close()
	client.session.HttpClientRequest("read_values", res.StatusCode, timer.End())
	if res.StatusCode != http.StatusOK {
		client.removeSession()
		return nil, fmt.Errorf("failed to read data, status code %d", res.StatusCode)
	}
	client.updateCookie(res.Cookies())
	data := &model.DataPollResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		client.removeSession()
		return data, fmt.Errorf("failed to parse json: %w", err)
	}
	return data, nil
}

// Experimental!
// Client to write settings to an iCON device.
type IconClientWriter interface {
	// Experimental!
	// Sets the thermostat settings
	SetThermostatSettings(tab int, thermosSettings model.ThermostatSettings) error
	// Experimental!
	// Sets the general settings
	SetGeneralSettings(tab int, generalSettings *model.GeneralSettings) error
}

// Experimental!
func (client *iconHttpClient) SetThermostatSettings(tab int, thermosSettings model.ThermostatSettings) error {
	timer := metrics.NewTimer()
	formData := getValues(thermosSettings.ToValues(tab))
	url, err := client.getPath("index.php")
	if err != nil {
		client.session.HttpClientRequest("set_thermostat_settings", 0, timer.End())
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		client.session.HttpClientRequest("set_thermostat_settings", 0, timer.End())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		client.session.HttpClientRequest("set_thermostat_settings", 0, timer.End())
		client.removeSession()
		return err
	}
	defer res.Body.Close()
	client.session.HttpClientRequest("set_thermostat_settings", res.StatusCode, timer.End())
	if res.StatusCode != 200 {
		client.removeSession()
		return fmt.Errorf("failed to read data, status code %d", res.StatusCode)
	}
	client.updateCookie(res.Cookies())
	data := &model.ActionResponse{}
	err = unmarshalBody(res, data)
	if err != nil {
		client.removeSession()
		return err
	}
	if !data.IsSuccess() {
		return data.CreateError()
	}
	return nil
}

// Experimental!
func (client *iconHttpClient) SetGeneralSettings(tab int, generalSettings *model.GeneralSettings) error {
	timer := metrics.NewTimer()
	formData := getValues(generalSettings.ToValues(tab))
	url, err := client.getPath("index.php")
	if err != nil {
		client.session.HttpClientRequest("set_general_settings", 0, timer.End())
		return err
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		client.session.HttpClientRequest("set_general_settings", 0, timer.End())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		client.session.HttpClientRequest("set_general_settings", 0, timer.End())
		client.removeSession()
		return err
	}
	defer res.Body.Close()
	client.session.HttpClientRequest("set_general_settings", res.StatusCode, timer.End())
	if res.StatusCode != http.StatusOK {
		client.removeSession()
		return fmt.Errorf("failed to read data, status code %d", res.StatusCode)
	}
	client.updateCookie(res.Cookies())
	data := &model.ActionResponse{}
	err = unmarshalBody(res, data)
	if err != nil {
		client.removeSession()
		return err
	}
	if !data.IsSuccess() {
		return data.CreateError()
	}
	return nil
}

// Experimental!
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
