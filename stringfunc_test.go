package main

import (
	"reflect"
	"testing"
)

//FilenameToTags 함수가 잘 작동하는지 테스트하는 함수
func Test_FilenameToTags(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{{
		in:   "s0010_c0010_ani_v001", // _ 포함 경우
		want: []string{"s0010", "c0010", "ani", "v001"},
	}, {
		in:   "s0010_c0010_ani_v001.mb", // 확장자 포함 경우
		want: []string{"s0010", "c0010", "ani", "v001"},
	}, {
		in:   "test_thumbnail.mb", // thumbnail 포함 경우
		want: []string{"test"},
	}, {
		in:   "ani/v001/test", // / 포함 경우
		want: []string{"ani", "v001", "test"},
	}, {
		in:   "rigging,shader", // , 포함 경우
		want: []string{"rigging", "shader"},
	}, {
		in:   "split by sign", // 공백 포함 경우
		want: []string{"split", "by", "sign"},
	}, {
		in:   "test", // 특수문자 포함 안한 경우
		want: []string{"test"},
	},
	}
	for _, c := range cases {
		b := FilenameToTags(c.in)
		if !reflect.DeepEqual(b, c.want) {
			t.Fatalf("Test_SplitbySign(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.in, c.want, b)
		}
	}
}
