package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	daprsdkclient "github.com/liuxd6825/dapr-go-sdk/client"
	pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
	"google.golang.org/grpc/status"
	"strings"

	// "google.golang.org/grpc/internal/status"
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
	MaxCallRecvMsgSize  int
}

// DaprClient
// @Description: Dapr客户端
type DaprClient interface {
	HttpGet(ctx context.Context, url string) *HttpResponse
	HttpPost(ctx context.Context, url string, reqData interface{}) *HttpResponse
	HttpPut(ctx context.Context, url string, reqData interface{}) *HttpResponse
	HttpDelete(ctx context.Context, url string, reqData interface{}) *HttpResponse

	WriteEventLog(ctx context.Context, req *WriteEventLogRequest) (resp *WriteEventLogResponse, resErr error)
	UpdateEventLog(ctx context.Context, req *UpdateEventLogRequest) (resp *UpdateEventLogResponse, resErr error)
	GetEventLogByCommandId(ctx context.Context, req *GetEventLogByCommandIdRequest) (resp *GetEventLogByCommandIdResponse, resErr error)
	WriteAppLog(ctx context.Context, req *WriteAppLogRequest) (resp *WriteAppLogResponse, resErr error)
	UpdateAppLog(ctx context.Context, req *UpdateAppLogRequest) (resp *UpdateAppLogResponse, resErr error)
	GetAppLogById(ctx context.Context, req *GetAppLogByIdRequest) (resp *GetAppLogByIdResponse, resErr error)

	InvokeService(ctx context.Context, appID, methodName, verb string, request interface{}, response interface{}) (interface{}, error)
	LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error)
	ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventResponse, error)
	CommitEvent(ctx context.Context, req *CommitRequest) (*CommitResponse, error)
	RollbackEvent(ctx context.Context, req *RollbackRequest) (*RollbackResponse, error)

	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
	GetRelations(ctx context.Context, req *GetRelationsRequest) (*GetRelationsResponse, error)
	GetEvents(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error)
	Client() (Client, error)
}

var _daprClient DaprClient

func GetDaprClient() DaprClient {
	return _daprClient
}

func SetDaprClient(client DaprClient) {
	_daprClient = client
}

func GetClient() (Client, error) {
	if _daprClient == nil {
		return nil, errors.New("not init client")
	}
	return _daprClient.Client()
}

type daprClient struct {
	host       string
	httpPort   int64
	grpcPort   int64
	httpClient *http.Client
	grpcClient daprsdkclient.Client
}

type Option func(options *DaprHttpOptions)

func newHttpOptions() *DaprHttpOptions {
	options := &DaprHttpOptions{
		MaxIdleConns:        DefaultMaxIdleConns,
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:     DefaultIdleConnTimeout,
		MaxCallRecvMsgSize:  GetMaxCallRecvMsgSize(),
	}
	return options
}

func NewDaprClient(ctx context.Context, host string, httpPort int64, grpcPort int64, opts ...Option) (DaprClient, error) {
	options := newHttpOptions()
	for _, opt := range opts {
		opt(options)
	}

	//重试5分钟进行Dapr连接
	grpcClient, err := newDaprClient(ctx, host, grpcPort, 5)
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
			MaxIdleConns:           options.MaxIdleConns,
			MaxIdleConnsPerHost:    options.MaxIdleConnsPerHost,
			IdleConnTimeout:        time.Second * time.Duration(options.IdleConnTimeout),
			MaxResponseHeaderBytes: 1024,
			WriteBufferSize:        1024 * 80,
			ReadBufferSize:         1024 * 80,
		},
	}

	return &daprClient{
		httpClient: httpClient,
		host:       host,
		httpPort:   httpPort,
		grpcPort:   grpcPort,
		grpcClient: grpcClient,
	}, nil
}

func newDaprClient(ctx context.Context, host string, grpcPort int64, retry uint) (daprsdkclient.Client, error) {
	port := strconv.FormatInt(grpcPort, 10)
	var grpcClient daprsdkclient.Client
	var err error
	addr := fmt.Sprintf("%s:%s", host, port)

	for i := 0; i <= int(retry); i++ {
		if grpcClient, err = newClientWithAddressContext(ctx, addr); err != nil {
			logs.Errorf(ctx, "", nil, "dapr %s", err.Error())
			continue
		}
		if grpcClient != nil {
			logs.Infof(ctx, "", nil, "dapr client connection success")
			break
		}
	}
	if err != nil {
		return nil, err
	}
	return grpcClient, nil
}

func (c *daprClient) tryCall(ctx context.Context, fun func() error, tryCount int, waitSecond time.Duration) error {
	var err error
	for i := 0; i < tryCount; i++ {
		err = fun()
		if errors.IsGrpcConnError(err) {
			time.Sleep(time.Second * waitSecond)
			grpcClient, err2 := newDaprClient(ctx, c.host, c.grpcPort, 1)
			if err2 != nil {
				return err2
			} else {
				c.grpcClient = grpcClient
			}
			continue
		}
		break
	}
	return err
}

// InvokeService
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param appID
//	@param methodName
//	@param verb
//	@param request
//	@param response
//	@return res
//	@return err
func (c *daprClient) InvokeService(ctx context.Context, appID, methodName, verb string, request interface{}, response interface{}) (res interface{}, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	var respBytes []byte
	if request != nil {
		reqBytes, err1 := json.Marshal(request)
		if err1 != nil {
			return nil, NewAppError(appID, err)
		}
		content := &daprsdkclient.DataContent{
			ContentType: "application/json",
			Data:        reqBytes,
		}
		err = c.tryCall(ctx, func() error {
			respBytes, err = c.grpcClient.InvokeMethodWithContent(ctx, appID, methodName, verb, content)
			return c.getError(err)
		}, 3, 1)

	} else {
		err = c.tryCall(ctx, func() error {
			respBytes, err = c.grpcClient.InvokeMethod(ctx, appID, methodName, verb)
			return c.getError(err)
		}, 3, 1)
	}
	if err != nil {
		return nil, NewAppError(appID, err)
	}
	if len(respBytes) > 0 {
		err = json.Unmarshal(respBytes, response)
		if err != nil {
			return nil, NewAppError(appID, err)
		}
		return response, nil
	}
	return nil, nil
}

func (c *daprClient) getError(err error) error {
	if err == nil {
		return err
	}
	st, ok := status.FromError(err)
	if ok {
		msg := ""
		details := st.Proto().GetDetails()
		if len(details) > 0 {
			for _, item := range details {
				value := string(item.Value)
				value = strings.ReplaceAll(value, "\n", "")
				msg = msg + value + "\n"
			}
			return errors.New(msg)
		}
	}
	return err
}

// CommitEvent
//
//	@Description:
//	@Description:
//	@receiver c
//	@param ctx
//	@param req
//	@return *CommitResponse
//	@return error
func (c *daprClient) CommitEvent(ctx context.Context, req *CommitRequest) (*CommitResponse, error) {
	resp := &CommitResponse{}
	in := &pb.CommitDomainEventsRequest{
		TenantId:  req.TenantId,
		SessionId: req.SessionId,
	}
	out, err := c.grpcClient.CommitDomainEvents(ctx, in)
	if out != nil {
		resp.Headers = NewResponseHeadersNil()
		resp.Headers.SetStatus(int32(out.Headers.Status))
		resp.Headers.SetValues(newMapStrings(out.Headers.Values))
		resp.Headers.SetMessage(out.Headers.Message)
	}
	return resp, err
}

func newPbMapStrings(mapValues map[string][]string) map[string]*pb.Strings {
	values := make(map[string]*pb.Strings)
	for k, v := range mapValues {
		values[k] = &pb.Strings{Values: v}
	}
	return values
}

func newMapStrings(mapValues map[string]*pb.Strings) map[string][]string {
	values := make(map[string][]string)
	for k, v := range mapValues {
		values[k] = v.Values
	}
	return values
}

// RollbackEvent
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param req
//	@return *RollbackResponse
//	@return error
func (c *daprClient) RollbackEvent(ctx context.Context, req *RollbackRequest) (*RollbackResponse, error) {
	resp := &RollbackResponse{}
	in := &pb.RollbackDomainEventsRequest{
		TenantId:  req.TenantId,
		SessionId: req.SessionId,
	}
	out, err := c.grpcClient.RollbackDomainEvents(ctx, in)
	if out != nil {
		resp.Headers = NewResponseHeadersNil()
		resp.Headers.SetStatus(int32(out.Headers.Status))
		resp.Headers.SetValues(newMapStrings(out.Headers.Values))
		resp.Headers.SetMessage(out.Headers.Message)
	}
	return resp, err
}

func (c *daprClient) Client() (Client, error) {
	return c.grpcClient, nil
}