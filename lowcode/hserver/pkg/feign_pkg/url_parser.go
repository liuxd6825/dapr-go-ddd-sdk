package feign_pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"strings"
)

type URLParser struct {
	Protocol    Protocol `json:"protocol,omitempty"`
	ServiceName string   `json:"serviceName,omitempty"`
	Port        string   `json:"port,omitempty"`
	Path        string   `json:"path,omitempty"`
}

// NewURLParser
//
//	@Description:
//	@param str   http://user-service:8080/api/v1/users
//	@return *URLParser
func NewURLParser(str string) (*URLParser, error) {
	protocol := Protocol_Dapr
	s := strings.ToLower(str)
	if strings.HasPrefix(s, Protocol_Dapr.String()+"://") {
		protocol = Protocol_Dapr
	} else if strings.HasPrefix(s, Protocol_Http.String()+"://") {
		protocol = Protocol_Http
	} else if strings.HasPrefix(s, Protocol_Grpc.String()+"://") {
		protocol = Protocol_Grpc
	} else if strings.HasPrefix(s, Protocol_Https.String()+"://") {
		protocol = Protocol_Https
	}
	serviceName := ""
	port := ""
	path := ""
	s = str[len(protocol.String())+3:]
	list := strings.Split(s, "/")
	count := len(list)
	if count == 0 {
		return nil, errors.New("")
	} else if count >= 1 { // user-service:8080
		l := strings.Split(list[0], ":")
		c := len(l)
		if c == 1 {
			serviceName = l[0]
		} else if c >= 2 {
			serviceName = l[0]
			port = l[1]
		}
		path = "/" + strings.Join(list[1:], "/")
	}

	return &URLParser{Protocol: Protocol(protocol), ServiceName: serviceName, Port: port, Path: path}, nil
}
