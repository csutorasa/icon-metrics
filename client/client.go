package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/csutorasa/icon-metrics/metrics"
)

type IconClient interface {
	io.Closer
	Login() error
	ReadValues() (*DataPollResponse, error)
	Logout() error
	IsLoggedIn() bool
	SysId() string
}

type httpIconClient struct {
	client    *http.Client
	url       *url.URL
	sysId     string
	password  string
	sessionId string
}

const phpSessionId = "PHPSESSID"

func NewIconClient(urlStr string, sysId string, password string) (IconClient, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	metrics.ConntectedGauge.WithLabelValues(sysId).Set(0)
	return &httpIconClient{
		client:    &http.Client{},
		url:       u,
		sysId:     sysId,
		password:  password,
		sessionId: "",
	}, nil
}

func (this *httpIconClient) Login() error {
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
	actionResponse := actionResponse{}
	err = unmarshalBody(res, &actionResponse)
	if err != nil {
		this.removeSession()
		return err
	}
	if actionResponse.Result != "success" {
		this.removeSession()
		return errors.New("Failed to login")
	}
	metrics.ConntectedGauge.WithLabelValues(this.sysId).Set(1)
	return nil
}

func (this *httpIconClient) ReadValues() (*DataPollResponse, error) {
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

func (this *httpIconClient) Logout() error {
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

func (this *httpIconClient) IsLoggedIn() bool {
	return this.sessionId != ""
}

func (this *httpIconClient) SysId() string {
	return this.sysId
}

func (this *httpIconClient) Close() error {
	err := this.Logout()
	metrics.ConntectedGauge.DeleteLabelValues(this.sysId)
	return err
}

func (this *httpIconClient) getPath(p string) *url.URL {
	u, err := url.Parse(this.url.String())
	if err != nil {

	}
	u.Path = path.Join(u.Path, p)
	return u
}

func (this *httpIconClient) removeSession() {
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

func (client *httpIconClient) updateCookie(cookies []*http.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == phpSessionId {
			client.sessionId = cookie.Value
			return true
		}
	}
	return false
}
