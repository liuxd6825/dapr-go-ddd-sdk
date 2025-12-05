package amount

import (
	"context"
	"fmt"
	"math"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
)

// NearMatchResult 存储匹配结果的详细信息
type NearMatchResult struct {
	BaseAmount        float64 // 用于计算的基准金额
	Multiple          int     // 匹配到的倍数 (N)
	TargetAmount      float64 // 理论上的目标金额 (N * BaseAmount)
	Difference        float64 // 实际金额与目标金额的差额
	DifferencePercent float64 // 差额百分比
}

// AmountNear
// @Description: 临近交易
type AmountNear struct {
	rule    model2.AmountNearRule
	percent float64
}

func NewAmountNear(rule model2.AmountNearRule) *AmountNear {
	return &AmountNear{
		rule:    rule,
		percent: rule.NearPercent / 100,
	}
}

func (s *AmountNear) IsEnable() bool {
	return s.rule.IsEnable
}

// DoAction 临近交易
func (s *AmountNear) DoAction(ctx context.Context, tx *model.Record, txIndex int, accTxs *model2.AccountRecords, result *action.AnalyseResult) {
	isNear, res := isNearMultiple(tx.Amount, s.rule.NearAmount, s.percent)
	if isNear {
		result.AddRecord(tx, model2.SuType_AmountNear, fmt.Sprintf("大额交易: 目标金额 %.2f, 差额 %.2f, 差额率 %.2f%%", res.TargetAmount, res.Difference, res.DifferencePercent))
	}
}

// isNearMultiple 检查单个金额是否是基准金额的近似倍数
// transactionAmount: 要检查的交易金额
// baseAmount: 基准金额（例如，审批阈值）
// tolerance: 容忍度，以小数表示（例如 0.01 代表 1%）
// 返回：是否匹配，以及匹配的详细信息
func isNearMultiple(transactionAmount, baseAmount float64, tolerance float64) (bool, NearMatchResult) {
	// --- 1. 输入校验 ---
	if baseAmount <= 0 {
		return false, NearMatchResult{} // 基准金额必须为正数
	}
	if transactionAmount < baseAmount*0.5 {
		// 如果交易金额连基准金额的一半都不到，不将其视为一个“倍数”
		return false, NearMatchResult{}
	}

	if tolerance < 0 || tolerance >= 1 {
		// 容忍度必须在 [0, 1) 区间内
		return false, NearMatchResult{}
	}

	// --- 2. 核心算法 ---
	// 计算比率
	ratio := transactionAmount / baseAmount

	// 四舍五入到最接近的整数倍
	closestN := int(math.Round(ratio))
	if closestN == 0 {
		return false, NearMatchResult{} // 避免0倍的情况
	}

	// 计算理论上的目标金额
	targetAmount := float64(closestN) * baseAmount

	// 计算实际差额
	difference := math.Abs(transactionAmount - targetAmount)

	// 计算允许的最大差额
	allowedDeviation := targetAmount * tolerance

	// --- 3. 判断并返回结果 ---
	if difference <= allowedDeviation {
		// 差额在容忍度范围内，判定为匹配
		return true, NearMatchResult{
			Multiple:          closestN,
			TargetAmount:      targetAmount,
			Difference:        difference,
			DifferencePercent: (difference / targetAmount) * 100,
		}
	}

	return false, NearMatchResult{}
}
