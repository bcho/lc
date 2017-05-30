package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrClientAuthNotPrepared = errors.New("client auth not prepared yet")
)

// ClientAuth provides LeanCloud dashboard authentication source.
type ClientAuth interface {
	// PrepareRequest adds auth settings to given request.
	PrepareRequest(req *http.Request) error

	// WriteTo serializes a client auth.
	WriteTo(io.Writer) (int64, error)
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

	resp, err := http.Post(UrlSignin(), "application/json", payload)
	if err != nil {
		return nil, err
	}

	return NewClientAuthFromCookies(resp.Cookies())
}

// NewClientAuthFromCookies creates an auth source from cookies.
func NewClientAuthFromCookies(cookies []*http.Cookie) (ClientAuth, error) {
	a := &clientAuthImpl{cookies: []*http.Cookie{}}
	for _, cookie := range cookies {
		a.cookies = append(a.cookies, &http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}
	return a, nil
}

// NewClientAuth creates an auth source from io.Reader.
func NewClientAuth(r io.Reader) (ClientAuth, error) {
	var cookies []*http.Cookie

	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("Cookie", string(raw))
	request := http.Request{Header: header}
	cookies = append(cookies, request.Cookies()...)

	return NewClientAuthFromCookies(cookies)
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

func (a clientAuthImpl) WriteTo(w io.Writer) (int64, error) {
	b := new(bytes.Buffer)
	for _, c := range a.cookies {
		b.WriteString(c.String() + ";")
	}
	return b.WriteTo(w)
}

// Client represents a LeanCloud dashboard API client.
type Client interface {
	// Get performs get request to given url.
	Get(string) (*http.Response, error)
	// Do performs a HTTP request.
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	Auth ClientAuth
}

// NewClient creates a API client from client auth.
func NewClient(a ClientAuth) Client {
	return &client{a}
}

func (c client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c client) Do(req *http.Request) (*http.Response, error) {
	if err := c.Auth.PrepareRequest(req); err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
