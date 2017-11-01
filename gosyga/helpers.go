package gosyga

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

// General json document with unknown fields
type JsonDoc map[string]interface{}

// This helper-structure holds data from more general repsonse of net/http library
type httpResponse struct {
	Code    int
	Body    []byte
	Cookies []*http.Cookie
}

/*
 * Some http methods for general api-class with logging.
 */
func (a *apiWithLogger) doGET(url string) (*httpResponse, error) {
	return a.sendJsonRequest("GET", url, nil)
}

func (a *apiWithLogger) doPOST(url string, data []byte) (*httpResponse, error) {
	return a.sendJsonRequest("POST", url, data)
}

func (a *apiWithLogger) doPUT(url string, data []byte) (*httpResponse, error) {
	return a.sendJsonRequest("PUT", url, data)
}

func (a *apiWithLogger) doDELETE(url string) (*httpResponse, error) {
	return a.sendRequest("DELETE", url, nil, false)
}

func (a *apiWithLogger) sendJsonRequest(method string, url string, data []byte) (*httpResponse, error) {
	return a.sendRequest(method, url, data, true)
}

func (a *apiWithLogger) sendRequest(method string, url string, data []byte, isJson bool) (*httpResponse, error) {
	a.log.Debugf("%s %s", method, url)

	var req *http.Request
	var err error

	if data != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	if isJson {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	a.log.Debug("response Status: ", resp.Status)
	a.log.Debug("response Headers: ", resp.Header)

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	a.log.Debug("response Body:", string(bytes))

	return &httpResponse{
		Code:    resp.StatusCode,
		Body:    bytes,
		Cookies: resp.Cookies(),
	}, nil
}
