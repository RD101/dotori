package main

import (
	"reflect"
	"testing"
)

//특수문자 기준으로 split하는 함수 SplitbySign이 잘 작동하는지 테스트하는 함수
func Test_SplitbySign(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{{
		in:   "s0010_c0010_ani_v001", // _만 포함된 경우
		want: []string{"s0010", "c0010", "ani", "v001"},
	}, {
		in:   "maya_thumbnail.test_check", // _, . 포함 경우
		want: []string{"maya", "thumbnail", "test", "check"},
	}, {
		in:   "アニメ/_の_ん_テスト", // 일본어, '/' 포함된 경우, 특수문자가 연속으로 있을 경우
		want: []string{"アニメ", "の", "ん", "テスト"},
	}, {
		in:   "挿絵_test_v002.mb", // 한자 포함된 경우
		want: []string{"挿絵", "test", "v002", "mb"},
	}, {
		in:   "SS0010_RIG_2.mb", // 대문자, 숫자 포함된 경우
		want: []string{"SS0010", "RIG", "2", "mb"},
	},
	}
	for _, c := range cases {
		b, _ := SplitBySign(c.in)
		if !reflect.DeepEqual(b, c.want) {
			t.Fatalf("Test_SplitbySign(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.in, c.want, b)
		}
	}
}
