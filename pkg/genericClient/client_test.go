package genericClient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	req  *http.Request
	resp *http.Response
	err  error
}

func (mt *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	mt.req = req
	return mt.resp, mt.err
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New())
}

func Test_client_Delete(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		mockTransport *mockTransport
		mockHandler   func(w http.ResponseWriter, r *http.Request)
		wantErr       bool
	}{
		{
			name:    "URL is empty",
			url:     "",
			wantErr: true,
		},
		{
			name: "unknown error",
			url:  "http://localhost:8080/anything/1",
			mockTransport: &mockTransport{
				resp: nil,
				err:  errors.New("unknown error"),
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {},
			wantErr:     true,
		},
		{
			name: "error calls get",
			url:  "http://localhost:8080/anything/1",
			mockTransport: &mockTransport{
				resp: &http.Response{
					Status:     fmt.Sprintf("%d", http.StatusInternalServerError),
					StatusCode: http.StatusInternalServerError,
				},
				err: nil,
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, fmt.Sprint(errors.New("unknown error")), http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "success",
			url:  "http://localhost:8080/anything",
			mockTransport: &mockTransport{
				resp: &http.Response{
					Status:     fmt.Sprintf("%d OK", http.StatusOK),
					StatusCode: http.StatusOK,
				},
				err: nil,
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupOptions()
			c := &client{
				httpClient: &http.Client{
					Transport: tt.mockTransport,
					Timeout:   time.Duration(r.TimeDuration) + 5*time.Second,
				},
				retryRoundOptions: r,
			}
			if err := c.Delete(context.Background(), tt.url); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Get(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		mockTransport *mockTransport
		mockHandler   func(w http.ResponseWriter, r *http.Request)
		wantErr       bool
	}{
		{
			name:    "URL is empty",
			url:     "",
			wantErr: true,
		},
		{
			name: "unknown error",
			url:  "http://localhost:8080/anything/1",
			mockTransport: &mockTransport{
				resp: nil,
				err:  errors.New("unknown error"),
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {},
			wantErr:     true,
		},
		{
			name: "error calls get",
			url:  "http://localhost:8080/anything/1",
			mockTransport: &mockTransport{
				resp: &http.Response{
					Status:     fmt.Sprintf("%d", http.StatusInternalServerError),
					StatusCode: http.StatusInternalServerError,
				},
				err: nil,
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, fmt.Sprint(errors.New("unknown error")), http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "success",
			url:  "http://localhost:8080/anything",
			mockTransport: &mockTransport{
				resp: &http.Response{
					Status:     fmt.Sprintf("%d OK", http.StatusOK),
					StatusCode: http.StatusOK,
				},
				err: nil,
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				// nolint:errcheck
				w.Write([]byte(`{
					"data": {
						"id": "f14956e9-a751-4879-9751-eb47001649b4",
						"version": 0,
						"organisation_id": ""
					}}`))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupOptions(WithTimeDuration(3))
			c := &client{
				httpClient: &http.Client{
					Transport: tt.mockTransport,
					Timeout:   time.Duration(r.TimeDuration) + 5*time.Second,
				},
				retryRoundOptions: r,
			}
			if _, err := c.Get(context.Background(), tt.url); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Post(t *testing.T) {
	body := []byte(`{
					"data": {
						"id": "f14956e9-a751-4879-9751-eb47001649b4",
						"version": 0,
						"organisation_id": ""
					}}`)
	tests := []struct {
		name          string
		url           string
		mockTransport *mockTransport
		mockHandler   func(w http.ResponseWriter, r *http.Request)
		header        []Header
		body          []byte
		wantErr       bool
		wantRetry     bool
	}{
		{
			name:    "URL is empty",
			url:     "",
			wantErr: true,
		},
		{
			name: "body is empty",
			url:  "http://localhost:8080/anything/1",
			mockTransport: &mockTransport{
				resp: nil,
				err:  errors.New("unknown error"),
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {},
			wantErr:     true,
		},
		{
			name: "unknown error",
			url:  "http://localhost:8080/anything/1/6",
			mockTransport: &mockTransport{
				resp: nil,
				err:  errors.New("unknown error"),
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {},
			body:        body,
			wantErr:     true,
		},
		{
			name: "error calls post with retry",
			url:  "http://localhost:8080/anything",
			mockTransport: &mockTransport{
				resp: &http.Response{
					Status:     fmt.Sprintf("%d", http.StatusInternalServerError),
					StatusCode: http.StatusInternalServerError,
				},
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, fmt.Sprint(errors.New("unknown error")), http.StatusInternalServerError)
			},
			body:      body,
			wantRetry: true,
			wantErr:   true,
		},
		{
			name: "success",
			url:  "http://localhost:8080/anything",
			mockTransport: &mockTransport{
				resp: &http.Response{
					Status:     fmt.Sprintf("%d OK", http.StatusOK),
					StatusCode: http.StatusOK,
				},
				err: nil,
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			},
			body: body,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupOptions()
			if tt.wantRetry {
				r = setupOptions(WithRetryPolicy(true),
					WithMaxRetryCount(3),
					WithBackoffPolicy(),
					WithTimeDuration(3),
				)
			}

			c := &client{
				httpClient: &http.Client{
					Transport: tt.mockTransport,
					Timeout:   time.Duration(r.TimeDuration) + 5*time.Second,
				},
				retryRoundOptions: r,
			}

			if _, err := c.Post(context.Background(), tt.url, tt.body, tt.header...); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
