package helpers

func FindInI64Array(arr []int64, x int64) int {
	for i, j := range arr {
		if j == x {
			return i
		}
	}
	return -1
}
