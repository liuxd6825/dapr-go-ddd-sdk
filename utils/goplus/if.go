package goplus

// IfElse 若expr成立，则返回a；否则返回b。
func IfElse[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}
