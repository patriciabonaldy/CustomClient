package jsend

import (
	"encoding/json"
	"net/http"
)

// Status constants
const (
	Error   Status = "error"
	Fail    Status = "fail"
	Success Status = "success"
)

type Status string

// response contains
type response struct {
	Status  Status      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// Ok returns a success response with the given data.
func Ok(w http.ResponseWriter, message string, data interface{}) {
	withSuccess(w, message, http.StatusOK, data)
}

// Created returns a success response with the given data.
func Created(w http.ResponseWriter, message string, data interface{}) {
	withSuccess(w, message, http.StatusCreated, data)
}

// NotContent returns a success response with the given data.
func NotContent(w http.ResponseWriter, message string, data interface{}) {
	withSuccess(w, message, http.StatusNoContent, data)
}

// InternalServerError returns a error response with given message.
func InternalServerError(w http.ResponseWriter, message string, data interface{}) {
	withError(w, message, http.StatusInternalServerError, data)
}

// NotFound returns a error response with given message.
func NotFound(w http.ResponseWriter, message string, data interface{}) {
	withFail(w, message, http.StatusNotFound, data)
}

func withSuccess(w http.ResponseWriter, message string, statusCode int, data interface{}) {
	body := response{
		Status:  Success,
		Message: message,
		Data:    data,
	}
	write(w, body, statusCode)
}

func withFail(w http.ResponseWriter, message string, statusCode int, data interface{}) {
	body := response{
		Status:  Fail,
		Message: message,
		Code:    statusCode,
		Data:    data,
	}

	write(w, body, statusCode)
}

func withError(w http.ResponseWriter, message string, statusCode int, data interface{}) {
	body := response{
		Status:  Error,
		Message: message,
		Code:    statusCode,
		Data:    data,
	}
	write(w, body, statusCode)
}

// write writes the response to http.ResponseWriter.
func write(w http.ResponseWriter, body response, statuses ...int) error {
	w.Header().Set("Content-Type", "application/json")

	if len(statuses) > 0 {
		w.WriteHeader(statuses[0])
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
