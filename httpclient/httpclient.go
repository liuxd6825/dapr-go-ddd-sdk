package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	DefaultMaxIdleConns        = 10
	DefaultMaxIdleConnsPerHost = 50
	DefaultIdleConnTimeout     = 5
)

type HttpClient struct {
	host   string
	port   int
	client *http.Client
}

type HttpOptions struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int
}

type Option func(options *HttpOptions)

func newHttpOptions() *HttpOptions {
	options := &HttpOptions{
		MaxIdleConns:        DefaultMaxIdleConns,
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:     DefaultIdleConnTimeout,
	}
	return options
}

func MaxIdleConns(val int) Option {
	return func(options *HttpOptions) {
		options.MaxIdleConns = val
	}
}

func MaxIdleConnsPerHost(val int) Option {
	return func(options *HttpOptions) {
		options.MaxIdleConnsPerHost = val
	}
}

func IdleConnTimeout(val int) Option {
	return func(options *HttpOptions) {
		options.IdleConnTimeout = val
	}
}

func NewHttpClient(host string, port int, opts ...Option) (*HttpClient, error) {
	options := newHttpOptions()
	for _, opt := range opts {
		opt(options)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        options.MaxIdleConns,
			MaxIdleConnsPerHost: options.MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Second * time.Duration(options.IdleConnTimeout),
		},
	}

	return &HttpClient{
		client: client,
		host:   host,
		port:   port,
	}, nil

}

func (c *HttpClient) Get(ctx context.Context, url string) *Response {
	getUlr := fmt.Sprintf("http://%s:%d/%s", c.host, c.port, url)
	resp, err := c.client.Get(getUlr)
	if err != nil {
		return NewResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewResponse(nil, errors.New(string(bs)))
	}
	return NewResponse(bs, err)
}

func (c *HttpClient) Post(ctx context.Context, url string, reqData interface{}) *Response {
	httpUrl := fmt.Sprintf("http://%s:%d/%s", c.host, c.port, url)
	jsonData, err := json.Marshal(reqData)
	resp, err := c.client.Post(httpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return NewResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewResponse(nil, errors.New(string(bs)))
	}
	return NewResponse(bs, err)
}

func (c *HttpClient) Put(ctx context.Context, url string, reqData interface{}) *Response {
	httpUrl := fmt.Sprintf("http://%s:%d/%s", c.host, c.port, url)
	jsonData, err := json.Marshal(reqData)
	resp, err := c.client.Post(httpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return NewResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewResponse(nil, errors.New(string(bs)))
	}
	return NewResponse(bs, err)
}

func (c *HttpClient) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	return bs, err
}
