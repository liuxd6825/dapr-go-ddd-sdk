package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
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

func (s *CodeService) New(ctx context.Context, typeCode string) (string, error) {
	// 1. 获取规则配置
	// 实际生产中，这一步应该加内存缓存(Redis/LocalCache)，避免每次都查库
	codeType, err := s.typeDAO.GetByCode(ctx, typeCode)
	if err != nil {
		return "", fmt.Errorf("db error: %v", err)
	}
	if codeType == nil {
		return "", fmt.Errorf("code rule not found for code: %s", typeCode)
	}

	// 2. 获取当前日期字符串 (Go layout: "060102" = YYMMDD)
	// 如果配置是 "20060102" 则是 YYYYMMDD
	now := time.Now()
	dateStr := now.Format(codeType.DateFormat)

	// 3. 原子递增获取序号
	seq, err := s.seqDAO.NextSeq(ctx, typeCode, dateStr)
	if err != nil {
		return "", fmt.Errorf("failed to generate sequence: %v", err)
	}

	// 4. 拼接最终编号
	// 格式: Prefix + Date + PaddedSeq (e.g. MH- + 250101 + 000001)
	fullCode := fmt.Sprintf("%s%s%0*d", codeType.Prefix, dateStr, codeType.SeqLength, seq)
	return fullCode, nil
}

func (s *CodeService) NewCode(ctx context.Context, typeCode string) string {
	fullCode, err := s.New(ctx, typeCode)
	if err != nil {
		panic(err)
	}
	return fullCode
}
