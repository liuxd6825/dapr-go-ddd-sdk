package daprclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dapr_sdk_client "github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultMaxIdleConns        = 10
	DefaultMaxIdleConnsPerHost = 50
	DefaultIdleConnTimeout     = 5
)

type DaprClient interface {
	HttpGet(ctx context.Context, url string) *Response
	HttpPost(ctx context.Context, url string, reqData interface{}) *Response
	HttpPut(ctx context.Context, url string, reqData interface{}) *Response
	InvokeService(ctx context.Context, appID, methodName, verb string, request interface{}, response interface{}) (interface{}, error)
}

type daprClient struct {
	host       string
	httpPort   int
	grpcPort   int
	client     *http.Client
	grpcClient dapr_sdk_client.Client
}

type DaprHttpOptions struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int
}

type Option func(options *DaprHttpOptions)

func newHttpOptions() *DaprHttpOptions {
	options := &DaprHttpOptions{
		MaxIdleConns:        DefaultMaxIdleConns,
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:     DefaultIdleConnTimeout,
	}
	return options
}

func MaxIdleConns(val int) Option {
	return func(options *DaprHttpOptions) {
		options.MaxIdleConns = val
	}
}

func MaxIdleConnsPerHost(val int) Option {
	return func(options *DaprHttpOptions) {
		options.MaxIdleConnsPerHost = val
	}
}

func IdleConnTimeout(val int) Option {
	return func(options *DaprHttpOptions) {
		options.IdleConnTimeout = val
	}
}

func NewClient(host string, httpPort int, grpcPort int, opts ...Option) (DaprClient, error) {
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

	grpcClient, err := dapr_sdk_client.NewClientWithPort(strconv.Itoa(grpcPort))
	if err != nil {
		return nil, err
	}

	return &daprClient{
		client:     client,
		host:       host,
		httpPort:   httpPort,
		grpcPort:   grpcPort,
		grpcClient: grpcClient,
	}, nil

}

func (c *daprClient) HttpGet(ctx context.Context, url string) *Response {
	getUlr := fmt.Sprintf("http://%s:%d/%s", c.host, c.httpPort, url)
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

func (c *daprClient) HttpPost(ctx context.Context, url string, reqData interface{}) *Response {
	httpUrl := fmt.Sprintf("http://%s:%d/%s", c.host, c.httpPort, url)
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

func (c *daprClient) HttpPut(ctx context.Context, url string, reqData interface{}) *Response {
	httpUrl := fmt.Sprintf("http://%s:%d/%s", c.host, c.httpPort, url)
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

func (c *daprClient) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	return bs, err
}

func (c *daprClient) InvokeService(ctx context.Context, appID, methodName, verb string, request interface{}, response interface{}) (interface{}, error) {
	var err error
	defer func() {
		if e := ddd_errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()
	var respBytes []byte

	if request != nil {
		reqBytes, err := json.Marshal(request)
		if err != nil {
			return nil, newAppError(appID, err)
		}
		content := &dapr_sdk_client.DataContent{
			ContentType: "application/json",
			Data:        reqBytes,
		}
		respBytes, err = c.grpcClient.InvokeMethodWithContent(ctx, appID, methodName, verb, content)
	} else {
		respBytes, err = c.grpcClient.InvokeMethod(ctx, appID, methodName, verb)
	}
	if err != nil {
		return nil, newAppError(appID, err)
	}
	if len(respBytes) > 0 {
		err = json.Unmarshal(respBytes, response)
		if err != nil {
			return nil, newAppError(appID, err)
		}
		return response, nil
	}
	return nil, nil
}

func newAppError(appID string, err error) error {
	msg := fmt.Sprintf("AppId is %s , %s", appID, err.Error())
	return errors.New(msg)
}
