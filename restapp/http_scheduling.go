package restapp

import (
	"errors"
	"github.com/dapr/go-sdk/service/common"
	"github.com/kataras/iris/v12/context"
)

// AddJobEventHandler
//
//	@Description:  Dapr服务方法
//	@receiver s
//	@param name
//	@param fn
//	@return error
func (s *HttpServer) AddJobEventHandler(name string, fn common.JobEventHandler) error {
	if name == "" {
		return errors.New("job event name cannot be empty")
	}

	if fn == nil {
		return errors.New("job event handler not supplied")
	}

	s.jobEventHandlers[name] = fn
	return nil
}

// OnJobEventAlpha1 is invoked by the sidecar following a scheduled job registered in
// the scheduler
func (s *HttpServer) OnJobEventAlpha1(ictx *context.Context) {
	/*
		// parse the job type from the method or name
		jobType, found := strings.CutPrefix(in.GetMethod(), "job/")
		if !found {
			if in.GetName() == "" {
				return &runtimepb.JobEventResponse{}, errors.New("unsupported invocation")
			}
			jobType = in.GetName()
		}

		if fn, ok := s.jobEventHandlers[jobType]; ok {
			e := &common.JobEvent{
				JobType: jobType,
				Data:    in.GetData().GetValue(),
			}
			if err := fn(ctx, e); err != nil {
				return nil, fmt.Errorf("error executing %s binding: %w", in.GetName(), err)
			}
			return &runtimepb.JobEventResponse{}, nil
		}
		return &runtimepb.JobEventResponse{}, errors.New("job event handler not found")

	*/
}
