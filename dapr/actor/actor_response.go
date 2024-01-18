package actor

import (
	"errors"
	"net/http"
)

type ResponseStatus = int

const (
	ResponseStatusSuccess ResponseStatus = http.StatusOK
	ResponseStatusError   ResponseStatus = http.StatusInternalServerError
)
const (
	ResponseMessageSuccess = "success"
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
	Data    any            `json:"data"`
}

func NewResponse(err error) *Response {
	resp := &Response{
		Status:  ResponseStatusSuccess,
		Message: ResponseMessageSuccess,
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

func (a *Response) SetError(err error) {
	if err != nil {
		a.Status = ResponseStatusError
		a.Message = err.Error()
	} else {
		a.Status = ResponseStatusSuccess
		a.Message = ResponseMessageSuccess
	}
}
