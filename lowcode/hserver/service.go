package hserver

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/orcaman/concurrent-map"
	"strings"
)

// Service 定义服务的基本结构
type Service struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	Code        string  `json:"code"`
	server      *Server `json:"-"`
	requests    cmap.ConcurrentMap
	data        cmap.ConcurrentMap
	srcFileName string
	srcPath     string
}

// Call 定义调用信息结构
type Call struct {
	Ref    string            `json:"ref"`
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

func NewService(server *Server, html []byte, srcFileName string, data map[string]any) (*Service, error) {

	srcPath, err := getFilePath(srcFileName)
	if err != nil {
		return nil, err
	}

	dataMap := cmap.New()
	for k, v := range data {
		dataMap.Set(k, v)
	}

	// 解析 Service 节点
	service := &Service{
		server:      server,
		data:        dataMap,
		srcFileName: srcFileName,
		srcPath:     srcPath,
		requests:    cmap.New(),
	}

	if err := service.parse(html); err != nil {
		return nil, err
	}

	if err := service.init(); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Service) AddRequest(r *Request) {
	s.requests.Set(r.Name, r)
}

func (s *Service) GetRequest(name string) *Request {
	r, ok := s.requests.Get(name)
	if !ok {
		return nil
	}
	return r.(*Request)
}

func (s *Service) GetRequestKeys() []string {
	return s.requests.Keys()
}

func (s *Service) GetRequestCount() int {
	return s.requests.Count()
}

func (s *Service) parse(html []byte) error {
	reader := bytes.NewReader(html)
	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}

	s.URL = doc.Find("service").AttrOr("url", "")
	s.Name = doc.Find("service").AttrOr("mame", "") // 注意拼写错误 "mame"
	s.Description = doc.Find("service").AttrOr("description", "")

	code := doc.Find("service > script").Text()
	if code != "" {
		s.Code = code
	}

	// 解析 Request 列表
	doc.Find("service > request").Each(func(i int, sel *goquery.Selection) {
		request := &Request{
			Type:        strings.ToUpper(sel.AttrOr("type", "")),
			Name:        sel.AttrOr("name", ""),
			URL:         sel.AttrOr("url", ""),
			AbsURl:      sel.AttrOr("abs-url", ""),
			Description: sel.AttrOr("description", ""),
			Params:      sel.Find("params").AttrOr("url", ""),
		}

		// 获取 <call> 子节点
		/*
			callSel := sel.Find("call")
			if callSel.Length() > 0 {
				call := Call{
					Ref:    callSel.AttrOr("ref", ""),
					Method: callSel.AttrOr("method", ""),
					Params: make(map[string]string),
				}
				callSel.Find("param").Each(func(i int, paramSel *goquery.Selection) {
					name := paramSel.AttrOr("name", "")
					value := paramSel.Text()
					call.Params[name] = value
				})
				request.Call = call
			}
		*/

		// 获取 <script> 子节点
		script := sel.Find("script").Text()
		if script != "" {
			request.Code = script
		}

		s.AddRequest(request)
	})

	return err
}

func (s *Service) init() error {
	if s.Code != "" {
		vm := runtime.NewRuntime()
		opts := &SetRuntimeOption{
			Server:  s.server,
			Service: s,
		}
		if err := setRuntime(vm, opts); err != nil {
			return err
		}
		if _, err := vm.RunString(s.Code); err != nil {
			return err
		}
	}
	return nil
}

// getFilePath 根据文件名获取文件路径
func getFilePath(fileName string) (string, error) {
	i := strings.LastIndex(fileName, "/")
	if i == -1 {
		return "", errors.New("file name not exist")
	}
	return fileName[:i], nil
}
