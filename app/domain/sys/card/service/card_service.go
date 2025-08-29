package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"sync"
)

type CardService struct {
	dao *dao.CardDao
	xbase.Service
}

var (
	_cardOnce          sync.Once
	_cardDomainService *CardService
)

func NewCardService() *CardService {
	_cardOnce.Do(func() {
		_cardDomainService = &CardService{
			dao: dao.NewCardDao(config.DBKey),
		}
	})
	return _cardDomainService
}

func (t *CardService) Create(ctx context.Context, cmd *command.CardCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *CardService) Delete(ctx context.Context, cmd *command.CardDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *CardService) Update(ctx context.Context, cmd *command.CardUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *CardService) CreateMany(ctx context.Context, cards []*model.Card) (int64, error) {
	if len(cards) == 0 {
		return 0, nil
	}
	res := t.dao.CreateMany(ctx, cards)
	return res.GetRowsAffected(), res.GetError()
}

func (t *CardService) UpdateMany(ctx context.Context, cards []*model.Card) (int64, error) {
	if len(cards) == 0 {
		return 0, nil
	}
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"group_id", "title_text", "subtitle_text", "icon", "url", "opened_target", "width", "height", "order_num", "card_level", "card_class", "card_mode", "card_html", "link_html"})
	res := t.dao.UpdateMany(ctx, cards, opts)
	return res.GetRowsAffected(), res.GetError()
}

func (t *CardService) DeleteMany(ctx context.Context, listId []string) (int64, error) {
	if len(listId) == 0 {
		return 0, nil
	}
	res := t.dao.DeleteByIds(ctx, listId)
	return res.GetRowsAffected(), res.GetError()
}

func (t *CardService) DeleteByListGroupId(ctx context.Context, listGroupId []string) (int64, error) {
	if len(listGroupId) == 0 {
		return 0, nil
	}
	var totalRows int64 = 0
	for _, id := range listGroupId {
		res := t.dao.DeleteByRSQL(ctx, fmt.Sprintf("group_id=='%s'", id))
		if res.GetError() != nil {
			return 0, res.GetError()
		}
		totalRows += res.GetRowsAffected()
	}
	return totalRows, nil
}

func (t *CardService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Card, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *CardService) FindPaging(ctx context.Context, caseId string, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Card], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *CardService) FindByHomeId(ctx context.Context, homeId string) ([]*model.Card, error) {
	rsql := fmt.Sprintf("home_id=='%s'", homeId)
	return t.dao.FindByRSQL(ctx, rsql)
}

func (t *CardService) FindByGroupId(ctx context.Context, groupId string) ([]*model.Card, error) {
	rsql := fmt.Sprintf("group_id=='%s'", groupId)
	return t.dao.FindByRSQL(ctx, rsql)
}
