package file

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/spf13/afero"
	"strings"
)

type Handler struct {
	fs        afero.Fs
	pageCache *types.CMap[bool]
	app       *iris.Application
}

func NewHandler(srcFs afero.Fs, app *iris.Application) *Handler {
	f := &Handler{
		app:       app,
		fs:        srcFs,
		pageCache: types.NewCMap[bool](),
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
	}

	// 检查目录中存在文件
	var isFound bool
	isFound, err = IsFileExist(h.fs, fileName)
	if err != nil {
		return
	}

	if isFound {
		isRender := false
		isHtml := strings.HasSuffix(fileName, ".html")
		if isHtml {
			if isRender, err = h.renderFile(ctx, ictx, fileName); err != nil {
				return
			}
		}
		if !isRender {
			err = h.writeFile(ictx, fileName)
		}

	} else {
		ictx.StatusCode(iris.StatusNotFound)
		_, err = ictx.WriteString("404 Not Found")
	}
}

func (h *Handler) renderFile(ctx context.Context, ictx iris.Context, fileName string) (bool, error) {
	isRender, err := isDynamicPage(ctx, h.fs, fileName)
	if isRender {
		// 动态渲染模板
		err = ictx.View(fileName, iris.Map{
			"Title": "Dynamic Page",
			"User":  "Dynamic User",
		})
		return true, err
	}
	return false, nil
}

func (h *Handler) writeFile(ictx iris.Context, fileName string) error {
	// 动态模板判断
	var content []byte
	content, err := afero.ReadFile(h.fs, fileName)
	if err != nil {
		ictx.StatusCode(iris.StatusInternalServerError)
		_, err = ictx.WriteString("Error reading template file")
		return err
	}
	_, err = ictx.Write(content)
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
