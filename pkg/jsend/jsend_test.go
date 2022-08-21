package jsend

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOk(t *testing.T) {
	type data struct {
		Name string `json:"name"`
	}
	type args struct {
		message      string
		data         data
		expectedBody response
	}

	params := args{
		data: data{
			Name: "Juan",
		},
		message: "ok",
		expectedBody: response{
			Status: Success,
			Data: map[string]interface{}{
				"name": "Juan",
			},
			Message: "ok",
		},
	}
	w := httptest.NewRecorder()
	Ok(w, params.message, params.data)

	gotContentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", gotContentType)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var gotBody response
	if err := json.Unmarshal(w.Body.Bytes(), &gotBody); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, params.expectedBody, gotBody)
}

func TestCreated(t *testing.T) {
	type data struct {
		Name string `json:"name"`
	}
	type args struct {
		message      string
		data         data
		expectedBody response
	}

	params := args{
		data: data{
			Name: "Juan",
		},
		message: "ok",
		expectedBody: response{
			Status: Success,
			Data: map[string]interface{}{
				"name": "Juan",
			},
			Message: "ok",
		},
	}
	w := httptest.NewRecorder()
	Created(w, params.message, params.data)

	gotContentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", gotContentType)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	var gotBody response
	if err := json.Unmarshal(w.Body.Bytes(), &gotBody); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, params.expectedBody, gotBody)
}

func TestNotContent(t *testing.T) {
	type data struct {
		Name string `json:"name"`
	}
	type args struct {
		message      string
		data         data
		expectedBody response
	}

	params := args{
		data: data{
			Name: "Juan",
		},
		message: "ok",
		expectedBody: response{
			Status: Success,
			Data: map[string]interface{}{
				"name": "Juan",
			},
			Message: "ok",
		},
	}
	w := httptest.NewRecorder()
	NotContent(w, params.message, params.data)

	gotContentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", gotContentType)
	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)

	var gotBody response
	if err := json.Unmarshal(w.Body.Bytes(), &gotBody); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, params.expectedBody, gotBody)
}

func TestNotFound(t *testing.T) {
	type args struct {
		message      string
		expectedBody response
	}

	params := args{
		message: "not found",
		expectedBody: response{
			Status:  Fail,
			Message: "not found",
			Code:    http.StatusNotFound,
		},
	}
	w := httptest.NewRecorder()
	NotFound(w, params.message, nil)

	gotContentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", gotContentType)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	var gotBody response
	if err := json.Unmarshal(w.Body.Bytes(), &gotBody); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, params.expectedBody, gotBody)
}

func TestInternalError(t *testing.T) {
	type args struct {
		message      string
		expectedBody response
	}

	params := args{
		message: "error trying do something...",
		expectedBody: response{
			Status:  Error,
			Message: "error trying do something...",
			Code:    http.StatusInternalServerError,
		},
	}
	w := httptest.NewRecorder()
	InternalServerError(w, params.message, nil)

	gotContentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", gotContentType)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	var gotBody response
	if err := json.Unmarshal(w.Body.Bytes(), &gotBody); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, params.expectedBody, gotBody)
}
