package light_rag

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type Client struct {
	httpClient *http.Client
	doc        *DocumentService
	query      *QueryService
	graph      *GraphService
	baseUrl    string
}

func NewClient(baseUrl string, httpClient *http.Client) *Client {
	cli := &Client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
	cli.doc = NewDocumentService(cli)
	cli.query = NewQueryService(cli)
	return cli
}

func (c *Client) Document() *DocumentService {
	return c.doc
}

func (c *Client) Query() *QueryService {
	return c.query
}

func (c *Client) Graph() *GraphService {
	return c.graph
}

func (c *Client) Post(ctx context.Context, url string, data any) (resp *http.Response, err error) {
	reader := c.newReader(ctx, data)
	request := c.newRequest(ctx, http.MethodPost, url, reader)
	return c.httpClient.Do(request)
}

func (c *Client) Put(ctx context.Context, url string, data any) (resp *http.Response, err error) {
	reader := c.newReader(ctx, data)
	request := c.newRequest(ctx, http.MethodPut, url, reader)
	return c.httpClient.Do(request)
}

func (c *Client) Delete(ctx context.Context, url string) (resp *http.Response, err error) {
	request := c.newRequest(ctx, http.MethodDelete, url, nil)
	return c.httpClient.Do(request)
}

func (c *Client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	request := c.newRequest(ctx, http.MethodGet, url, nil)
	return c.httpClient.Do(request)
}

func (c *Client) newRequest(ctx context.Context, methodType string, url string, body io.Reader) *http.Request {
	request := &http.Request{
		Method: methodType,
		URL:    c.newUrl(ctx, url),
		Header: http.Header{},
		Body:   NewReadCloser(body),
	}
	return request
}

func (c *Client) newUrl(ctx context.Context, aUrl string) *url.URL {
	rawURL := c.baseUrl + aUrl
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return parsedURL
}

func (c *Client) newReader(ctx context.Context, data any) *strings.Reader {
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	reader := strings.NewReader(string(d))
	return reader
}

func (c *Client) getJsonData(ctx context.Context, resp *http.Response, data any) error {
	if resp == nil || resp.Body == nil {
		return errors.New("http response or body is nil")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(bytes, data)
}

func (c *Client) uploadFile(ctx context.Context, url string, fileName string, fs afero.Fs) error {
	// 打开文件
	fileData, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	fileReader := bytes.NewReader(fileData)
	// 创建一个缓冲区用于存储 multipart 数据
	body := &bytes.Buffer{}
	// 创建一个新的 multipart writer
	writer := multipart.NewWriter(body)

	// 创建一个文件字段并将文件数据写入其中
	part, err := writer.CreateFormFile("file", filepath.Base(fileName))
	if err != nil {
		return fmt.Errorf("unable to create form file: %v", err)
	}

	// 将文件内容复制到表单字段中
	_, err = io.Copy(part, fileReader)
	if err != nil {
		return fmt.Errorf("unable to copy file contents: %v", err)
	}

	// 结束 multipart 写入
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("unable to close writer: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, c.baseUrl+url, body)
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}

	// 设置 Content-Type 头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %v", resp.Status)
	}

	// 成功上传
	fmt.Println("File uploaded successfully!")
	return nil
}
