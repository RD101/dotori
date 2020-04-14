package main

import (
	"errors"
	"log"
	"path/filepath"
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

// PathToTags 는 경로를 받아서 태그를 반환한다.
func PathToTags(path string) ([]string, error) {
	var returnTags []string
	filename := strings.TrimSuffix(path, filepath.Ext(path)) // 확장자 제거
	tags := regexSplitBySign.Split(filename, -1)
	for _, tag := range tags {
		if tag != "thumbnail" {
			returnTags = append(returnTags, tag)
		}
	}
	if len(returnTags) == 0 {
		return returnTags, errors.New("태그리스트가 비어있습니다")
	}
	return returnTags, nil
}
