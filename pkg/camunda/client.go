package camunda

import (
	"fmt"
	camundaclientgo "github.com/citilinkru/camunda-client-go/v3"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"time"
)

type NewConfig struct {
	UserAgent           string
	Url                 string
	Timeout             time.Duration
	ApiUser             string
	ApiPassword         string
	AuthorizationHeader string
}

func NewClient(cfg *NewConfig) *camundaclientgo.Client {
	if cfg == nil {
		panic(errors.New("config is nil"))
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = time.Second * 10
	}
	client := camundaclientgo.NewClient(camundaclientgo.ClientOptions{
		EndpointUrl: cfg.Url,
		ApiUser:     cfg.ApiUser,
		ApiPassword: cfg.ApiPassword,
		Timeout:     cfg.Timeout,
	})
	return client
}

func NewDaprUrl(daprHost string, daprPort int, serviceName string) string {
	return fmt.Sprintf("http://%s:%d/v1.0/invoke/%s/method/engine-rest", daprHost, daprPort, serviceName)
}

func NewTestClient() *camundaclientgo.Client {
	cfg := &NewConfig{
		Url:         "http://localhost:8080/engine-rest",
		ApiUser:     "demo",
		ApiPassword: "demo",
		Timeout:     time.Second * 10,
	}
	return NewClient(cfg)
}

func String(str string) *string {
	return &str
}
