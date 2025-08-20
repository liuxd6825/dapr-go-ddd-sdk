package amount

import "testing"

func Test_isNearMultiple(t *testing.T) {
	// 场景一：检查是否有意规避 10 万元的审批限额
	approvalLimit := 100000.0
	tolerancePercent := 0.03 // 3%的容忍度

	amounts := []float64{
		99850.00,  // 非常接近 10万 的 1倍
		150000.00, // 精确等于 10万 的 1.5倍，不应匹配整数倍
		199500.00, // 非常接近 10万 的 2倍
		205000.00, // 在 10万 的 2倍 的容忍度范围内（2.5%）
		240000.00, // 距离任何倍数都较远
		301000.00, // 非常接近 10万 的 3倍
		85000.00,  // 低于基准金额，但接近1倍
		499900.00, // 非常接近 50万 的 1倍，也接近 10万 的 5倍
	}
	for _, amount := range amounts {
		isNear, res := isNearMultiple(amount, approvalLimit, tolerancePercent)
		printNear(t, amount, isNear, res)
	}
}

func printNear(t *testing.T, name float64, isNear bool, res NearMatchResult) {
	t.Log(name)
	t.Log("isNear:", isNear, " differencePercent:", res.DifferencePercent, " difference ", res.Difference)
}
