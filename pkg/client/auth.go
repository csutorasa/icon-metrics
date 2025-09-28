package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

type IconSession struct {
	sessionId string
}

func (s *IconSession) Valid() bool {
	return s.sessionId != ""
}

// Manages session to an iCON device.
type IconAuthClient interface {
	Login(ctx context.Context, sysId string, password string) (*IconSession, error)
	Logout(ctx context.Context, session *IconSession) error
}

func NewAuthClient(c *http.Client, baseUrl *url.URL) IconAuthClient {
	return &IconPhpAuthClient{
		client: c,
		url:    baseUrl,
	}
}

// Manages PHP session to an iCON device.
type IconPhpAuthClient struct {
	client *http.Client
	url    *url.URL
}

func (client *IconPhpAuthClient) Login(ctx context.Context, sysId string, password string) (*IconSession, error) {
	req, err := postFormRequestContext(ctx, client.url.String(), loginRequest(sysId, password))
	if err != nil {
		return nil, err
	}
	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = validateAction(res)
	if err != nil {
		return nil, err
	}
	session, ok := findSesson(res)
	if !ok {
		return nil, errors.New("no session was found")
	}
	return session, nil
}

func (client *IconPhpAuthClient) Logout(ctx context.Context, session *IconSession) error {
	req, err := postFormRequestContextWithSession(ctx, getUrlWithPath(client.url, "index.php"), logoutRequest(), session)
	if err != nil {
		return err
	}
	res, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = validateStatusCode(res)
	if err != nil {
		return err
	}
	return nil
}

func loginRequest(sysId string, password string) url.Values {
	return url.Values{
		"sysid":    []string{sysId},
		"password": []string{password},
		"lang":     []string{"hu"},
		"tab":      []string{"login"},
		"form":     []string{"login"},
	}
}

func logoutRequest() url.Values {
	return url.Values{
		"logout": []string{"true"},
	}
}
