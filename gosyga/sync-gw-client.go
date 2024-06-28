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

	apiWithLogger
}

type sessionPasswordRequest struct {
	UserName     string `json:"name"`
	UserPassword string `json:"password"`
}

func NewClientApi(url string, bucket string, user, password string) *ClientApi {
	return &ClientApi{
		bucket:        bucket,
		url:           url,
		apiWithLogger: newNullApiLogger(user, password),
	}
}

func (c *ClientApi) CreateSession(username string, password string) (*SessionToken, error) {
	url := c.url + "/" + url.QueryEscape(c.bucket) + "/_session"

	sessReq := sessionPasswordRequest{
		UserName:     username,
		UserPassword: password,
	}

	data, err := json.Marshal(sessReq)
	if err != nil {
		return nil, err
	}

	resp, err := c.doPOST(url, data)
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

	token := &SessionToken{
		CookieName: sessCookie.Name,
		SessionId:  sessCookie.Value,
		Expires:    sessCookie.Expires.UTC().Format("2006-01-02T15:04:05-0700"),
	}

	return token, nil
}
