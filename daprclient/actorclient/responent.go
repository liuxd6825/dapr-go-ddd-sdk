package actorclient

import (
	"errors"
	"net/http"
)

type ResponseStatus = int

const (
	ResponseStatusSuccess ResponseStatus = http.StatusOK
	ResponseStatusError   ResponseStatus = http.StatusInternalServerError
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

func NewResponse(err error) *Response {
	resp := &Response{
		Status:  ResponseStatusSuccess,
		Message: "success",
	}
	if err != nil {
		resp.Status = ResponseStatusError
		resp.Message = err.Error()
	}
	return resp
}

func (a *Response) GetError() error {
	if a.Status != ResponseStatusSuccess {
		return errors.New(a.Message)
	}
	return nil
}
