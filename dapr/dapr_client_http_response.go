package dapr

import "encoding/json"

type HttpResponse struct {
	data []byte
	err  error
}

func NewHttpResponse(data []byte, err error) *HttpResponse {
	return &HttpResponse{
		data: data,
		err:  err,
	}
}

func (r *HttpResponse) GetData() []byte {
	return r.data
}

func (r *HttpResponse) GetError() error {
	return r.err
}

func (r *HttpResponse) OnSuccess(data interface{}, fun func() error) *HttpResponse {
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

func (r *HttpResponse) OnError(fun func(err error)) *HttpResponse {
	if r.err != nil {
		fun(r.err)
	}
	return r
}
