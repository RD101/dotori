package main

import (
	"testing"
)

// 시간형식을 테스트하기 위한 함수
func Test_checkTime(t *testing.T) {
	cases := []struct {
		time string
		want bool
	}{{
		time: "2019-09-13T22:04:32+09:00",
		want: true,
	}, {
		time: "2019-09-32T22:04:32+09:00", // 날짜가 틀렸을 경우
		want: false,
	}, {
		time: "2019-09-13T22:04:75+09:00", // 초가 틀렸을 경우
		want: false,
	}, {
		time: "2019-09-13T22:64:32+09:00", // 분이 틀렸을 경우
		want: false,
	}, {
		time: "2019-13-13T22:04:32+09:00", // 월이 틀렸을 경우
		want: false,
	},
	}

	for _, c := range cases {
		b := regexRFC3339Time.MatchString(c.time)
		if c.want != b {
			t.Fatalf("Test_checkTime(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.time, c.want, b)
		}
	}
}
