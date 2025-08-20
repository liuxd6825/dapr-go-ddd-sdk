package times

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/suspicious/action"
	"time"
)

// NonWorkingHours 非工时交易
type NonWorkingHours struct {
	rule model.TimeNonWorkingHoursRule
}

func NewNonWorkingHours(rule model.TimeNonWorkingHoursRule) *NonWorkingHours {
	return &NonWorkingHours{
		rule: rule,
	}
}

func (s *NonWorkingHours) IsEnable() bool {
	return s.rule.IsEnable
}

func (s *NonWorkingHours) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model.AccountRecords, result *action.AnalyseResult) {
	if tx.Amount < s.rule.TxAmount {
		return
	}
	t := tx.Date
	// 非工作日交易
	if s.rule.IsNonWorkingSunday {
		if t.Weekday() == time.Saturday {
			result.AddRecord(tx, model.SuType_TimeNonWorking, fmt.Sprintf("非工作时间交易, 周六%s", t.Format(time.DateTime)))
			return
		}
		if t.Weekday() == time.Sunday {
			result.AddRecord(tx, model.SuType_TimeNonWorking, fmt.Sprintf("非工作时间交易, 周日%s", t.Format(time.DateTime)))
			return
		}
	}
	// 非工作时间交易
	if s.rule.IsNonWorkingHours {
		if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 { // 排除时间是00:00:00的无时间部分交易
			return
		}
		if t.Hour() < s.rule.NonWorkingHoursMin || t.Hour() > s.rule.NonWorkingHoursMax {
			result.AddRecord(tx, model.SuType_TimeNonWorking, fmt.Sprintf("非工作时间交易, 早晚%s", t.Format(time.DateTime)))
			return
		}
	}

}
