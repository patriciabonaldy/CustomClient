package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrNotFound = errors.New("resource not found")
var ErrURLIsEmpty = errors.New("request does not have url")
var ErrBodyIsEmpty = errors.New("request does not have body")

type header struct {
	Key   string
	Value string
}

type Client interface {
	Delete(ctx context.Context, url string, body io.Reader) error
	Get(ctx context.Context, url string) (resp *http.Response, err error)
	Post(ctx context.Context, url string, body io.Reader, headers ...header) (resp *http.Response, err error)
}

// Client defines the communication client.
type client struct {
	httpClient *http.Client
}

func New() Client {
	return &client{
		httpClient: &http.Client{Transport: &http.Transport{}},
	}
}

func (c *client) Delete(ctx context.Context, url string, body io.Reader) error {
	req, err := c.createHeader(ctx, http.MethodDelete, url, nil, nil)
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
	req, err := c.createHeader(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err = c.do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) Post(ctx context.Context, url string, body io.Reader, headers ...header) (resp *http.Response, err error) {
	req, err := c.createHeader(ctx, http.MethodPost, url, body, headers)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, ErrBodyIsEmpty
	}

	resp, err = c.do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) createHeader(ctx context.Context, method, url string, body io.Reader, headers []header) (*http.Request, error) {
	if url == "" {
		return nil, ErrURLIsEmpty
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to make request [%s:%s]: %w", req.Method, req.URL.String(), err)
	}

	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	return req, nil
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed doing request [%s:%s]: %w", req.Method, req.URL.String(), err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return resp, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("failed to do request, %d status code received", resp.StatusCode)
	}
}
