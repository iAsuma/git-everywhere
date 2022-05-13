package qslice

// ContainsInSliceString 判断字符串是否在 slice 中
func ContainsInSliceString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
