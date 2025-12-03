package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/intutils"
)

type CodeService struct {
	typeDAO *dao.CodeTypeDao
	seqDAO  *dao.CodeSequenceDao
}

var codeService *CodeService
var codeServiceOnce sync.Once

func NewCodeService() *CodeService {
	codeServiceOnce.Do(func() {
		codeService = &CodeService{
			typeDAO: dao.NewCodeTypeDao(config.DBKey),
			seqDAO:  dao.NewCodeSequenceDao(config.DBKey),
		}
	})
	return codeService
}

func (s *CodeService) NewCode(ctx context.Context, typeCode string) string {
	cmd := command.NewCodeNewCommand()
	cmd.CommandId = idutils.NewId()
	cmd.Data.Type = typeCode
	cmd.Data.Style = model.CodeStyle_Global
	cmd.Data.Length = 4

	fullCode, err := s.New(ctx, cmd)
	if err != nil {
		panic(err)
	}
	return fullCode
}

func (s *CodeService) NewCodeLen(ctx context.Context, typeCode string, codeLen int) string {
	cmd := command.NewCodeNewCommand()
	cmd.CommandId = idutils.NewId()
	cmd.Data.Type = typeCode
	cmd.Data.Style = model.CodeStyle_Global
	cmd.Data.Length = codeLen

	fullCode, err := s.New(ctx, cmd)
	if err != nil {
		panic(err)
	}
	return fullCode
}

// New
// @Description:
// @receiver s
// @param ctx
// @param typeCode
// @param codeStyle
// @param numLength
// @return string
// @return error
func (s *CodeService) New(ctx context.Context, cmd *command.CodeNewCommand) (string, error) {
	if cmd == nil {
		return "", errors.New("cmd is nil")
	}
	typeCode := cmd.Data.Type
	style := cmd.Data.Style
	numLength := cmd.Data.Length

	dateStr := ""
	if cmd.Data.Style != model.CodeStyle_Global {
		now := time.Now()
		dateStr = now.Format(style.DateFormat())
	}

	// 原子递增获取序号
	seq, err := s.seqDAO.NextSeq(ctx, typeCode, dateStr)
	if err != nil {
		return "", fmt.Errorf("failed to generate sequence: %v", err)
	}

	num36 := intutils.IntToBase36(seq)
	if len(num36) < numLength {
		num36 = strings.Repeat("0", numLength-len(num36)) + num36
	}

	// 拼接最终编号
	fullCode := fmt.Sprintf("%s%s%s", typeCode, dateStr, num36)
	return fullCode, nil
}
