package service

import (
	"context"

	sysconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/configs"
	client "github.com/ory/kratos-client-go"
)

type OryService struct {
	ory *client.APIClient
	cfg sysconfigs.KratosConfig
}

func (t *OryService) init() {
	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{
		{URL: t.cfg.Url},
	}
	t.ory = client.NewAPIClient(cfg)
}

func NewOryService(cfg sysconfigs.KratosConfig) *OryService {
	service := &OryService{
		cfg: cfg,
	}
	service.init()
	return service
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

func (t *OryService) GetIdentityByCode(ctx context.Context, code string) (*client.Identity, error) {
	list, _, err := t.ory.IdentityAPI.ListIdentities(ctx).CredentialsIdentifier(code).Execute()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
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

func (t *OryService) GetSession(ctx context.Context, id string) (*client.SuccessfulNativeLogin, error) {
	session, _, err := t.ory.IdentityAPI.GetSession(ctx, id).Execute()
	if err != nil {
		return nil, err
	}

	snl := client.NewSuccessfulNativeLogin(*session)

	return snl, nil
}
