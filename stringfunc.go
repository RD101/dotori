package main

import (
	"log"
	"strings"
)

// SplitBySpace 는 string 문자열을 공백을 기준으로 split하여 리스트를 반환하는 함수이다.
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

// StringToMap 함수는 "key:value,key:value" 형식의 문자열을 map 형으로 변환하는 함수이다.
func StringToMap(str string) map[string]string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}

	var result map[string]string
	result = make(map[string]string)

	if !regexMap.MatchString(str) { // 전달받은 str이 key:value,key:value 형식이 맞는지 확인
		log.Fatal("map 형식이 아닙니다")
		return result
	}

	for _, s := range strings.Split(str, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		key := strings.Split(s, ":")[0]
		value := strings.Split(s, ":")[1]
		result[key] = value
	}

	return result
}

// SplitBySign 는 string 문자열을 특수문자 기준으로 split하여 리스트를 반환하는 함수이다.
func SplitBySign(str string) ([]string, error) {
	var result []string
	if !regexSplitbySign.MatchString(str) {
		log.Fatal("string 형식이 아닙니다")
	}
	result = regexSplitbySign.FindAllString(str, -1)
	if len(result) == 0 {
		log.Fatal("빈 리스트를 반환했습니다")
		return result, nil
	}
	return result, nil
}
