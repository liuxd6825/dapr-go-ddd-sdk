package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	case_dao "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/dao"
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
	caseDao *case_dao.CaseDao
}

var codeService *CodeService
var codeServiceOnce sync.Once

func NewCodeService() *CodeService {
	codeServiceOnce.Do(func() {
		codeService = &CodeService{
			typeDAO: dao.NewCodeTypeDao(config.DBKey),
			seqDAO:  dao.NewCodeSequenceDao(config.DBKey),
			caseDao: case_dao.NewCaseDao(config.DBKey),
		}
	})
	return codeService
}

func (s *CodeService) getCase(ctx context.Context, caseId string) (string, error) {
	vo, err := s.caseDao.FindById(ctx, caseId)
	if err != nil {
		return "", err
	}
	if vo.Code == "" {
		return "", errors.New("Case Code Not Found case.id:%s", caseId)
	}
	return vo.Code, nil
}

func (s *CodeService) NewBillCode(ctx context.Context, caseId, billType string) string {
	caseCode, err := s.getCase(ctx, caseId)
	if err != nil {
		panic(err)
	}
	return s.NewCode(ctx, fmt.Sprintf("%s-%s-", caseCode, billType))
}

// NewSuCode
// @Description: 新建可疑任务编号
func (s *CodeService) NewSuCode(ctx context.Context, caseId string) string {
	return s.NewBillCode(ctx, caseId, "SU")
}

// NewHumanCode
// @Description: 新建人员编号
func (s *CodeService) NewHumanCode(ctx context.Context, caseId string) string {
	return s.NewBillCode(ctx, caseId, "HM")
}

// NewCompanyCode
// @Description: 新建公司编号
func (s *CodeService) NewCompanyCode(ctx context.Context, caseId string) string {
	return s.NewBillCode(ctx, caseId, "CP")
}

// NewAccountCode
// @Description: 新建账号编号
func (s *CodeService) NewAccountCode(ctx context.Context, caseId string) string {
	return s.NewBillCode(ctx, caseId, "AC")
}

// NewContractCode
// @Description: 新建合同编号
func (s *CodeService) NewContractCode(ctx context.Context, caseId string) string {
	return s.NewBillCode(ctx, caseId, "CT")
}

// NewProductCode
// @Description: 新建产品编号
func (s *CodeService) NewProductCode(ctx context.Context, caseId string) string {
	return s.NewBillCode(ctx, caseId, "PR")
}

func (s *CodeService) NewCaseCode(ctx context.Context) string {
	cmd := command.NewCodeNewCommand()
	cmd.Data.Type = "CS"
	cmd.Data.Style = model.CodeStyle_Global
	cmd.Data.NoHead = true
	cmd.Data.Length = 3
	code, err := s.New(ctx, cmd)
	if err != nil {
		panic(err)
	}
	return code
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

func (s *CodeService) NewCodeLen(ctx context.Context, typeCode string, codeLen int, addHead bool) string {
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
	count := cmd.Data.Count
	if count <= 0 {
		count = 1
	}

	dateStr := ""
	if cmd.Data.Style != model.CodeStyle_Global {
		now := time.Now()
		dateStr = now.Format(style.DateFormat())
	}

	// 原子递增获取序号
	seq, err := s.seqDAO.NextSeq(ctx, typeCode, dateStr, count)
	if err != nil {
		return "", fmt.Errorf("failed to generate sequence: %v", err)
	}

	num36 := intutils.IntToBase36(seq)
	if len(num36) < numLength {
		num36 = strings.Repeat("0", numLength-len(num36)) + num36
	}

	if cmd.Data.NoHead && dateStr == "" {
		return num36, nil
	} else if cmd.Data.NoHead && dateStr != "" {
		return fmt.Sprintf("%s%s", dateStr, num36), nil
	}
	// 拼接最终编号
	fullCode := fmt.Sprintf("%s%s%s", typeCode, dateStr, num36)
	return fullCode, nil
}
