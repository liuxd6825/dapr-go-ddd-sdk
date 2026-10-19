package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
}

type ClientOptions struct {
	Timeout *time.Duration
	Header  map[string][]string
}

func (o *ClientOptions) SetTimeout(timeout time.Duration) *ClientOptions {
	o.Timeout = &timeout
	return o
}

func (o *ClientOptions) AddHeader(header string, values []string) *ClientOptions {
	o.Header[header] = values
	return o
}

func NewClientOptions(opts ...*ClientOptions) *ClientOptions {
	o := &ClientOptions{
		Header: make(map[string][]string),
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Timeout != nil {
			o.Timeout = opt.Timeout
		}
		if opt.Header != nil {
			for k, v := range opt.Header {
				o.Header[k] = v
			}
		}
	}
	return o
}

func NewClient() *Client {
	return &Client{}
}

func (s *Client) newClient(opt *ClientOptions) *http.Client {
	client := &http.Client{}
	if opt != nil {
		if opt.Timeout != nil {
			client.Timeout = *opt.Timeout
		}
	}
	return client
}

func (s *Client) Delete(url string, request any, response any, opts ...*ClientOptions) (*http.Response, error) {
	return s.Do("DELETE", url, request, response, opts...)
}

func (s *Client) Get(url string, request any, response any, opts ...*ClientOptions) (*http.Response, error) {
	return s.Do("GET", url, request, response, opts...)
}

func (s *Client) Post(url string, request any, response any, opts ...*ClientOptions) (*http.Response, error) {
	return s.Do("POST", url, request, response, opts...)
}

func (s *Client) Do(methodType string, url string, request any, response any, opts ...*ClientOptions) (*http.Response, error) {
	opt := NewClientOptions(opts...)
	// 2. 创建自定义的 http.Client（建议设置超时时间）
	client := s.newClient(opt)

	// 如果需要发送请求体，可以将最后一个参数替换为 bytes.NewBuffer(jsonData)
	var reader io.Reader
	if request != nil {
		reqJson, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(reqJson)
	}

	req, err := http.NewRequest(methodType, url, reader)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
	}
	if opt.Header != nil {
		req.Header = opt.Header
	}

	// 5. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return resp, err
	}
	// 务必记得关闭 Body
	defer func() {
		_ = resp.Body.Close()
	}()

	// 6. 读取并打印响应状态码与内容
	fmt.Printf("响应状态码: %s\n", resp.Status)

	if response != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("读取响应体失败: %v\n", err)
		}
		err = json.Unmarshal(body, response)
		if err != nil {
			return resp, err
		}
	}
	return resp, nil
}
