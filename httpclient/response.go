package httpclient

import "encoding/json"

type Response struct {
	data []byte
	err  error
}

func NewResponse(data []byte, err error) *Response {
	return &Response{
		data: data,
		err:  err,
	}
}

func (r *Response) GetData() []byte {
	return r.data
}

func (r *Response) GetError() error {
	return r.err
}

func (r *Response) OnSuccess(data interface{}, fun func() error) *Response {
	if r.data == nil {
		return r
	}
	if r.err != nil {
		return r
	}
	err := json.Unmarshal(r.data, data)
	if err != nil {
		r.err = err
	} else {
		r.err = fun()
	}
	return r
}

func (r *Response) OnError(fun func(err error)) *Response {
	if r.err != nil {
		fun(r.err)
	}
	return r
}
