package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/csutorasa/icon-metrics/pkg/model"
)

// Join path to the base URL.
func getUrlWithPath(baseUrl *url.URL, p string) string {
	u := new(url.URL)
	*u = *baseUrl
	u.Path = path.Join(u.Path, p)
	return u.String()
}

func postFormRequestContext(ctx context.Context, url string, formData url.Values) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func postFormRequestContextWithSession(ctx context.Context, url string, formData url.Values, session *IconSession) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	err = addSesson(req, session)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// Default body size is <1kB and each DP is <1kB.
// Each error should be <1kB.
// 1 MB should be a safe hard limit.
const maxReadBytes = 1024*1024 + 1

// Unmarshal JSON content from http response body.
func unmarshalBody(res *http.Response, v any) error {
	err := validateStatusCode(res)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, maxReadBytes))
	if err != nil {
		return err
	}
	if len(body) >= maxReadBytes {
		return &ErrHttpBodyUnmarshal{Cause: fmt.Errorf("too long response body")}
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return &ErrHttpBodyUnmarshal{Cause: err}
	}
	return nil
}

func validateStatusCode(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		return &ErrHttpStatus{StatusCode: res.StatusCode}
	}
	return nil
}

func validateAction(res *http.Response) error {
	err := validateStatusCode(res)
	if err != nil {
		return err
	}
	data := model.ActionResponse{}
	err = unmarshalBody(res, &data)
	if err != nil {
		return err
	}
	if !data.IsSuccess() {
		return data.CreateError()
	}
	return nil
}

// session cookie name
const phpSessionId = "PHPSESSID"

// Finds session from response.
func findSesson(res *http.Response) (*IconSession, bool) {
	for _, cookie := range res.Cookies() {
		if cookie.Name == phpSessionId {
			return &IconSession{
				sessionId: cookie.Value,
			}, true
		}
	}
	return nil, false
}

// Updates session from response.
func updateSession(res *http.Response, session *IconSession) {
	if IsAuthErrorStatusCode(res.StatusCode) {
		session.sessionId = ""
		return
	}
	s, ok := findSesson(res)
	if ok {
		*session = *s
	}
}

// Adds session to request.
func addSesson(req *http.Request, session *IconSession) error {
	if !session.Valid() {
		return fmt.Errorf("session is invalid")
	}
	sessionIds := req.CookiesNamed(phpSessionId)
	if len(sessionIds) == 0 {
		req.AddCookie(&http.Cookie{Name: phpSessionId, Value: session.sessionId})
		return nil
	}
	if len(sessionIds) == 1 {
		sessionIds[0].Value = session.sessionId
		return nil
	}
	return fmt.Errorf("multiple session cookies found")
}
