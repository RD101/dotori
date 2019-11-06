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

// path를 테스트하기 위한 함수
func Test_CheckPath(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{{
		path: "/LIBRARY_3D/asset/", // 정상 경로
		want: true,
	}, {
		path: "//LIBRARY_3D/asset/", // 맨앞에 '/'가 붙은 경우
		want: true,
	}, {
		path: "C:/LIBRARY_3D/asset/", // C드라이브에 있을 경우 (모든 경로는 삼바UNC 또는 유닉스 경로여야 한다.
		want: false,
	}, {
		path: "D:/LIBRARY_3D/asset/", // D드라이브로 시작하는 경로
		want: false,
	}, {
		path: "Q:/LIBRARY_3D/asset/", // Q드라이브로 시작하는 경로
		want: false,
	}, {
		path: "/바탕화면/LIBRARY_3D/asset/", // 경로에 한글이 섞이면 안된다.
		want: false,
	}, {
		path: "/library_3d/asset/", // 전부 소문자인 경우
		want: true,
	}, {
		path: "/LIBRARY_3D/ASSET/", // 전부 대문자인 경우
		want: true,
	}, {
		path: "/Library_3d/Asset/", // 맨앞만 대문자인 경우
		want: true,
	}, {
		path: "/LIBRARY_3D/asset/src.png", // 확장자가 png인 경우
		want: true,
	}, {
		path: "/LIBRARY_3D/asset/src.mp4", // 확장자가 mp4인 경우
		want: true,
	}, {
		path: "/LIBRARY_3D/asset/src", // src 경로
		want: true,
	}, {
		path: "/LIBRARY_3D/asset/.image.png", // 숨김파일 .image.png 경로의 경우
		want: true,
	}, {
		path: "/#LIBRARY_3D★/asset/", // 경로에 특수문자가 들어간경우
		want: false,
	}, {
		path: "/LIBRARY_3D學問/asset/", // 경로에 한문이 들어간경우.
		want: false,
	}, {
		path: "/LIBRARY 3D/asset/", // 경로에 공백문자가 들어간경우. 에러는 아니지만 추후 연산을 위해 파이프라인툴에서 공백을 허용하지 않는다.
		want: false,
	}, {
		path: "/LIBRARY.3D/asset/", // 경로에 '.' 문자가 들어간경우
		want: true,
	}, {
		path: "\\\\LIBRARY.3D\\asset", // \\문자로 시작하는 경로
		want: false,
	}, {
		path: "//10.0.20.30/library", // unc path
		want: true,
	}, {
		path: "10.0.20.30/library", // /문자로 시작하지 않을 경우
		want: false,
	},
	}

	for _, c := range cases {
		b := regexPath.MatchString(c.path)
		if c.want != b {
			t.Fatalf("Test_checkPath(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.path, c.want, b)
		}
	}
}
