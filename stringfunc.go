package main

import (
	"strings"
)

// SplitSpace 는 string 문자열을 공백을 기준으로 split하여 리스트를 반환하는 함수이다.
func SplitBySpace(str string) []string {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return nil
	}

	return strings.Split(str, " ")
}