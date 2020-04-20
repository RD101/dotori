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

// Int2Status 는 입력받은 ItemStatus에 해당하는 문자열을 반환한다.
func Int2Status(i ItemStatus) string {
	switch i {
	case 0:
		return "Ready"
	case 1:
		return "Copying"
	case 2:
		return "Copied"
	case 3:
		return "CreatingThumbnail"
	case 4:
		return "CreatedThumbnail"
	case 5:
		return "CreatingContainer"
	case 6:
		return "CreatedContainer"
	case 7:
		return "Done"
	default:
		return ""
	}
}
