package light_rag

import (
	"net/http"
	"sync"
)

var cli *Client
var cliOnce sync.Once

func client() *Client {
	cliOnce.Do(func() {
		cli = NewClient("http://localhost:9621", &http.Client{})
	})
	return cli
}
