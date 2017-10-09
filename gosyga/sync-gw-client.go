package gosyga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ClientApi struct {
	bucket string
	url    string
}

type SessionPasswordRequest struct {
	UserName     string `json:"name"`
	UserPassword string `json:"password"`
}

func NewClientApi(url string, bucket string) *ClientApi {
	return &ClientApi{
		bucket: bucket,
		url:    url,
	}
}

func (c *ClientApi) CreateSession(username string, password string) (*SessionResponse, error) {
	url := c.url + "/" + url.QueryEscape(c.bucket) + "/_session"

	sessReq := SessionPasswordRequest{
		UserName:     username,
		UserPassword: password,
	}

	data, err := json.Marshal(sessReq)
	if err != nil {
		return nil, err
	}

	resp, err := Do_POST(url, data)
	if err != nil {
		return nil, err
	}

	if resp.Code == 401 || resp.Code == 403 {
		return nil, nil // wrong username/password pair
	}

	var sessCookie *http.Cookie
	for _, cookie := range resp.Cookies {
		// TODO: select first cookie instead?
		if cookie.Name == "SyncGatewaySession" {
			sessCookie = cookie
		}
	}

	if sessCookie == nil {
		return nil, fmt.Errorf("Can't find SyncGatewaySession cookie in server response headers")
	}

	sessionResponse := &SessionResponse{
		CookieName: sessCookie.Name,
		SessionId:  sessCookie.Value,
		Expires:    sessCookie.Expires.UTC().Format("2006-01-02T15:04:05-0700"),
	}

	return sessionResponse, nil
}
