package utils

func FindLongestStringInSlice(str []string) string {
	if len(str) == 0 {
		return ""
	}
	maxLen := len(str[0])
	result := str[0]
	for i := 0; i < len(str); i++ {
		if len(str[i]) > maxLen {
			maxLen = len(str[i])
			result = str[i]
		}
		if len(str[i]) == maxLen && str[i] > result {
			maxLen = len(str[i])
			result = str[i]
		}
	}
	return result
}
