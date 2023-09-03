package utils

import "unicode/utf8"

func FindLongestStringInSlice(str []string) string {
	if len(str) == 0 {
		return ""
	}
	maxLen := utf8.RuneCountInString(str[0])
	result := str[0]
	for i := 0; i < len(str); i++ {
		if utf8.RuneCountInString(str[i]) > maxLen {
			maxLen = utf8.RuneCountInString(str[i])
			result = str[i]
		}
		if len(str[i]) == maxLen && str[i] > result {
			maxLen = utf8.RuneCountInString(str[i])
			result = str[i]
		}
	}
	return result
}
