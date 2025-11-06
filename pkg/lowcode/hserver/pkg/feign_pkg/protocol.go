package feign_pkg

type Protocol string

const (
	Protocol_Dapr  Protocol = "dapr"
	Protocol_Http  Protocol = "http"
	Protocol_Grpc  Protocol = "grpc"
	Protocol_Https Protocol = "https"
)

func (p Protocol) String() string {
	return string(p)
}
