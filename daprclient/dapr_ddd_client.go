package daprclient

import (
	"context"
	"encoding/json"
	"fmt"
	dapr_sdk_client "github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultMaxIdleConns        = 10
	DefaultMaxIdleConnsPerHost = 50
	DefaultIdleConnTimeout     = 5
	Protocol                   = "http"
)

type DaprHttpOptions struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int
}

type DaprDddClient interface {
	HttpGet(ctx context.Context, url string) *Response
	HttpPost(ctx context.Context, url string, reqData interface{}) *Response
	HttpPut(ctx context.Context, url string, reqData interface{}) *Response
	InvokeService(ctx context.Context, appID, methodName, verb string, request interface{}, response interface{}) (interface{}, error)
	LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error)
	ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventResponse, error)
	CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error)
	DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error)
	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
	DaprClient() (dapr_sdk_client.Client, error)
}

var _daprClient DaprDddClient

func GetDaprDDDClient() DaprDddClient {
	return _daprClient
}

func SetDaprDddClient(client DaprDddClient) {
	_daprClient = client
}

type daprDddClient struct {
	host           string
	httpPort       int64
	grpcPort       int64
	httpClient     *http.Client
	grpcDaprClient dapr_sdk_client.Client
}

type Option func(options *DaprHttpOptions)

func newHttpOptions() *DaprHttpOptions {
	options := &DaprHttpOptions{
		MaxIdleConns:        DefaultMaxIdleConns,
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:     DefaultIdleConnTimeout,
	}
	return options
}

func NewDaprDddClient(host string, httpPort int64, grpcPort int64, opts ...Option) (DaprDddClient, error) {
	options := newHttpOptions()
	for _, opt := range opts {
		opt(options)
	}

	grpcDaprClient, err := newDaprClient(host, grpcPort)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        options.MaxIdleConns,
			MaxIdleConnsPerHost: options.MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Second * time.Duration(options.IdleConnTimeout),
		},
	}

	return &daprDddClient{
		httpClient:     httpClient,
		host:           host,
		httpPort:       httpPort,
		grpcPort:       grpcPort,
		grpcDaprClient: grpcDaprClient,
	}, nil
}

func newDaprClient(host string, grpcPort int64) (dapr_sdk_client.Client, error) {

	// 三次试错创建daprClient
	port := strconv.FormatInt(grpcPort, 10)
	var grpcClient dapr_sdk_client.Client
	var err error
	var waitSecond time.Duration = 5
	addr := fmt.Sprintf("%s:%s", host, port)

	for i := 0; i < 4; i++ {
		if grpcClient, err = dapr_sdk_client.NewClientWithAddress(addr); err != nil {
			log.Infoln(fmt.Sprintf("dapr client connection error, address=%s", addr), err)
			continue
		}
		if grpcClient != nil {
			log.Infoln(fmt.Sprintf("dapr client connection success, address=%s", addr))
			break
		}
		time.Sleep(time.Second * waitSecond)
		waitSecond = 3
	}
	if err != nil {
		return nil, err
	}
	return grpcClient, nil
}

func (c *daprDddClient) InvokeService(ctx context.Context, appID, methodName, verb string, request interface{}, response interface{}) (interface{}, error) {
	var err error
	defer func() {
		if e := ddd_errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()
	var respBytes []byte

	if request != nil {
		reqBytes, err := json.Marshal(request)
		if err != nil {
			return nil, ddd_utils.NewAppError(appID, err)
		}
		content := &dapr_sdk_client.DataContent{
			ContentType: "application/json",
			Data:        reqBytes,
		}
		respBytes, err = c.grpcDaprClient.InvokeMethodWithContent(ctx, appID, methodName, verb, content)
	} else {
		respBytes, err = c.grpcDaprClient.InvokeMethod(ctx, appID, methodName, verb)
	}
	if err != nil {
		return nil, ddd_utils.NewAppError(appID, err)
	}
	if len(respBytes) > 0 {
		err = json.Unmarshal(respBytes, response)
		if err != nil {
			return nil, ddd_utils.NewAppError(appID, err)
		}
		return response, nil
	}
	return nil, nil
}

func (c *daprDddClient) DaprClient() (dapr_sdk_client.Client, error) {
	return c.grpcDaprClient, nil
}
