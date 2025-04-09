package env

// HServer
// @Description: 脚本服务配置
// @Author:       liuxd
// @Date:         2021/10/18 10:57
type HServer struct {
	Enable       bool           `yaml:"enable" json:"enable"`             // 是否启用脚本服务
	SrcName      string         `yaml:"srcName" json:"srcName"`           // API源码文件系统名称
	WebName      string         `yaml:"webName" json:"webName"`           // Web源码文件系统名称
	BasePath     string         `yaml:"basePath" json:"basePath"`         // 脚本文件路径
	Reload       bool           `yaml:"reload" json:"reload"`             // 是否自动加载脚本
	WatchRestart bool           `yaml:"watchRestart" json:"watchRestart"` // 检查文件变化，重新启动
	Metadata     map[string]any `yaml:"metadata" json:"metadata"`
	Npm          Npm            `yaml:"npm" json:"npm"`
}

type Npm struct {
	Links []*NpmLink `yaml:"links" json:"links"`
}

type NpmLink struct {
	Name string `yaml:"name" json:"name"`
	Path string `yaml:"path" json:"path"`
}
