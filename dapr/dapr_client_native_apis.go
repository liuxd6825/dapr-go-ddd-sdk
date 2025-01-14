package dapr

import (
	"context"
	pb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/actor/config"
	daprsdkclient "github.com/dapr/go-sdk/client"
	"google.golang.org/grpc"
	"io"
	"time"
)

func (c *daprClient) InvokeBinding(ctx context.Context, in *daprsdkclient.InvokeBindingRequest) (out *daprsdkclient.BindingEvent, err error) {
	return c.grpcClient.InvokeBinding(ctx, in)
}

func (c *daprClient) InvokeOutputBinding(ctx context.Context, in *daprsdkclient.InvokeBindingRequest) error {
	return c.grpcClient.InvokeOutputBinding(ctx, in)
}

func (c *daprClient) InvokeMethod(ctx context.Context, appID, methodName, verb string) (out []byte, err error) {
	return c.grpcClient.InvokeMethod(ctx, appID, methodName, verb)
}

func (c *daprClient) InvokeMethodWithContent(ctx context.Context, appID, methodName, verb string, content *daprsdkclient.DataContent) (out []byte, err error) {
	return c.grpcClient.InvokeMethodWithContent(ctx, appID, methodName, verb, content)
}

func (c *daprClient) InvokeMethodWithCustomContent(ctx context.Context, appID, methodName, verb string, contentType string, content interface{}) (out []byte, err error) {
	return c.grpcClient.InvokeMethodWithCustomContent(ctx, appID, methodName, verb, contentType, content)
}

func (c *daprClient) GetMetadata(ctx context.Context) (metadata *daprsdkclient.GetMetadataResponse, err error) {
	return c.grpcClient.GetMetadata(ctx)
}

func (c *daprClient) SetMetadata(ctx context.Context, key, value string) error {
	return c.grpcClient.SetMetadata(ctx, key, value)
}

func (c *daprClient) PublishEvent(ctx context.Context, pubsubName, topicName string, data interface{}, opts ...daprsdkclient.PublishEventOption) error {
	return c.grpcClient.PublishEvent(ctx, pubsubName, topicName, data, opts...)
}

func (c *daprClient) PublishEventfromCustomContent(ctx context.Context, pubsubName, topicName string, data interface{}) error {
	return c.grpcClient.PublishEventfromCustomContent(ctx, pubsubName, topicName, data)
}

func (c *daprClient) PublishEvents(ctx context.Context, pubsubName, topicName string, events []interface{}, opts ...daprsdkclient.PublishEventsOption) daprsdkclient.PublishEventsResponse {
	return c.grpcClient.PublishEvents(ctx, pubsubName, topicName, events, opts...)
}

func (c *daprClient) GetSecret(ctx context.Context, storeName, key string, meta map[string]string) (data map[string]string, err error) {
	return c.grpcClient.GetSecret(ctx, storeName, key, meta)
}

func (c *daprClient) GetBulkSecret(ctx context.Context, storeName string, meta map[string]string) (data map[string]map[string]string, err error) {
	return c.grpcClient.GetBulkSecret(ctx, storeName, meta)
}

func (c *daprClient) SaveState(ctx context.Context, storeName, key string, data []byte, meta map[string]string, so ...daprsdkclient.StateOption) error {
	return c.grpcClient.SaveState(ctx, storeName, key, data, meta, so...)
}

func (c *daprClient) SaveStateWithETag(ctx context.Context, storeName, key string, data []byte, etag string, meta map[string]string, so ...daprsdkclient.StateOption) error {
	return c.grpcClient.SaveStateWithETag(ctx, storeName, key, data, etag, meta, so...)
}

func (c *daprClient) SaveBulkState(ctx context.Context, storeName string, items ...*daprsdkclient.SetStateItem) error {
	return c.grpcClient.SaveBulkState(ctx, storeName, items...)
}

func (c *daprClient) GetState(ctx context.Context, storeName, key string, meta map[string]string) (item *daprsdkclient.StateItem, err error) {
	return c.grpcClient.GetState(ctx, storeName, key, meta)
}

func (c *daprClient) GetStateWithConsistency(ctx context.Context, storeName, key string, meta map[string]string, sc daprsdkclient.StateConsistency) (item *daprsdkclient.StateItem, err error) {
	return c.grpcClient.GetStateWithConsistency(ctx, storeName, key, meta, sc)
}

func (c *daprClient) GetBulkState(ctx context.Context, storeName string, keys []string, meta map[string]string, parallelism int32) ([]*daprsdkclient.BulkStateItem, error) {
	return c.grpcClient.GetBulkState(ctx, storeName, keys, meta, parallelism)
}

func (c *daprClient) QueryStateAlpha1(ctx context.Context, storeName, query string, meta map[string]string) (*daprsdkclient.QueryResponse, error) {
	return c.grpcClient.QueryStateAlpha1(ctx, storeName, query, meta)
}

func (c *daprClient) DeleteState(ctx context.Context, storeName, key string, meta map[string]string) error {
	return c.grpcClient.DeleteState(ctx, storeName, key, meta)
}

func (c *daprClient) DeleteStateWithETag(ctx context.Context, storeName, key string, etag *daprsdkclient.ETag, meta map[string]string, opts *daprsdkclient.StateOptions) error {
	return c.grpcClient.DeleteStateWithETag(ctx, storeName, key, etag, meta, opts)
}

func (c *daprClient) ExecuteStateTransaction(ctx context.Context, storeName string, meta map[string]string, ops []*daprsdkclient.StateOperation) error {
	return c.grpcClient.ExecuteStateTransaction(ctx, storeName, meta, ops)
}

func (c *daprClient) GetConfigurationItem(ctx context.Context, storeName, key string, opts ...daprsdkclient.ConfigurationOpt) (*daprsdkclient.ConfigurationItem, error) {
	return c.grpcClient.GetConfigurationItem(ctx, storeName, key, opts...)
}

func (c *daprClient) GetConfigurationItems(ctx context.Context, storeName string, keys []string, opts ...daprsdkclient.ConfigurationOpt) (map[string]*daprsdkclient.ConfigurationItem, error) {
	return c.grpcClient.GetConfigurationItems(ctx, storeName, keys, opts...)
}

func (c *daprClient) SubscribeConfigurationItems(ctx context.Context, storeName string, keys []string, handler daprsdkclient.ConfigurationHandleFunction, opts ...daprsdkclient.ConfigurationOpt) (string, error) {
	return c.grpcClient.SubscribeConfigurationItems(ctx, storeName, keys, handler, opts...)
}

func (c *daprClient) UnsubscribeConfigurationItems(ctx context.Context, storeName string, id string, opts ...daprsdkclient.ConfigurationOpt) error {
	return c.grpcClient.UnsubscribeConfigurationItems(ctx, storeName, id, opts...)
}

func (c *daprClient) Subscribe(ctx context.Context, opts daprsdkclient.SubscriptionOptions) (*daprsdkclient.Subscription, error) {
	return c.grpcClient.Subscribe(ctx, opts)
}

func (c *daprClient) SubscribeWithHandler(ctx context.Context, opts daprsdkclient.SubscriptionOptions, handler daprsdkclient.SubscriptionHandleFunction) (func() error, error) {
	return c.grpcClient.SubscribeWithHandler(ctx, opts, handler)
}

func (c *daprClient) DeleteBulkState(ctx context.Context, storeName string, keys []string, meta map[string]string) error {
	return c.grpcClient.DeleteBulkState(ctx, storeName, keys, meta)
}

func (c *daprClient) DeleteBulkStateItems(ctx context.Context, storeName string, items []*daprsdkclient.DeleteStateItem) error {
	return c.grpcClient.DeleteBulkStateItems(ctx, storeName, items)
}

func (c *daprClient) TryLockAlpha1(ctx context.Context, storeName string, request *daprsdkclient.LockRequest) (*daprsdkclient.LockResponse, error) {
	return c.grpcClient.TryLockAlpha1(ctx, storeName, request)
}

func (c *daprClient) UnlockAlpha1(ctx context.Context, storeName string, request *daprsdkclient.UnlockRequest) (*daprsdkclient.UnlockResponse, error) {
	return c.grpcClient.UnlockAlpha1(ctx, storeName, request)
}

func (c *daprClient) Encrypt(ctx context.Context, in io.Reader, opts daprsdkclient.EncryptOptions) (io.Reader, error) {
	return c.grpcClient.Encrypt(ctx, in, opts)
}

func (c *daprClient) Decrypt(ctx context.Context, in io.Reader, opts daprsdkclient.DecryptOptions) (io.Reader, error) {
	return c.grpcClient.Decrypt(ctx, in, opts)
}

func (c *daprClient) Shutdown(ctx context.Context) error {
	return c.grpcClient.Shutdown(ctx)
}

func (c *daprClient) Wait(ctx context.Context, timeout time.Duration) error {
	return c.grpcClient.Wait(ctx, timeout)
}

func (c *daprClient) WithTraceID(ctx context.Context, id string) context.Context {
	return c.grpcClient.WithTraceID(ctx, id)
}

func (c *daprClient) WithAuthToken(token string) {
	c.grpcClient.WithAuthToken(token)
}

func (c *daprClient) Close() {
	c.grpcClient.Close()
}

func (c *daprClient) RegisterActorTimer(ctx context.Context, req *daprsdkclient.RegisterActorTimerRequest) error {
	return c.grpcClient.RegisterActorTimer(ctx, req)
}

func (c *daprClient) UnregisterActorTimer(ctx context.Context, req *daprsdkclient.UnregisterActorTimerRequest) error {
	return c.grpcClient.UnregisterActorTimer(ctx, req)
}

func (c *daprClient) RegisterActorReminder(ctx context.Context, req *daprsdkclient.RegisterActorReminderRequest) error {
	return c.grpcClient.RegisterActorReminder(ctx, req)
}

func (c *daprClient) UnregisterActorReminder(ctx context.Context, req *daprsdkclient.UnregisterActorReminderRequest) error {
	return c.grpcClient.UnregisterActorReminder(ctx, req)
}

func (c *daprClient) InvokeActor(ctx context.Context, req *daprsdkclient.InvokeActorRequest) (*daprsdkclient.InvokeActorResponse, error) {
	return c.grpcClient.InvokeActor(ctx, req)
}

func (c *daprClient) GetActorState(ctx context.Context, req *daprsdkclient.GetActorStateRequest) (data *daprsdkclient.GetActorStateResponse, err error) {
	return c.grpcClient.GetActorState(ctx, req)
}

func (c *daprClient) SaveStateTransactionally(ctx context.Context, actorType, actorID string, operations []*daprsdkclient.ActorStateOperation) error {
	return c.grpcClient.SaveStateTransactionally(ctx, actorType, actorID, operations)
}

func (c *daprClient) ImplActorClientStub(actorClientStub actor.Client, opt ...config.Option) {
	c.grpcClient.ImplActorClientStub(actorClientStub, opt...)
}

func (c *daprClient) StartWorkflowBeta1(ctx context.Context, req *daprsdkclient.StartWorkflowRequest) (*daprsdkclient.StartWorkflowResponse, error) {
	return c.grpcClient.StartWorkflowBeta1(ctx, req)
}

func (c *daprClient) GetWorkflowBeta1(ctx context.Context, req *daprsdkclient.GetWorkflowRequest) (*daprsdkclient.GetWorkflowResponse, error) {
	return c.grpcClient.GetWorkflowBeta1(ctx, req)
}

func (c *daprClient) PurgeWorkflowBeta1(ctx context.Context, req *daprsdkclient.PurgeWorkflowRequest) error {
	return c.grpcClient.PurgeWorkflowBeta1(ctx, req)
}

func (c *daprClient) TerminateWorkflowBeta1(ctx context.Context, req *daprsdkclient.TerminateWorkflowRequest) error {
	return c.grpcClient.TerminateWorkflowBeta1(ctx, req)
}

func (c *daprClient) PauseWorkflowBeta1(ctx context.Context, req *daprsdkclient.PauseWorkflowRequest) error {
	return c.grpcClient.PauseWorkflowBeta1(ctx, req)
}

func (c *daprClient) ResumeWorkflowBeta1(ctx context.Context, req *daprsdkclient.ResumeWorkflowRequest) error {
	return c.grpcClient.ResumeWorkflowBeta1(ctx, req)
}

func (c *daprClient) RaiseEventWorkflowBeta1(ctx context.Context, req *daprsdkclient.RaiseEventWorkflowRequest) error {
	return c.grpcClient.RaiseEventWorkflowBeta1(ctx, req)
}

func (c *daprClient) ScheduleJobAlpha1(ctx context.Context, req *daprsdkclient.Job) error {
	return c.grpcClient.ScheduleJobAlpha1(ctx, req)
}

func (c *daprClient) GetJobAlpha1(ctx context.Context, name string) (*daprsdkclient.Job, error) {
	return c.grpcClient.GetJobAlpha1(ctx, name)
}

func (c *daprClient) DeleteJobAlpha1(ctx context.Context, name string) error {
	return c.grpcClient.DeleteJobAlpha1(ctx, name)
}

func (c *daprClient) GrpcClient() pb.DaprClient {
	return c.grpcClient.GrpcClient()
}

func (c *daprClient) GrpcClientConn() *grpc.ClientConn {
	return c.grpcClient.GrpcClientConn()
}

func (c *daprClient) LoadDomainEvents(ctx context.Context, request *pb.LoadDomainEventRequest) (*pb.LoadDomainEventResponse, error) {
	return c.grpcClient.LoadDomainEvents(ctx, request)
}

func (c *daprClient) SaveDomainEventSnapshot(ctx context.Context, request *pb.SaveDomainEventSnapshotRequest) (*pb.SaveDomainEventSnapshotResponse, error) {
	return c.grpcClient.SaveDomainEventSnapshot(ctx, request)
}

func (c *daprClient) CommitDomainEvents(ctx context.Context, request *pb.CommitDomainEventsRequest) (*pb.CommitDomainEventsResponse, error) {
	return c.grpcClient.CommitDomainEvents(ctx, request)
}

func (c *daprClient) RollbackDomainEvents(ctx context.Context, request *pb.RollbackDomainEventsRequest) (*pb.RollbackDomainEventsResponse, error) {
	return c.grpcClient.RollbackDomainEvents(ctx, request)
}

func (c *daprClient) ApplyDomainEvent(ctx context.Context, request *pb.ApplyDomainEventRequest) (*pb.ApplyDomainEventResponse, error) {
	return c.grpcClient.ApplyDomainEvent(ctx, request)
}

func (c *daprClient) GetDomainEventRelations(ctx context.Context, request *pb.GetDomainEventRelationsRequest) (*pb.GetDomainEventRelationsResponse, error) {
	return c.grpcClient.GetDomainEventRelations(ctx, request)
}

func (c *daprClient) GetDomainEvents(ctx context.Context, request *pb.GetDomainEventsRequest) (*pb.GetDomainEventsResponse, error) {
	return c.grpcClient.GetDomainEvents(ctx, request)
}

func (c *daprClient) WriteAppEventLog(ctx context.Context, request *pb.WriteAppEventLogRequest) (*pb.WriteAppEventLogResponse, error) {
	return c.grpcClient.WriteAppEventLog(ctx, request)
}

func (c *daprClient) UpdateAppEventLog(ctx context.Context, request *pb.UpdateAppEventLogRequest) (*pb.UpdateAppEventLogResponse, error) {
	return c.grpcClient.UpdateAppEventLog(ctx, request)
}

func (c *daprClient) GetAppEventLogByCommandId(ctx context.Context, request *pb.GetAppEventLogByCommandIdRequest) (*pb.GetAppEventLogByCommandIdResponse, error) {
	return c.grpcClient.GetAppEventLogByCommandId(ctx, request)
}

func (c *daprClient) WriteAppLog(ctx context.Context, request *pb.WriteAppLogRequest) (*pb.WriteAppLogResponse, error) {
	return c.grpcClient.WriteAppLog(ctx, request)
}

func (c *daprClient) UpdateAppLog(ctx context.Context, request *pb.UpdateAppLogRequest) (*pb.UpdateAppLogResponse, error) {
	return c.grpcClient.UpdateAppLog(ctx, request)
}

func (c *daprClient) GetAppLogById(ctx context.Context, request *pb.GetAppLogByIdRequest) (*pb.GetAppLogByIdResponse, error) {
	return c.grpcClient.GetAppLogById(ctx, request)
}

func (c *daprClient) Command(ctx context.Context, request *pb.CommandRequest) (*pb.CommandResponse, error) {
	return c.grpcClient.Command(ctx, request)
}
