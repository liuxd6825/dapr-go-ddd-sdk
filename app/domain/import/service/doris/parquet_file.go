package doris

import (
	"fmt"

	"github.com/parquet-go/parquet-go"
	"github.com/parquet-go/parquet-go/compress/snappy"
	"github.com/spf13/afero"
)

// ParquetFile 泛型 Parquet 文件写入器（基于 Apache parquet-go）。
//
// 使用示例：
//
//	pf, err := NewParquetFile[*Record]("/path/to/file.parquet")
//	if err != nil { return err }
//	for _, r := range records {
//	    if err := pf.Write(r); err != nil { return err }
//	}
//	if err := pf.WriteStop(); err != nil { return err }  // flush + footer + close 一步完成
//
// Close/WriteStop 等价。WriteStop 保留以兼容旧调用方。
type ParquetFile[T any] struct {
	filePath string
	// file holds the underlying afero.File that was opened by NewParquetFile.
	// The handle must be kept and closed here: with a streaming Fs such as
	// miniofs the act of closing the file is what finalises the upload, and
	// dropping the handle would leak both the file descriptor and the upload
	// goroutine, leaving the object permanently absent from the backend.
	file afero.File
	pw    *parquet.GenericWriter[T]
	closed bool
	fs     afero.Fs
}

// NewParquetFile 创建 Parquet 文件写入器，自动从 T 推断 schema，使用 SNAPPY 压缩。
//
// 注：parquet-go 不需要显式构造 schema 文件，直接从 Go struct 反射得到列名/类型/nullable
// 等元数据，输出完全符合 Apache Parquet 规范（Doris 原生兼容）。
func NewParquetFile[T any](fs afero.Fs, filePath string) (*ParquetFile[T], error) {
	f, err := fs.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file %s: %w", filePath, err)
	}
	pw := parquet.NewGenericWriter[T](f,
		parquet.Compression(&snappy.Codec{}),
	)
	return &ParquetFile[T]{
		filePath: filePath,
		file:     f,
		fs:       fs,
		pw:       pw,
	}, nil
}

// Write 写入一条记录。
func (p *ParquetFile[T]) Write(data T) error {
	if p.closed {
		return fmt.Errorf("parquet file %s already closed", p.filePath)
	}
	if _, err := p.pw.Write([]T{data}); err != nil {
		return fmt.Errorf("write parquet row: %w", err)
	}
	return nil
}

// WriteStop 写入 footer 并关闭文件（parquet-go 的 Close 自动完成 flush + footer + close）。
// 与 Close 等价；保留方法名以兼容旧调用方。
func (p *ParquetFile[T]) WriteStop() error {
	return p.Close()
}

// Close 写入 footer 并关闭文件。
//
// The order matters: parquet-go's Close flushes its buffered pages and writes
// the file footer to the underlying file, so it must run before the file is
// closed. Closing the file afterwards finalises the streaming upload (e.g. for
// miniofs: it closes the io.Pipe, lets the PutObject goroutine drain, and
// surfaces any upload error here).
func (p *ParquetFile[T]) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true
	pwErr := p.pw.Close()
	// Even when the parquet writer failed we still want to release the
	// underlying handle so the streaming PutObject (if any) is not left
	// waiting on a pipe that will never be closed.
	fileErr := p.file.Close()
	if pwErr != nil {
		return fmt.Errorf("close parquet writer: %w", pwErr)
	}
	if fileErr != nil {
		return fmt.Errorf("close file %s: %w", p.filePath, fileErr)
	}
	return nil
}

// GetFileName 返回文件路径。
func (p *ParquetFile[T]) GetFileName() string {
	return p.filePath
}
