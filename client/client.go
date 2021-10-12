package client

import (
	"encoding/json"
	"errors"
	"fmt"
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

const phpSessionId = "PHPSESSID"

func NewIconClient(urlStr string, sysId string, password string) (*IconClient, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	metrics.ConntectedGauge.WithLabelValues(sysId).Set(0)
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

func (this *IconClient) Login() error {
	timer := metrics.NewTimer()
	formData := url.Values{
		"sysid":    []string{this.sysId},
		"password": []string{this.password},
		"lang":     []string{"hu"},
		"tab":      []string{"login"},
		"form":     []string{"login"},
	}
	req, err := http.NewRequest("POST", this.url.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "login", "0").Observe(timer.End().Seconds())
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := this.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "login", "0").Observe(timer.End().Seconds())
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(this.sysId, "login", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed to login, status code %d", res.StatusCode))
	}
	if !this.updateCookie(res.Cookies()) {
		return errors.New("No session was found")
	}
	data := actionResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		this.removeSession()
		return err
	}
	if !data.IsSuccess() {
		this.removeSession()
		return data.CreateError()
	}
	metrics.ConntectedGauge.WithLabelValues(this.sysId).Set(1)
	return nil
}

func (this *IconClient) ReadValues() (*DataPollResponse, error) {
	timer := metrics.NewTimer()
	fomrData := url.Values{
		"tab": []string{"datapoll"},
	}
	req, err := http.NewRequest("POST", this.getPath("index.php").String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "read_values", "0").Observe(timer.End().Seconds())
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: this.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := this.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "read_values", "0").Observe(timer.End().Seconds())
		this.removeSession()
		return nil, err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(this.sysId, "read_values", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != 200 {
		this.removeSession()
		return nil, errors.New(fmt.Sprintf("Failed to read data, status code %d", res.StatusCode))
	}
	this.updateCookie(res.Cookies())
	data := &DataPollResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		this.removeSession()
		return data, err
	}
	return data, nil
}

func (this *IconClient) SetThermostatSettings(tab int, thermosSettings ThermostatSettings) error {
	timer := metrics.NewTimer()
	fomrData := getValues(thermosSettings.ToValues(tab))
	req, err := http.NewRequest("POST", this.getPath("index.php").String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "set_thermostat_settings", "0").Observe(timer.End().Seconds())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: this.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := this.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "set_thermostat_settings", "0").Observe(timer.End().Seconds())
		this.removeSession()
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(this.sysId, "set_thermostat_settings", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != 200 {
		this.removeSession()
		return errors.New(fmt.Sprintf("Failed to read data, status code %d", res.StatusCode))
	}
	this.updateCookie(res.Cookies())
	data := &actionResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		this.removeSession()
		return err
	}
	if !data.IsSuccess() {
		return data.CreateError()
	}
	return nil
}

func (this *IconClient) SetGeneralSettings(tab int, generalSettings *GeneralSettings) error {
	timer := metrics.NewTimer()
	fomrData := getValues(generalSettings.ToValues(tab))
	req, err := http.NewRequest("POST", this.getPath("index.php").String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "set_general_settings", "0").Observe(timer.End().Seconds())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: this.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := this.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "set_general_settings", "0").Observe(timer.End().Seconds())
		this.removeSession()
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(this.sysId, "set_general_settings", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	if res.StatusCode != 200 {
		this.removeSession()
		return errors.New(fmt.Sprintf("Failed to read data, status code %d", res.StatusCode))
	}
	this.updateCookie(res.Cookies())
	data := &actionResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		this.removeSession()
		return err
	}
	if !data.IsSuccess() {
		return data.CreateError()
	}
	return nil
}

func (this *IconClient) Logout() error {
	timer := metrics.NewTimer()
	fomrData := url.Values{
		"logout": []string{"true"},
	}
	req, err := http.NewRequest("POST", this.getPath("index.php").String(), strings.NewReader(fomrData.Encode()))
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "logout", "0").Observe(timer.End().Seconds())
		return err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionId, Value: this.sessionId})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := this.client.Do(req)
	if err != nil {
		metrics.HttpGauge.WithLabelValues(this.sysId, "logout", "0").Observe(timer.End().Seconds())
		return err
	}
	defer res.Body.Close()
	metrics.HttpGauge.WithLabelValues(this.sysId, "logout", strconv.Itoa(res.StatusCode)).Observe(timer.End().Seconds())
	this.removeSession()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed to logout, status code %d", res.StatusCode))
	}
	return nil
}

func (this *IconClient) IsLoggedIn() bool {
	return this.sessionId != ""
}

func (this *IconClient) SysId() string {
	return this.sysId
}

func (this *IconClient) Close() error {
	err := this.Logout()
	metrics.ConntectedGauge.DeleteLabelValues(this.sysId)
	return err
}

func (this *IconClient) getPath(p string) *url.URL {
	u, err := url.Parse(this.url.String())
	if err != nil {

	}
	u.Path = path.Join(u.Path, p)
	return u
}

func (this *IconClient) removeSession() {
	metrics.ConntectedGauge.WithLabelValues(this.sysId).Set(0)
	this.sessionId = ""
}

func unmarshalBody(res *http.Response, v interface{}) error {
	body, err := getBodyBytes(res)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func getBodyBytes(res *http.Response) ([]byte, error) {
	body := make([]byte, 8092)
	read, err := res.Body.Read(body)
	if err.Error() != "EOF" {
		return nil, err
	}
	return body[0:read], nil
}

func (client *IconClient) updateCookie(cookies []*http.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == phpSessionId {
			client.sessionId = cookie.Value
			return true
		}
	}
	return false
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
