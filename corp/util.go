package corp

// RemoveDuplicateString 删除重复的字符串元素
func RemoveDuplicateString(a []string) []string {
	exists := make(map[string]struct{})
	result := make([]string, 0, len(a))
	for _, s := range a {
		if _, ok := exists[s]; ok {
			continue
		}
		result = append(result, s)
		exists[s] = struct{}{}
	}
	return result
}
