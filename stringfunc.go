package main

import (
	"errors"
	"path/filepath"
	"sort"
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

// Str2List 함수는 string 문자열을 " ". , 기준으로 tag 리스트를 반환하는 함수이다.
func Str2List(str string) []string {
	var result []string
	if str == "" {
		return result
	}
	str = strings.TrimSpace(str) // 앞뒤에 붙어있는 Space를 제거한다.
	// 빈 문자열은 리스트에서 제외
	for _, s := range strings.Split(str, " ") {
		for _, w := range strings.Split(s, ",") {
			if w != "" {
				result = append(result, strings.TrimSpace(w))
			}
		}
	}
	// 태그순으로 정렬한다.
	sort.Strings(result)
	return result
}

// StringToMap 함수는 "key:value,key:value" 형식의 문자열을 map 형으로 변환하는 함수이다.
func StringToMap(str string) (map[string]string, error) {
	result := make(map[string]string)

	str = strings.TrimSpace(str)
	if str == "" {
		return result, nil
	}

	if !regexMap.MatchString(str) { // 전달받은 str이 key:value,key:value 형식이 맞는지 확인
		return result, errors.New("map 형식이 아닙니다")
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

	return result, nil
}

// FilenameToTags 는 경로를 받아서 태그를 반환한다.
func FilenameToTags(path string) []string {
	var returnTags []string
	filename := strings.TrimSuffix(path, filepath.Ext(path)) // 확장자 제거
	tags := regexSplitBySign.Split(filename, -1)
	for _, tag := range tags {
		if tag != "thumbnail" {
			returnTags = append(returnTags, tag)
		}
	}
	if len(returnTags) == 0 {
		return returnTags
	}
	return returnTags
}
