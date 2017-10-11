package gosyga

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// General json document with unknown fields
type JsonDoc map[string]interface{}

// This helper-structure holds data from more general repsonse of net/http library
type Response struct {
	Code    int
	Body    []byte
	Cookies []*http.Cookie
}

func Do_GET(url string) (*Response, error) {
	// TODO: add logrus to all fmt.Println...
	fmt.Println("GET>", url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println("response Body:", string(bytes))

	return &Response{
		Code:    resp.StatusCode,
		Body:    bytes,
		Cookies: resp.Cookies(),
	}, nil
}

func Do_POST(url string, data []byte) (*Response, error) {
	return sendJsonRequest("POST", url, data)
}

func Do_PUT(url string, data []byte) (*Response, error) {
	return sendJsonRequest("PUT", url, data)
}

func sendJsonRequest(method string, url string, data []byte) (*Response, error) {
	// TODO: add logrus to all fmt.Println...
	fmt.Println(method+">", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println("response Body:", string(bytes))

	return &Response{
		Code:    resp.StatusCode,
		Body:    bytes,
		Cookies: resp.Cookies(),
	}, nil
}