package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrClientAuthNotPrepared = errors.New("client auth not prepared yet")
)

// ClientAuth provides LeanCloud dashboard authentication source.
type ClientAuth interface {
	// PrepareRequest adds auth settings to given request.
	PrepareRequest(req *http.Request) error
}

type clientAuthImpl struct {
	// cookie login
	cookies []*http.Cookie
}

// NewClientAuthFromLogin creates an auth source from password based credentials.
func NewClientAuthFromLogin(email, password string) (ClientAuth, error) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(map[string]string{
		"email":    email,
		"password": password,
	}) // should not fail here

	resp, err := http.Post(UrlSignin, "application/json", payload)
	if err != nil {
		return nil, err
	}

	return &clientAuthImpl{resp.Cookies()}, nil
}

func (a clientAuthImpl) PrepareRequest(req *http.Request) error {
	if len(a.cookies) == 0 {
		return ErrClientAuthNotPrepared
	}

	for _, c := range a.cookies {
		req.AddCookie(c)
	}

	return nil
}

// Client represents a LeanCloud dashboard API client.
type Client struct {
	Auth ClientAuth
}

// NewClient creates a API client from client auth.
func NewClient(a ClientAuth) *Client {
	return &Client{a}
}

// Get performs get request to given url.
func (c Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Do performs a HTTP request.
func (c Client) Do(req *http.Request) (*http.Response, error) {
	if err := c.Auth.PrepareRequest(req); err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
