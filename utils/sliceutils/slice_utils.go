package sliceutils

// 数组长度
func Length[T any](data []T) int {
	return len(data) // Slice
}

// 数组长度是否大于0
func Has[T any](data []T) bool {
	return len(data) > 0 // Slice
}
