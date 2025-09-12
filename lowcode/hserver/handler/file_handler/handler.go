package file_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
	engine2 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/handler/file_handler/engine"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/spf13/afero"
	"path/filepath"
	"strings"
)

type Config struct {
	Env         *env.Env
	ServerFs    afero.Fs
	WebFs       afero.Fs
	NodeModules []afero.Fs
}

type Handler struct {
	cfg          *Config             // 配置
	pageCache    *types.CMap[bool]   // 页面缓存
	app          *iris.Application   // iris APP对象
	globalValues map[string]any      // 全局view模板数据
	prodMode     bool                // 是生产模式
	webCfg       *WebConfig          //
	htmlJs       *types.CMap[string] // 将html中的<script>import * </script>转为js代码文件的内容
}

func NewHandler(app *iris.Application, globalValues map[string]any, cfg *Config) *Handler {
	f := &Handler{
		app:          app,
		cfg:          cfg,
		pageCache:    types.NewCMap[bool](),
		globalValues: globalValues,
		prodMode:     env.GetEnv().GetProdMode(),
		htmlJs:       types.NewCMap[string](),
	}
	ctx := context.Background()
	if err := f.preloadDynamicPages(ctx, nil, "/"); err != nil {
		panic(err)
	}

	app.Use(f.PathInterceptor)

	// 初始化 Pongo2 模板引擎，使用 Afero 文件系统
	en := engine2.NewEngine(cfg.Env, cfg.ServerFs, cfg.WebFs, ".html")

	// 注册模板引擎到 Iris
	app.RegisterView(en)
	return f
}

// PathInterceptor 路径拦截器（支持.html.js和普通路径）
func (h *Handler) PathInterceptor(ctx iris.Context) {
	// 获取原始路径（如 /human/subdir/bill.html.js）
	rawPath := ctx.Path()

	// 仅拦截特定后缀的请求
	if strings.HasSuffix(rawPath, ".html.js") {
		ctx.StatusCode(200)
		ctx.ContentType("application/javascript; charset=utf-8")
		str, ok := h.htmlJs.Get(rawPath)
		if ok {
			_, _ = ctx.WriteString(str)
		}
		return
	}
	ctx.Next() // 继续执行后续中间件/路由
}

func newWebConfig(fs afero.Fs, fileName string) *WebConfig {
	exists, err := afero.Exists(fs, fileName)
	if err != nil {
		panic(err)
	}
	if !exists {
		return &WebConfig{Npm: &Npm{Links: make([]*NpmLink, 0)}}
	}
	data, err := afero.ReadFile(fs, fileName)
	if err != nil {
		panic(err)
	}
	cfg := &WebConfig{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (h *Handler) Handle(ictx iris.Context) {
	var err error
	defer func() {
		err = utils.RecoverError(err, recover())
		if err != nil {
			utils.SetError(ictx, err)
		}
	}()

	// 获取文件路径
	fileName := "/" + ictx.Params().Get("file")
	if fileName == "" || fileName == "/" {
		fileName = "/index.html" // 默认页面
	} else if !strings.Contains(fileName, ".") {
		fileName = fileName + ".html"
	}
	ctx := appctx.NewWebContext(context.Background(), ictx)
	logs.Infofmt(ctx, "file handle %s", fileName)
	_ = h.RenderView(ctx, ictx, fileName)
}

func (h *Handler) RenderView(ctx context.Context, ictx iris.Context, fileName string, viewData ...map[string]any) error {
	if fileName == "" {
		fileName = ictx.Request().URL.Path
	}
	// 检查目录中存在文件
	fs, _, isFound, err := h.isFileExist(fileName)
	if err != nil {
		return err
	}
	// 对js和ts文件进行转换
	if !isFound {
		fs, fileName, isFound, err = h.GetJsTsFile(fileName)
	}
	if err != nil {
		return err
	}

	if !isFound {
		ictx.StatusCode(iris.StatusNotFound)
		_, err = ictx.WriteString("404 Not Found")
		return err
	}

	// 不设置 Content-Type，浏览器将自动推断 MIME 类型
	setContentType(ictx, fileName)

	err = h.render(ctx, ictx, fs, fileName, viewData...)
	if err != nil {
		ictx.StatusCode(iris.StatusInternalServerError)
		ictx.SetErr(err)
	}
	return err
}

func (h *Handler) isFileExist(fileName string) (fs afero.Fs, resFileName string, exist bool, err error) {
	resFileName = fileName
	exist, err = afero.Exists(h.cfg.WebFs, fileName)
	if exist {
		fs = h.cfg.WebFs
		return
	}
	for _, f := range h.cfg.NodeModules {
		exist, err = afero.Exists(f, fileName)
		if exist {
			fs = f
			return
		}
	}
	return
}

func (h *Handler) GetJsTsFile(fileName string) (fs afero.Fs, resFileName string, exist bool, err error) {
	if strings.HasSuffix(fileName, ".js") {
		fileName = fileName[:len(fileName)-len(".js")] + ".ts"
		fs, resFileName, exist, err = h.isFileExist(fileName)
	} else if strings.HasSuffix(fileName, ".ts") {
		fileName = fileName[:len(fileName)-len(".ts")] + ".js"
		fs, resFileName, exist, err = h.isFileExist(fileName)
	}
	return fs, resFileName, exist, err
}

// render
// @Description: 渲染HTML
// @receiver h
// @param ctx 上下文
// @param ictx Web上下文
// @param fs 文件系统
// @param fileName 文件名称
// @param viewData 数据
// @return error
func (h *Handler) render(ctx context.Context, ictx iris.Context, fs afero.Fs, fileName string, viewData ...map[string]any) error {
	var err error
	isRender := false
	isHtml := strings.HasSuffix(fileName, ".html")
	if isHtml {
		fsData, err := afero.ReadFile(fs, fileName)
		if err != nil {
			return err
		}
		if !h.prodMode {
			reader := bytes.NewReader(fsData)
			doc, err := goquery.NewDocumentFromReader(reader)
			if err != nil {
				return err
			}
			fsData, err = parserHtml(doc, &env.GetEnv().App.HServer.Npm, func(scripts *goquery.Selection, sb *strings.Builder) {
				if sb != nil {
					scripts.Remove()
					// 插入方式1：插入到#container末尾
					doc.Find("body").AppendHtml(fmt.Sprintf(`<script type="module" src="%s"></script>`, fileName+".js"))
					h.htmlJs.Set(fileName+".js", sb.String())
				}
			})
			if err != nil {
				return err
			}
			doc.Find("meta").Each(func(i int, s *goquery.Selection) {
				name, exists := s.Attr("name")
				if exists && name == "ssr" {
					isRender = true
				}
			})
		}
		if !isRender {
			isRender, err = h.isDynamicPage(ctx, ictx, fs, fileName, h.prodMode)
			if err != nil {
				return err
			}
		}

		if isRender {
			vData := make(map[string]interface{})
			for k, v := range h.globalValues {
				vData[k] = v
			}
			for _, dataMap := range viewData {
				for k, v := range dataMap {
					vData[k] = v
				}
			}
			vData["ctx"] = ctx
			if fsData != nil {
				vData["templateData"] = fsData
			}
			err = ictx.View(fileName, vData)
			return err
		}
	}

	if isRender == false {
		err = h.writeFile(ictx, fs, fileName)
	}
	return err
}

func (h *Handler) writeFile(ictx iris.Context, fs afero.Fs, fileName string) error {
	file, err := fs.Open(fileName)
	if err != nil {
		err = errors.New("Error reading template file %s ", fileName)
		return err
	}

	//defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if strings.HasSuffix(fileName, ".ts") {
		tsData, err := afero.ReadFile(fs, fileName)
		if err != nil {
			return err
		}
		jsData, err := CompileTS(tsData)
		if err != nil {
			return err
		}
		reader := bytes.NewReader(jsData) // 直接返回 io.ReadSeeker
		ictx.ServeContentWithRate(reader, fileName, info.ModTime(), 0, 0)
		return err
	}

	// 将文件流写入响应
	//_, err = io.Copy(ictx, file)
	ictx.ServeContentWithRate(file, fileName, info.ModTime(), 0, 0)
	return err
}

// preloadDynamicPages
//
//	@Description: 缓存加载动态HTML文件
//	@receiver h
//	@param ctx
//	@param fs
//	@param path
//	@return error
func (h *Handler) preloadDynamicPages(ctx context.Context, ictx iris.Context, path string) error {
	srcFs := h.cfg.WebFs
	files, _ := afero.ReadDir(srcFs, path)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			fileName := path + file.Name()
			isDynamic, err := h.isDynamicPage(ctx, ictx, srcFs, fileName, env.GetEnv().GetProdMode())
			if err != nil {
				return err
			}
			h.pageCache.Set(fileName, isDynamic)
		}
	}
	return nil
}

// 获取文件扩展名并设置Content-Type
func setContentType(ctx iris.Context, filename string) {
	// 获取文件扩展名（不区分大小写）
	ext := strings.ToLower(filepath.Ext(filename))

	// 根据扩展名设置Content-Type
	switch ext {
	case ".html":
		ctx.ContentType("text/html; charset=utf-8")
	case ".css":
		ctx.ContentType("text/css; charset=utf-8")
	case ".js":
		ctx.ContentType("application/javascript; charset=utf-8")
	case ".ts":
		ctx.ContentType("application/javascript; charset=utf-8")
	case ".jpg", ".jpeg":
		ctx.ContentType("image/jpeg; charset=utf-8")
	case ".png":
		ctx.ContentType("image/png; charset=utf-8")
	case ".gif":
		ctx.ContentType("image/gif; charset=utf-8")
	case ".json":
		ctx.ContentType("application/json; charset=utf-8")
	case ".xml":
		ctx.ContentType("application/xml; charset=utf-8")
	default:
		ctx.ContentType("application/octet-stream; charset=utf-8") // 默认二进制流
	}
}
