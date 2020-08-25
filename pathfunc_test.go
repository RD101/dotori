package main

import (
	"testing"
)

//HasWildcard 함수가 잘 작동하는지 테스트하는 함수
func Test_HasWildcard(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "/project/test.*.exr", //* 포함 경우
		want: true,
	}, {
		in:   "/project/test.????.exr", // ? 포함 경우
		want: true,
	}, {
		in:   "/project/test.exr", // 아무것도 없는경우
		want: false,
	},
	}
	for _, c := range cases {
		b := HasWildcard(c.in)
		if c.want != b {
			t.Fatalf("Test_HasWildcard(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.in, c.want, b)
		}
	}
}
