package main

import (
	"strings"
)

// StringToMap 함수는 "key:value,key:value" 형식의 문자열을 map 형으로 변환하는 함수이다.
func StringToMap(str string) map[string]string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}

	var result map[string]string
	result = make(map[string]string)
	for _, s := range strings.Split(str, ",") {
		if s == "" {
			continue
		}

		key := strings.Split(s, ":")[0]
		value := strings.Split(s, ":")[1]
		result[key] = value
	}

	return result
}
