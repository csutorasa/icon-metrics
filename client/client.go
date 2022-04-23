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
	"strconv"
	"strings"
	"time"

	"github.com/csutorasa/icon-metrics/metrics"
)

type IconClient struct {
	client    *http.Client
	url       *url.URL
	sysId     string
	password  string
	sessionId string
}

// session cookie name
const phpSessionId = "PHPSESSID"

// Default body size is <1kB and each DP is <1kB.
// Each error should be <1kB.
// 1 MB should be a safe hard limit.
const maxReadBytes = 1024 * 1024

// Creates a new client to fetch data from an iCON device.
func NewIconClient(urlStr string, sysId string, password string) (*IconClient, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return &IconClient{
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
	}, nil
}

// Logs in and creates a session.
func (client *IconClient) Login() error {
	timer := metrics.NewTimer()
	formData := url.Values{
		"sysid":    []string{client.sysId},
		"password": []string{client.password},
		"lang":     []string{"hu"},
		"tab":      []string{"login"},
		"form":     []string{"login"},
	}
	req, err := http.NewRequest("POST", client.url.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "login", "0").Observe(timer.End().Seconds())
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "login", "0").Observe(timer.End().Seconds())
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(client.sysId, "login", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login, status code %d", res.StatusCode)
	}
	if !client.updateCookie(res.Cookies()) {
		return errors.New("no session was found")
	}
	data := actionResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		client.removeSession()
		return err
	}
	if !data.IsSuccess() {
		client.removeSession()
		return data.CreateError()
	}
	return nil
}

// Reads data from the device.
func (client *IconClient) ReadValues() (*DataPollResponse, error) {
	timer := metrics.NewTimer()
	fomrData := url.Values{
		"tab": []string{"datapoll"},
	}
	url, err := client.getPath("index.php")
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "read_values", "0").Observe(timer.End().Seconds())
		return nil, err
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "read_values", "0").Observe(timer.End().Seconds())
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "read_values", "0").Observe(timer.End().Seconds())
		client.removeSession()
		return nil, err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(client.sysId, "read_values", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != http.StatusOK {
		client.removeSession()
		return nil, fmt.Errorf("failed to read data, status code %d", res.StatusCode)
	}
	client.updateCookie(res.Cookies())
	data := &DataPollResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		client.removeSession()
		return data, err
	}
	return data, nil
}

// Experimental!
func (client *IconClient) SetThermostatSettings(tab int, thermosSettings ThermostatSettings) error {
	timer := metrics.NewTimer()
	fomrData := getValues(thermosSettings.ToValues(tab))
	url, err := client.getPath("index.php")
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "set_thermostat_settings", "0").Observe(timer.End().Seconds())
		return err
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "set_thermostat_settings", "0").Observe(timer.End().Seconds())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "set_thermostat_settings", "0").Observe(timer.End().Seconds())
		client.removeSession()
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(client.sysId, "set_thermostat_settings", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != 200 {
		client.removeSession()
		return fmt.Errorf("failed to read data, status code %d", res.StatusCode)
	}
	client.updateCookie(res.Cookies())
	data := &actionResponse{}
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
func (client *IconClient) SetGeneralSettings(tab int, generalSettings *GeneralSettings) error {
	timer := metrics.NewTimer()
	fomrData := getValues(generalSettings.ToValues(tab))
	url, err := client.getPath("index.php")
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "set_general_settings", "0").Observe(timer.End().Seconds())
		return err
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "set_general_settings", "0").Observe(timer.End().Seconds())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "set_general_settings", "0").Observe(timer.End().Seconds())
		client.removeSession()
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(client.sysId, "set_general_settings", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != http.StatusOK {
		client.removeSession()
		return fmt.Errorf("failed to read data, status code %d", res.StatusCode)
	}
	client.updateCookie(res.Cookies())
	data := &actionResponse{}
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

// Closes a session.
func (client *IconClient) Logout() error {
	timer := metrics.NewTimer()
	fomrData := url.Values{
		"logout": []string{"true"},
	}
	url, err := client.getPath("index.php")
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "logout", "0").Observe(timer.End().Seconds())
		return err
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "logout", "0").Observe(timer.End().Seconds())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: client.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(client.sysId, "logout", "0").Observe(timer.End().Seconds())
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(client.sysId, "logout", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	client.removeSession()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to logout, status code %d", res.StatusCode)
	}
	return nil
}

// Returns if there is a session.
func (client *IconClient) IsLoggedIn() bool {
	return client.sessionId != ""
}

// Returns the system ID.
func (client *IconClient) SysId() string {
	return client.sysId
}

// Cleans up the client.
func (client *IconClient) Close() error {
	err := client.Logout()
	return err
}

// Join path to the base URL.
func (client *IconClient) getPath(p string) (*url.URL, error) {
	u, err := url.Parse(client.url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to creat url from %s %w", client.url.String(), err)
	}
	u.Path = path.Join(u.Path, p)
	return u, nil
}

// Removes metrics and session data.
func (client *IconClient) removeSession() {
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

// Unpdates session data from cookies.
func (client *IconClient) updateCookie(cookies []*http.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == phpSessionId {
			client.sessionId = cookie.Value
			return true
		}
	}
	return false
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
