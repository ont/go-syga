package gosyga

import (
	"encoding/json"
	"fmt"
	"net/url"
)

//
// TODO: replace this file with swagger autogenerated code?
//

type AdminApi struct {
	bucket string
	url    string

	apiWithLogger
}

type User struct {
	Name          string   `json:"name"`
	Password      string   `json:"password,omitempty"`
	AdminChannels []string `json:"admin_channels"`
	AllChannels   []string `json:"all_channels"`
	AdminRoles    []string `json:"admin_roles"`
	Email         string   `json:"email,omitempty"`
	Disabled      bool     `json:"disabled"`
}

type SessionToken struct {
	CookieName string `json:"cookie_name"`
	Expires    string `json:"expires"`
	SessionId  string `json:"session_id"`
}

type SessionInfo struct {
	Valid bool `json:"ok"`
	User  struct {
		Username string         `json:"name"`
		Channels map[string]int `json:"channels"`
	} `json:"userCtx"`
}

type sessionRequest struct {
	UserName string `json:"name"`
	TTL      int    `json:"ttl"`
}

func NewAdminApi(url string, bucket string) *AdminApi {
	return &AdminApi{
		bucket:        bucket,
		url:           url,
		apiWithLogger: newNullApiLogger(),
	}
}

func (a *AdminApi) GetUser(uuid string) (*User, error) {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/_user/" + url.QueryEscape(uuid)

	resp, err := a.doGET(url)

	if err != nil {
		return nil, err
	}

	if resp.Code == 404 {
		return nil, nil // user doesn't exists in database
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("Can't get user, got non-200 response code: %d", resp.Code)
	}

	var user User
	err = json.Unmarshal(resp.Body, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *AdminApi) CreateUser(uuid string, password string) (*User, error) {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/_user/"

	user := User{
		Name:     uuid,
		Password: password,
		Disabled: false,
	}

	data, err := json.Marshal(user)

	if err != nil {
		return nil, err
	}

	_, err = a.doPOST(url, data)

	if err != nil {
		return nil, err
	}

	return a.GetUser(uuid)
}

func (a *AdminApi) UpdateUser(user *User) (*User, error) {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/_user/" + url.QueryEscape(user.Name)

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	_, err = a.doPUT(url, data)
	if err != nil {
		return nil, err
	}

	return a.GetUser(user.Name)
}

func (a *AdminApi) CreateSession(username string) (*SessionToken, error) {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/_session"

	sessReq := sessionRequest{
		UserName: username,
		TTL:      24 * 3600, // 24 hours
	}

	data, err := json.Marshal(sessReq)
	if err != nil {
		return nil, err
	}

	resp, err := a.doPOST(url, data)

	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("Can't create session, got non-200 response code: %d", resp.Code)
	}

	var token SessionToken
	err = json.Unmarshal(resp.Body, &token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Simple version of GET /{db}/{doc} call.
// Unmarshal any valid response into variable "v".
func (a *AdminApi) GetDoc(docId string, v interface{}) (found bool, err error) {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/" + docId

	resp, err := a.doGET(url)
	if err != nil {
		return false, err
	}

	if resp.Code == 404 {
		return false, nil
	}

	if resp.Code != 200 {
		return false, fmt.Errorf("Can't get document %s, got non-200 response code: %d", docId, resp.Code)
	}

	err = json.Unmarshal(resp.Body, v)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *AdminApi) GetSession(sessionId string) (*SessionInfo, error) {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/_session/" + sessionId

	resp, err := a.doGET(url)
	if err != nil {
		return nil, err
	}

	if resp.Code == 404 {
		return nil, nil
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("Can't get session %s, got non-200 response code: %d", sessionId, resp.Code)
	}

	var sessionInfo SessionInfo
	err = json.Unmarshal(resp.Body, &sessionInfo)

	return &sessionInfo, nil
}

// Method retrieves document from database then pass it to the "fields" callback as json string.
// Return value of "fields" callback must be field-value pairs which should be updated in original document.
// This method will try to save document in database several times before returing error.
func (a *AdminApi) UpdateDoc(docId string, fields func(bytes []byte) (JsonDoc, error)) error {
	errs := make([]string, 0)

	for try := 0; ; try++ {
		// TODO: make number of retries configurable
		if try > 5 {
			return fmt.Errorf("Too many tries during document %s update: %#v", docId, errs)
		}

		bytes, err := a.GetRawDoc(docId)
		if err != nil {
			return err
		}

		if bytes == nil {
			return fmt.Errorf("Document with id %s is not found", docId)
		}

		fs, err := fields(bytes)
		if err != nil {
			return err
		}

		var doc JsonDoc
		err = json.Unmarshal(bytes, &doc)
		if err != nil {
			return err
		}

		for key, value := range fs {
			doc[key] = value
		}

		bytes, err = json.Marshal(doc)
		if err != nil {
			return err
		}

		err = a.UpdateRawDoc(docId, bytes)
		if err == nil {
			return nil
		}
		errs = append(errs, err.Error())
	}
}

func (a *AdminApi) GetRawDoc(docId string) ([]byte, error) {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/" + docId
	resp, err := a.doGET(url)

	if err != nil {
		return nil, err
	}

	if resp.Code == 404 {
		return nil, nil // document doesn't exists in database
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("Can't get document, got non-200 response code: %d", resp.Code)
	}

	return resp.Body, nil
}

func (a *AdminApi) UpdateRawDoc(docId string, bytes []byte) error {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/" + docId

	resp, err := a.doPUT(url, bytes)

	if err != nil {
		return err
	}

	if resp.Code != 200 && resp.Code != 201 {
		return fmt.Errorf("Error during updating document %s: http response code not in (200, 201): %d", resp.Code)
	}

	return nil
}

func (a *AdminApi) DeleteDoc(docId string) error {
	url := a.url + "/" + url.QueryEscape(a.bucket) + "/" + docId

	var doc struct {
		Revision string `json:"_rev"`
	}

	found, err := a.GetDoc(docId, &doc)

	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("Can't find document with id %s in database", docId)
	}

	resp, err := a.doDELETE(url + "?rev=" + doc.Revision)

	if resp.Code != 200 && resp.Code != 201 {
		return fmt.Errorf("Error during deleting document %s: http response code not in (200, 201): %d", resp.Code)
	}

	return err
}
