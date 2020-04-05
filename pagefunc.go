package main

import "strconv"

// TotalPage 함수는 아이템의 갯수를 입력받아 필요한 총 페이지 수를 구한다.
func TotalPage(itemNum, limitnum int64) int64 {
	page := itemNum / limitnum
	if itemNum%limitnum != 0 {
		page++
	}
	return page
}

// PageToInt 함수는 페이지 문자를 받아서 Int형 페이지수를 반환한다.
func PageToInt(page string) int64 {
	n, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return 1 // 변환할 수 없는 문자라면, 1페이지를 반환한다.
	}
	return n
}

// PageToString 함수는 페이지 문자를 받아서 String형 페이지수를 반환한다.
func PageToString(page string) string {
	// url에서 "&page=1"이 아닌 "&page=1#" 로 실수로 입력했을 때,
	// 변환할 수 없는 문자라면, 1페이지를 반환하도록 하기 위해 이 함수가 존재한다.
	_, err := strconv.Atoi(page)
	if err != nil {
		return "1"
	}
	return page
}

// PreviousPage 함수는 이전 페이지를 반환한다.
func PreviousPage(current, maxnum int64) int64 {
	if maxnum < current {
		return maxnum
	}
	if current == 1 {
		return 1
	}
	return current - 1
}

// NextPage 함수는 다음 페이지를 반환한다.
func NextPage(current, maxnum int64) int64 {
	if maxnum <= current {
		return maxnum
	}
	return current + 1
}
