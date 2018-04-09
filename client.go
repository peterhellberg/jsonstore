package jsonstore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Errors
var (
	ErrNotFound            = fmt.Errorf("not found")
	ErrInternalServerError = fmt.Errorf("internal server error")
	ErrUnexpectedStatus    = fmt.Errorf("unexpected status")
)

const (
	defaultScheme    = "https"
	defaultHost      = "www.jsonstore.io"
	defaultUserAgent = "jsonstore/client.go (godoc.org/github.com/peterhellberg/jsonstore)"
	defaultTimeout   = 10 * time.Second
)

// Client for the www.jsonstore.io API
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	userAgent  string
	secret     string
}

// Option is a type of function used to configure the client
type Option func(*Client)

// New creates a www.jsonstore.io Client
func New(options ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: &url.URL{
			Scheme: defaultScheme,
			Host:   defaultHost,
		},
		userAgent: defaultUserAgent,
	}

	for _, option := range options {
		option(c)
	}

	if c.secret == "" {
		c.secret = NewSecret()
	}

	return c
}

// HTTPClient changes the HTTP client used by the client to the provided *http.Client
func HTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// BaseURL changes the base URL used by the client to the URL parsed from rawurl
func BaseURL(rawurl string) Option {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// Secret sets the secret used by the client to the provided string
func Secret(s string) Option {
	return func(c *Client) {
		c.secret = s
	}
}

// NewSecret generates a new secret based on the current time
func NewSecret() string {
	return sha256sum(time.Now().Format(time.RFC3339Nano))
}

// Secret returns the client secret
func (c *Client) Secret() string {
	return c.secret
}

// URL returns the URL used by the client
func (c *Client) URL(segments ...string) *url.URL {
	return c.baseURL.ResolveReference(&url.URL{
		Path: "/" + c.secret + "/" + strings.TrimPrefix(strings.Join(segments, "/"), "/"),
	})
}

// Get response from jsonstore
func (c *Client) Get(ctx context.Context, path string, v interface{}) error {
	req, err := c.request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	return c.do(req, v)
}

// Post to jsonstore
func (c *Client) Post(ctx context.Context, path string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	req, err := c.request(ctx, http.MethodPost, path, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// Put update to jsonstore
func (c *Client) Put(ctx context.Context, path string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	req, err := c.request(ctx, http.MethodPut, path, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// Delete from jsonstore
func (c *Client) Delete(ctx context.Context, path string) error {
	req, err := c.request(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) request(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse("/" + c.secret + "/" + strings.TrimPrefix(path, "/"))
	if err != nil {
		return nil, err
	}

	rawurl := c.baseURL.ResolveReference(rel).String()

	req, err := http.NewRequest(method, rawurl, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.CopyN(ioutil.Discard, resp.Body, 1024)
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusInternalServerError:
		return ErrInternalServerError
	default:
		return ErrUnexpectedStatus
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(&response{v})
	}

	return nil
}

func sha256sum(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

type response struct {
	Result interface{} `json:"result"`
}
