package restapp

type HttpType string

const (
	HttpGet    HttpType = "GET"
	HttpPost   HttpType = "POST"
	HttpDelete HttpType = "DELETE"
	HttpPut    HttpType = "PUT"
)

const (
	ContentType = "Content-Type"
)
