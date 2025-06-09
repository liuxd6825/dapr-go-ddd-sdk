package light_rag

import "io"

type readCloserWrapper struct {
	io.Reader
}

func (r *readCloserWrapper) Close() error {
	return nil // strings.Reader无需实际关闭操作
}

func NewReadCloser(r io.Reader) io.ReadCloser {
	return &readCloserWrapper{r}
}
