package types

type FileInfo struct {
	IsDir     bool        `json:"isDir,omitempty"`
	Name      string      `json:"name,omitempty"`
	Size      int64       `json:"size,omitempty"`
	Path      string      `json:"path,omitempty"`
	SubFiles  []*FileInfo `json:"subFiles,omitempty"`
	SizeTitle string      `json:"sizeTitle,omitempty"`
}
