package common

import (
	"io"
	"mime/multipart"
)

type FormFile struct {
	file       multipart.File
	fileHeader *multipart.FileHeader
}

func NewFormFileResult(file multipart.File, fileHeader *multipart.FileHeader) *FormFile {
	return &FormFile{file: file, fileHeader: fileHeader}
}

func (f *FormFile) FileName() string {
	return f.fileHeader.Filename
}

func (f *FormFile) FileType() string {
	return f.fileHeader.Header.Get("Content-Type")
}

func (f *FormFile) FileSize() int64 {
	return f.fileHeader.Size
}

func (f *FormFile) Header() *multipart.FileHeader {
	return f.fileHeader
}

func (f *FormFile) Close() error {
	return f.file.Close()
}

func (f *FormFile) ReadAt(p []byte, off int64) (int, error) {
	return f.file.ReadAt(p, off)
}

func (f *FormFile) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

func (f *FormFile) ReadAll() []byte {
	// 创建一个字节缓冲区并读取文件内容
	fileBytes := make([]byte, f.FileSize())
	if _, err := io.ReadFull(f.file, fileBytes); err != nil {
		panic(err)
	}
	return fileBytes
}
