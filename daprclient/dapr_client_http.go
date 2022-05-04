package daprclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *daprClient) HttpGet(ctx context.Context, url string) *Response {
	httpUrl := c.getFullUrl(url)
	resp, err := c.client.Get(httpUrl)
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
	httpUrl := c.getFullUrl(url)
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
	httpUrl := c.getFullUrl(url)
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

func (c *daprClient) getFullUrl(url string) string {
	res := fmt.Sprintf("%s://%s:%d", Protocol, c.host, c.httpPort)
	if strings.HasPrefix(url, "/") {
		return fmt.Sprintf("%s%s", res, url)
	}
	return fmt.Sprintf("%s/%s", res, url)
}

func (c *daprClient) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	return bs, err
}
