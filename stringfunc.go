package main

import (
	"strings"
)

// SplitSpace 는 string 문자열을 공백을 기준으로 split하여 리스트를 반환하는 함수이다.
func SplitBySpace(str string) []string {
	str = strings.TrimSpace(str)
	var result []string
	if str == "" {
		return result
	}
	// 빈 문자열은 리스트에서 제외
	for _, s := range strings.Split(str, " ") {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
