package main

import (
	"strings"
)

// Tags2str 템플릿함수는 태그리스트를 공백으로 분리된 문자열로 만든다.
func Tags2str(tags []string) string {
	var newtags []string
	for _, tag := range tags {
		if tag != "" {
			newtags = append(newtags, tag)
		}
	}
	return strings.Join(newtags, " ")
}

// add함수는 입력받은 두 정수를 더한 값을 반환한다.
func add(a, b int) int {
	return (a + b)
}

// RmRootpath 템플릿함수는 path가 rootpath로 시작하면 rootpath 문자열을 제거한다.
func RmRootpath(path, rootpath string) string {
	return strings.TrimLeft(path, rootpath)
}

// LastLog 함수는 마지막 로그를 반환한다.
func LastLog(logs []string) string {
	if len(logs) == 0 {
		return ""
	}
	return logs[len(logs)-1]
}

// SplitTimeData 함수는 createData의 T를 기준으로 나누어 년,월,일 만 반환한다.
func SplitTimeData(data string) string {
	splitData := strings.Split(data, "T")
	return splitData[0]
}

// ItemListLength 함수는 Item형 리스트 전체의 개수를 반환한다.
func ItemListLength(items []Item) int {
	return len(items)
}

// StringListLength 함수는 string형 리스트 전체의 개수를 반환한다.
func StringListLength(strs []string) int {
	return len(strs)
}
