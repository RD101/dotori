package main

import "strconv"

// TotalPage 함수는 아이템의 갯수를 입력받아 필요한 총 페이지 수를 구한다.
func TotalPage(itemNum int) int {
	page := itemNum / *flagPagenum
	if itemNum%*flagPagenum != 0 {
		page++
	}
	return page
}

// PageToInt 함수는 페이지 문자를 받아서 Int형 페이지수를 반환한다.
func PageToInt(page string) int {
	n, err := strconv.Atoi(page)
	if err != nil {
		return 1 // 변환할 수 없는 문자라면, 1페이지를 반환한다.
	}
	return n
}

// PageToString 함수는 페이지 문자를 받아서 String형 페이지수를 반환한다.
func PageToString(page string) string {
	_, err := strconv.Atoi(page)
	if err != nil {
		return "1" // 변환할 수 없는 문자라면, 1페이지를 반환한다.
	}
	return page
}
