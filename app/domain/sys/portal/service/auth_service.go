package service

import (
	"context"

	service "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/oryservice"
	client "github.com/ory/kratos-client-go"
)

type AuthService struct {
	ory *service.OryService
}

func NewAuthService(dbKey string) *AuthService {
	s := &AuthService{
		ory: service.NewOryService(dbKey),
	}
	return s
}

func (t *AuthService) CreateLoginFlow(ctx context.Context) (*client.LoginFlow, error) {
	return t.ory.CreateNativeLoginFlow(ctx)
}

func (t *AuthService) Login(ctx context.Context, flowId string, body client.UpdateLoginFlowBody) (*client.SuccessfulNativeLogin, error) {
	return t.ory.UpdateLoginFlow(ctx, flowId, body)
}

func (t *AuthService) Logout(ctx context.Context, flowId string, body client.UpdateLoginFlowBody) (*client.SuccessfulNativeLogin, error) {
	return t.ory.UpdateLoginFlow(ctx, flowId, body)
}
