package dapr

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

func (c *daprClient) HttpGet(ctx context.Context, url string) *HttpResponse {
	httpUrl := c.getFullUrl(url)
	resp, err := c.httpClient.Get(httpUrl)
	if err != nil {
		return NewHttpResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewHttpResponse(nil, errors.New(string(bs)))
	}
	return NewHttpResponse(bs, err)
}

func (c *daprClient) HttpPost(ctx context.Context, url string, reqData interface{}) *HttpResponse {
	httpUrl := c.getFullUrl(url)
	jsonData, err := json.Marshal(reqData)
	resp, err := c.httpClient.Post(httpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return NewHttpResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewHttpResponse(nil, errors.New(string(bs)))
	}
	return NewHttpResponse(bs, err)
}

func (c *daprClient) HttpPut(ctx context.Context, url string, reqData interface{}) *HttpResponse {
	httpUrl := c.getFullUrl(url)
	jsonData, err := json.Marshal(reqData)
	resp, err := c.httpClient.Post(httpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return NewHttpResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewHttpResponse(nil, errors.New(string(bs)))
	}
	return NewHttpResponse(bs, err)
}

func (c *daprClient) HttpDelete(ctx context.Context, url string, reqData interface{}) *HttpResponse {
	httpUrl := c.getFullUrl(url)
	jsonData, err := json.Marshal(reqData)
	resp, err := c.httpClient.Post(httpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return NewHttpResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewHttpResponse(nil, errors.New(string(bs)))
	}
	return NewHttpResponse(bs, err)
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
