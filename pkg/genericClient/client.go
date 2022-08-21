package genericClient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New("resource not found")
	// ErrURLIsEmpty error
	ErrURLIsEmpty = errors.New("request does not have url")
	// ErrBodyIsEmpty error
	ErrBodyIsEmpty = errors.New("request does not have body")
)

// Header represents Header in the request.
type Header struct {
	Key   string
	Value string
}

// Client handler different http methods
type Client interface {
	Delete(ctx context.Context, url string, headers ...Header) error
	Get(ctx context.Context, url string) (resp *http.Response, err error)
	Post(ctx context.Context, url string, data []byte, headers ...Header) (resp *http.Response, err error)
}

// Client defines the communication client.
type client struct {
	httpClient        *http.Client
	retryRoundOptions *Options
}

// New create a new client
func New(options ...Option) Client {
	r := setupOptions(options...)
	return &client{
		httpClient: &http.Client{
			Transport: &http.Transport{},
			Timeout:   time.Duration(r.TimeDuration) + 5*time.Second,
		},
		retryRoundOptions: r,
	}
}

func setupOptions(options ...Option) *Options {
	r := &Options{
		TimeDuration: 10,
	}
	for _, fn := range options {
		fn(r)
	}

	return r
}

func (c *client) Delete(ctx context.Context, url string, headers ...Header) error {
	req, err := c.withHeader(ctx, http.MethodDelete, url, nil, headers)
	if err != nil {
		return err
	}

	_, err = c.do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := c.withHeader(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err = c.do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) Post(ctx context.Context, url string, data []byte, headers ...Header) (resp *http.Response, err error) {
	if len(data) == 0 {
		return nil, ErrBodyIsEmpty
	}

	req, err := c.withHeader(ctx, http.MethodPost, url, data, headers)
	if err != nil {
		return nil, err
	}

	resp, err = c.do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type Request struct {
	*http.Request
	method  string
	url     string
	headers []Header
	data    []byte
}

func (c *client) withHeader(ctx context.Context, method, url string, data []byte, headers []Header) (*Request, error) {
	if url == "" {
		return nil, ErrURLIsEmpty
	}

	body := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to make request [%s:%s]: %w", req.Method, req.URL.String(), err)
	}

	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}

	return &Request{
		Request: req,
		method:  method,
		url:     url,
		headers: headers,
		data:    data,
	}, nil
}

func (c *client) do(req *Request) (resp *http.Response, err error) {
	for {
		resp, err = c.httpClient.Do(req.Request)
		if err != nil {
			return nil, fmt.Errorf("failed doing request [%s:%s]: %w", req.Method, req.URL.String(), err)
		}

		if !checkRetry(c, resp) {
			break
		}

		resetBody(req)
		c.retryRoundOptions.MaxRetryCount--
		<-time.After(c.retryRoundOptions.CalculateBackoff(c.retryRoundOptions.MaxRetryCount))
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("failed to do request, %d status code received", resp.StatusCode)
	}
}

func checkRetry(c *client, resp *http.Response) bool {
	if !c.retryRoundOptions.ShouldRetry {
		return false
	}

	if c.retryRoundOptions.MaxRetryCount <= 0 {
		return false
	}

	return resp.StatusCode == 0 || resp.StatusCode >= 500
}

func resetBody(req *Request) {
	req.Request.Body = io.NopCloser(bytes.NewBuffer(req.data))
	req.Request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(req.data)), nil
	}
}
