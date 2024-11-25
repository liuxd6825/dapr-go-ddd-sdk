package hserver

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"strings"
)

// Service 定义服务的基本结构
type Service struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	server      *Server `json:"-"`
	requests    cmap.ConcurrentMap
	data        cmap.ConcurrentMap
	srcFileName string
	srcPath     string
	code        string
	codeType    string
	scripts     *ScriptManager
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
		scripts:     NewScriptManager(server.logger),
	}

	if err := service.parse(html); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Service) InitVM(vm *goja.Runtime) error {
	return nil
}

func (s *Service) AddRequest(r *Request) {
	s.requests.Set(r.opts.Name, r)
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
	serviceEl := doc.Find("body service")
	if serviceEl == nil {
		return errors.New("no service found")
	}
	s.URL = serviceEl.AttrOr("url", "")
	s.Name = serviceEl.AttrOr("mame", "") // 注意拼写错误 "mame"
	s.Description = serviceEl.AttrOr("description", "")
	optsList := []RequestOptions{}
	serviceEl.Children().Each(func(i int, child *goquery.Selection) {
		node := child.Get(0) // 获取当前节点的 *html.Node
		if node.Type != 3 {
			return
		}
		switch node.Data {
		case "script":
			s.code = child.Text()
			s.codeType = child.AttrOr("type", "")
		case "request":
			// 解析 Request 列表
			opts := RequestOptions{
				Type:        strings.ToUpper(child.AttrOr("type", "")),
				Name:        child.AttrOr("name", ""),
				URL:         child.AttrOr("url", ""),
				AbsURl:      child.AttrOr("abs-url", ""),
				Description: child.AttrOr("description", ""),
				ParamsUrl:   child.Find("params").AttrOr("url", ""),
			}
			// 获取 <script> 子节点
			scriptEl := child.Find("script")
			if scriptEl != nil {
				opts.Code = scriptEl.Text()
				opts.CodeType = scriptEl.AttrOr("type", "")
			}
			optsList = append(optsList, opts)
		}

	})

	for _, opts := range optsList {
		request, err := NewRequest(s.server, s, opts)
		if err != nil {
			return err
		}
		s.AddRequest(request)
	}
	return err
}

func (s *Service) Run() error {
	s.runInitScript()
	for _, key := range s.GetRequestKeys() {
		r := s.GetRequest(key)
		if r != nil {
			err := r.Run()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) runInitScript() error {
	if s.code != "" {
		fileName := fmt.Sprintf("service_%s_init.js", s.Name)
		opts := &RuntimeOption{
			Server:  s.server,
			Service: s,
		}
		err := RunInitScript(s.scripts, s.code, s.codeType, fileName, opts, s.GetLogger())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetLogger() logrus.FieldLogger {
	return s.server.logger
}

// getFilePath 根据文件名获取文件路径
func getFilePath(fileName string) (string, error) {
	i := strings.LastIndex(fileName, "/")
	if i == -1 {
		return "", errors.New("file name not exist")
	}
	return fileName[:i], nil
}
