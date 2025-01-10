package file_handler

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
	"strings"
)

type Handler struct {
	fs        afero.Fs
	pageCache *types.CMap[bool]
	app       *iris.Application
	vdata     map[string]any
}

func NewHandler(srcFs afero.Fs, app *iris.Application, data map[string]any) *Handler {
	f := &Handler{
		app:       app,
		fs:        srcFs,
		pageCache: types.NewCMap[bool](),
		vdata:     data,
	}
	ctx := context.Background()
	if err := f.preloadDynamicPages(ctx, "/"); err != nil {
		panic(err)
	}
	// 初始化 Pongo2 模板引擎，使用 Afero 文件系统
	engine := NewEngine(srcFs, ".html")
	// 注册模板引擎到 Iris
	app.RegisterView(engine)
	return f
}

func (h *Handler) Handle(ictx iris.Context) {
	var err error
	ctx := context.Background()

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

	// 检查目录中存在文件
	var isFound bool
	isFound, err = IsFileExist(h.fs, fileName)
	if err != nil {
		return
	}

	if !isFound {
		ictx.StatusCode(iris.StatusNotFound)
		_, err = ictx.WriteString("404 Not Found")
		return
	}

	// 不设置 Content-Type，浏览器将自动推断 MIME 类型
	setContentType(ictx, fileName)

	isRender := false
	isHtml := strings.HasSuffix(fileName, ".html")
	if isHtml {
		if isRender, err = h.renderFile(ctx, ictx, fileName); err != nil {
			ictx.StatusCode(iris.StatusInternalServerError)
			ictx.SetErr(err)
			return
		}

	}
	if !isRender {
		err = h.writeFile(ictx, fileName)
	}
	if err != nil {
		ictx.StatusCode(iris.StatusInternalServerError)
		ictx.SetErr(err)
	}

}

func (h *Handler) renderFile(ctx context.Context, ictx iris.Context, fileName string) (bool, error) {
	isRender, err := isDynamicPage(ctx, h.fs, fileName)
	if isRender {
		// 动态渲染模板
		err = ictx.View(fileName, h.vdata)
		return true, err
	}
	return false, nil
}

func (h *Handler) writeFile(ictx iris.Context, fileName string) error {
	// 动态模板判断
	/*
		var content []byte
		content, err := afero.ReadFile(h.fs, fileName)
		if err != nil {
			ictx.StatusCode(iris.StatusInternalServerError)
			_, err = ictx.WriteString("Error reading template file")
			return err
		}
		_, err = ictx.Write(content)
	*/

	file, err := h.fs.Open(fileName)
	if err != nil {
		ictx.StatusCode(iris.StatusInternalServerError)
		_, err = ictx.WriteString("Error reading template file")
		return err
	}
	defer file.Close()
	// 将文件流写入响应
	_, err = io.Copy(ictx, file)
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
func (h *Handler) preloadDynamicPages(ctx context.Context, path string) error {
	files, _ := afero.ReadDir(h.fs, path)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			fileName := path + file.Name()
			isDynamic, err := isDynamicPage(ctx, h.fs, fileName)
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
	case ".min.js", ".js":
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
