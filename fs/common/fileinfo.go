package common

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	name    string      // base name of the file
	size    int64       // length in bytes for regular files; system-dependent for others
	mode    fs.FileMode // file mode bits
	modTime time.Time   // modification time
	isDir   bool        // abbreviation for Mode().IsDir()
	sys     any         // underlying data source (can return nil)
}

func NewFileInfo(name string, size int64, mode fs.FileMode, modTime time.Time, isDir bool, sys any) fs.FileInfo {
	return &FileInfo{name, size, mode, modTime, isDir, sys}
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) Mode() fs.FileMode {
	return f.mode
}

func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *FileInfo) IsDir() bool {
	return f.isDir
}

func (f *FileInfo) Sys() any {
	return f.sys
}

func (f *FileInfo) SetSys(sys any) {
	f.sys = sys
}

func (f *FileInfo) SetMode(mode fs.FileMode) {
	f.mode = mode
}

func (f *FileInfo) SetModTime(modTime time.Time) {
	f.modTime = modTime
}

func (f *FileInfo) SetIsDir(isDir bool) {
	f.isDir = isDir
}

func (f *FileInfo) SetSize(size int64) {
	f.size = size
}

func (f *FileInfo) SetName(name string) {
	f.name = name
}
