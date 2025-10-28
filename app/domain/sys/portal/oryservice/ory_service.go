package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	client "github.com/ory/kratos-client-go"
	"sync"
)

type OryService struct {
	ory    *client.APIClient
	oryCfg *config.OryConfig
}

var (
	_oryOnce          sync.Once
	_oryDomainService *OryService
)

func (t *OryService) loadConfig() {
	e := env.GetEnv()
	oryMeta := e.App.Meta["ory"]
	var err error
	if oryMeta == nil {
		panic("ory not found in env.app")
	}

	t.oryCfg, err = config.ReadOryConfig(oryMeta)
	if err != nil {
		panic("read OryConfig error" + err.Error())
	}
}

func (t *OryService) init() {
	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{
		{URL: fmt.Sprintf(t.oryCfg.Url)},
	}
	t.ory = client.NewAPIClient(cfg)
}

func NewOryService() *OryService {
	_oryOnce.Do(func() {
		_oryDomainService = &OryService{}
		_oryDomainService.loadConfig()
		_oryDomainService.init()
	})
	return _oryDomainService
}

func (t *OryService) CreateNativeLoginFlow(ctx context.Context) (*client.LoginFlow, error) {
	flow, _, err := t.ory.FrontendAPI.CreateNativeLoginFlow(ctx).Execute()
	if err != nil {
		return nil, err
	}
	return flow, nil
}

func (t *OryService) UpdateLoginFlow(ctx context.Context, flowId string, body client.UpdateLoginFlowBody) (*client.SuccessfulNativeLogin, error) {
	flow, _, err := t.ory.FrontendAPI.UpdateLoginFlow(ctx).Flow(flowId).UpdateLoginFlowBody(body).Execute()
	if err != nil {
		return nil, err
	}
	return flow, nil
}

func (t *OryService) CreateIdentityExecute(ctx context.Context, createIdentityBody client.CreateIdentityBody) (*client.Identity, error) {
	ident, _, err := t.ory.IdentityAPI.CreateIdentity(ctx).CreateIdentityBody(createIdentityBody).Execute()
	if err != nil {
		return nil, err
	}
	return ident, nil
}

func (t *OryService) GetIdentity(ctx context.Context, id string) (*client.Identity, error) {
	ident, _, err := t.ory.IdentityAPI.GetIdentity(ctx, id).Execute()
	if err != nil {
		return nil, err
	}
	return ident, nil
}

func (t *OryService) UpdateIdentity(ctx context.Context, id string, updateIdentityBody client.UpdateIdentityBody) (*client.Identity, error) {
	ident, _, err := t.ory.IdentityAPI.UpdateIdentity(ctx, id).UpdateIdentityBody(updateIdentityBody).Execute()
	if err != nil {
		return nil, err
	}
	return ident, nil
}

func (t *OryService) DeleteIdentity(ctx context.Context, id string) error {
	_, err := t.ory.IdentityAPI.DeleteIdentity(ctx, id).Execute()
	if err != nil {
		return err
	}
	return nil
}

func (t *OryService) CreateNativeSettingsFlow(ctx context.Context) (*client.SettingsFlow, error) {
	flow, _, err := t.ory.FrontendAPI.CreateNativeSettingsFlow(ctx).Execute()
	if err != nil {
		return nil, err
	}
	return flow, nil
}

func (t *OryService) UpdateSettingsFlow(ctx context.Context, flowId string, updateSettingsFlowBody client.UpdateSettingsFlowBody) (*client.SettingsFlow, error) {
	flow, _, err := t.ory.FrontendAPI.UpdateSettingsFlow(ctx).Flow(flowId).UpdateSettingsFlowBody(updateSettingsFlowBody).Execute()
	if err != nil {
		return nil, err
	}
	return flow, nil
}
