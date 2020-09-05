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

func Test_SingleQuotePath(t *testing.T) {
	cases := []struct {
		Itemtype string
		want     bool
	}{{
		Itemtype: "'/show/test'", // 작은 따옴표가 들어간 경로
		want:     true,
	}, {
		Itemtype: "/show/test", // 일반 문자열 경로
		want:     false,
	}, {
		Itemtype: "'/show/test df'", // 띄어쓰기 포함
		want:     true,
	}, {
		Itemtype: "'/show/project/assets/현장 데이터 사진/20200830'", // 띄어쓰기 + 한글 포함
		want:     true,
	},
	}
	for _, c := range cases {
		b := regexSingleQuotesPath.MatchString(c.Itemtype)
		if c.want != b {
			t.Fatalf("Test_SingleQuotePath(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Itemtype, c.want, b)
		}
	}
}

func Test_DoubleQuotePath(t *testing.T) {
	cases := []struct {
		Itemtype string
		want     bool
	}{{
		Itemtype: "\"/show/test\"", // 큰 따옴표가 들어간 경로
		want:     true,
	}, {
		Itemtype: "/show/test", // 일반 문자열 경로
		want:     false,
	}, {
		Itemtype: "\"/show/test df\"", // 띄어쓰기 포함
		want:     true,
	}, {
		Itemtype: "\"/show/project/assets/현장 데이터 사진/20200830\"", // 띄어쓰기 및 한글포함
		want:     true,
	},
	}
	for _, c := range cases {
		b := regexDoubleQuotesPath.MatchString(c.Itemtype)
		if c.want != b {
			t.Fatalf("Test_DoubleQuotePath(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Itemtype, c.want, b)
		}
	}
}

func Test_HasQuotes(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "'/project/test.*.exr'", // 작은 따옴표 포함
		want: true,
	}, {
		in:   "\"/project/test.????.exr\"", // 큰 따옴표 포함
		want: true,
	}, {
		in:   "/project/test.exr", // 아무것도 없는경우
		want: false,
	},
	}
	for _, c := range cases {
		b := HasQuotes(c.in)
		if c.want != b {
			t.Fatalf("Test_HasQuotes(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.in, c.want, b)
		}
	}
}

func Test_QuotesPaths2Paths(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{{
		in:   "'/project/test.0010.exr'", // 작은 따옴표 포함
		want: []string{"/project/test.0010.exr"},
	}, {
		in:   "\"/project/dublequotes.0010.exr\"", // 큰 따옴표 포함
		want: []string{"/project/dublequotes.0010.exr"},
	}, {
		in:   "/project/test.exr", // 아무것도 없는경우
		want: []string{"/project/test.exr"},
	}, {
		in:   "\"/project/case1.exr\" \"/project/case2.exr\"", // 큰 따옴표로 구성된 다중경로
		want: []string{"/project/case1.exr", "/project/case2.exr"},
	}, {
		in:   "\"/project/case 1.exr\" \"/project/case2.exr\"", // 큰 따옴표 + 스페이스로 구성된 다중경로
		want: []string{"/project/case 1.exr", "/project/case2.exr"},
	}, {
		in:   "'/project/case 1.exr' \"/project/case2.exr\"", // 작은 따옴표 + 큰 따옴표 + 스페이스로 구성된 다중경로
		want: []string{"/project/case 1.exr", "/project/case2.exr"},
	}, {
		in:   "/project/case1.exr /project/case2.exr", // 띄어쓰기로 구성된 다중경로
		want: []string{"/project/case1.exr", "/project/case2.exr"},
	},
	}
	for _, c := range cases {
		results := QuotesPaths2Paths(c.in)
		if !testIsEqualSlice(c.want, results) {
			t.Fatalf("Test_HasQuotes(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.in, c.want, results)
		}
	}
}
