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
		path: "/바탕화면/LIBRARY_3D/asset/", // 경로에 한글이 섞인경우
		want: true,
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
		path: "/LIBRARY 3D/asset/", // 경로에 공백문자가 들어간경우.
		want: true,
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

// IP형식을 테스트하기 위한 함수
func Test_checkIp(t *testing.T) {
	cases := []struct {
		ip   string
		want bool
	}{{
		ip:   "10.20.30.230",
		want: true,
	}, {
		ip:   "255.255.255.255",
		want: true,
	}, {
		ip:   "256.0.0.1", // 0~255 범위를 넘겼을 경우
		want: false,
	}, {
		ip:   "10.20.30", // -.-.-.- 형식이 아닌 경우
		want: false,
	}, {
		ip:   "...", // 숫자를 입력하지 않고 .만 찍었을 경우
		want: false,
	}, {
		ip:   "10.20.30.", // 마지막에 숫자를 입력하지 않은 경우
		want: false,
	},
	}

	for _, c := range cases {
		b := regexIPv4.MatchString(c.ip)
		if c.want != b {
			t.Fatalf("Test_checkIp(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.ip, c.want, b)
		}
	}
}

// Itemtype을 테스트하기위한 함수
func Test_Itemtype(t *testing.T) {
	cases := []struct {
		Itemtype string
		want     bool
	}{{
		Itemtype: "maya", // 정상 영문
		want:     true,
	}, {
		Itemtype: "ma야", // 한글포함
		want:     false,
	}, {
		Itemtype: "fusion360", // 숫자포함
		want:     true,
	},
	}

	for _, c := range cases {
		b := regexLower.MatchString(c.Itemtype)
		if c.want != b {
			t.Fatalf("Test_Itemtype(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Itemtype, c.want, b)
		}
	}
}

// Dbname을 테스트하기위한 함수
func Test_Dbname(t *testing.T) {
	cases := []struct {
		Dbname string
		want   bool
	}{{
		Dbname: "dotori", // 정상 영문
		want:   true,
	}, {
		Dbname: "doto리", // 한글포함
		want:   false,
	}, {
		Dbname: "fusion360", // 숫자포함
		want:   true,
	},
	}

	for _, c := range cases {
		b := regexLower.MatchString(c.Dbname)
		if c.want != b {
			t.Fatalf("Test_Dbname(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Dbname, c.want, b)
		}
	}
}

// idToPath를 테스트하기위한 함수
func Test_idToPath(t *testing.T) {
	cases := []struct {
		Dbname string
		want   bool
	}{{
		Dbname: "54759eb3c090d83494e2d804", // 정상 소문자와 숫자포함
		want:   true,
	}, {
		Dbname: "129638926139621982386219", // 정상 숫자
		want:   true,
	}, {
		Dbname: "oweiruioqrjkldafieuqwrri", // 정상 소문자
		want:   true,
	}, {
		Dbname: "54759eb3c090d83494e2d8/4", // 특수문자"/" 포함
		want:   false,
	}, {
		Dbname: "54759eB3c090d83494E2d804", // 대문자 포함
		want:   false,
	}, {
		Dbname: "54759eb3c090d83494e2d.04", // 특수문자"." 포함
		want:   false,
	}, {
		Dbname: "54759eb3c090d83494 e2d804", // 공백 포함
		want:   false,
	},
	}

	for _, c := range cases {
		b := regexObjectID.MatchString(c.Dbname)
		if c.want != b {
			t.Fatalf("Test_Dbname(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Dbname, c.want, b)
		}
	}
}

// Attributes을 테스트하기위한 함수
func Test_Attributes(t *testing.T) {
	cases := []struct {
		Attributes string
		want       bool
	}{{
		Attributes: "Nuke:10", // 값이 하나만 있는 경우
		want:       true,
	}, {
		Attributes: "Nuke:10,size:2048x858", // 값이 두개인 경우
		want:       true,
	}, {
		Attributes: "Nuke:10,", // , 뒤에 값이 없는 경우
		want:       false,
	}, {
		Attributes: "Nuke:10,size:", // : 뒤에 값이 없는 경우
		want:       false,
	}, {
		Attributes: "Nuke:10,size", // key까지만 입력한 경우
		want:       false,
	}, {
		Attributes: "Nuke,size", // key만 입력한 경우
		want:       false,
	}, {
		Attributes: ":10,:2048x858", // value만 입력한 경우
		want:       false,
	}, {
		Attributes: "Nuke:10,,size:2048x858", // 중간에 값이 없는 경우
		want:       false,
	}, {
		Attributes: "Nuke:10.5,size:2048x858", // value에 특수문자가 포함된 경우
		want:       true,
	}, {
		Attributes: ",Nuke:10", // ,으로 시작한 경우
		want:       false,
	},
	}

	for _, c := range cases {
		b := regexMap.MatchString(c.Attributes)
		if c.want != b {
			t.Fatalf("Test_Attributes(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Attributes, c.want, b)
		}
	}
}

// tag를 테스트하기위한 함수
func Test_Tag(t *testing.T) {
	cases := []struct {
		Tag  string
		want bool
	}{{
		Tag:  "tag", // 영문만 있는 경우
		want: true,
	}, {
		Tag:  "태그", // 한글만 있는 경우
		want: true,
	}, {
		Tag:  "tag태그", // 영문, 한글이 포함된 경우
		want: true,
	}, {
		Tag:  "비밀유지서약서", // 긴 길이의 태그
		want: true,
	}, {
		Tag:  "tag2", // 숫자가 포함된 경우
		want: true,
	}, {
		Tag:  "tag@", // 특수문자가 포함된 경우
		want: false,
	}, {
		Tag:  ",tag", // 특수문자로 시작하는 경우
		want: false,
	},
	}

	for _, c := range cases {
		b := regexTag.MatchString(c.Tag)
		if c.want != b {
			t.Fatalf("Test_Tag(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Tag, c.want, b)
		}
	}
}

//permission을 테스트하기 위한 함수
func Test_Permission(t *testing.T) {
	cases := []struct {
		Perm string
		want bool
	}{{
		Perm: "0777",
		want: true,
	}, {
		Perm: "0778", //permission의 범위를 넘어가는 수가 포함된 경우
		want: false,
	}, {
		Perm: "07a7", //수가 아닌 문자가 포함된 경우 / 영문
		want: false,
	}, {
		Perm: "07가7", //수가 아닌 문자가 포함된 경우 / 한글
		want: false,
	}, {
		Perm: "3777", // 첫번째 자리에 0이 아닌 다른 수가 오는 경우
		want: false,
	}, {
		Perm: "07777", // 4자리를 넘는 경우
		want: false,
	},
	}
	for _, c := range cases {
		b := regexPermission.MatchString(c.Perm)
		if c.want != b {
			t.Fatalf("Test_Permission(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Perm, c.want, b)
		}
	}
}

//SpliBySign 정규표현식을 테스트하기 위한 함수
func Test_SplitBySign(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "s0010_c0010_fx_v021.mb", // _ 포함
		want: true,
	}, {
		in:   "test,code,split", // , 포함
		want: true,
	}, {
		in:   "model geo", // 공백 포함
		want: true,
	}, {
		in:   "pr/test/rigged", // / 포함
		want: true,
	}, {
		in:   "lightingtest.mb", // . 포함
		want: false,
	}, {
		in:   "lightingte<stmb", // < 포함
		want: false,
	},
	}
	for _, c := range cases {
		b := regexSplitBySign.MatchString(c.in)
		if c.want != b {
			t.Fatalf("Test_SplitBySign(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.in, c.want, b)
		}
	}
}

// title를 테스트하기위한 함수
func Test_Title(t *testing.T) {
	cases := []struct {
		Title string
		want  bool
	}{{
		Title: "title", // 영문만 있는 경우
		want:  true,
	}, {
		Title: "123", // 숫자만 있는 경우
		want:  true,
	}, {
		Title: "title2", // 영문 + 숫자 조합
		want:  true,
	}, {
		Title: "타이틀", // 한글만 있는 경우
		want:  true,
	}, {
		Title: "타이틀 생성", // 공백 포함
		want:  true,
	}, {
		Title: "title타이틀", // 영문 + 한글 조합
		want:  true,
	}, {
		Title: "title@", // 특수문자가 포함된 경우
		want:  false,
	}, {
		Title: ",title", // 특수문자로 시작하는 경우
		want:  false,
	},
	}

	for _, c := range cases {
		b := regexTitle.MatchString(c.Title)
		if c.want != b {
			t.Fatalf("Test_Title(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.Title, c.want, b)
		}
	}
}

func Test_Str2List(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{{
		input: "title", // 태그가 하나일 때
		want:  []string{"title"},
	}, {
		input: "tag1 tag2", // 태그가 띄어쓰기 일 때
		want:  []string{"tag1", "tag2"},
	}, {
		input: "tag1,tag2", // 쉼표로 구분되어 있을 때
		want:  []string{"tag1", "tag2"},
	}, {
		input: "tag1, tag2", // 쉼표와 띄어쓰기로 구분되어 있을 때
		want:  []string{"tag1", "tag2"},
	}, {
		input: "tag1, tag2 tag3", // 쉼표와 띄어쓰기가 다른곳에서 구분되어 있을 때
		want:  []string{"tag1", "tag2", "tag3"},
	}, {
		input: "tag1, tag2     tag3", // 다수의 스페이스가 포함된 경우
		want:  []string{"tag1", "tag2", "tag3"},
	}, {
		input: "tag1, tag2     		tag3 태그", // 다수의 스페이스와 탭이 섞여있는 경우
		want: []string{"tag1", "tag2", "tag3", "태그"},
	},
	}
	for _, c := range cases {
		v := Str2List(c.input)
		if !testIsEqualSlice(c.want, v) {
			t.Fatalf("Test_Str2List(): 입력 값: %v, 원하는 값: %v, 얻은 값: %v\n", c.input, c.want, v)
		}
	}
}

// testIsEqualSlice 함수는 2개의 리스트를 받아서 구조가 같은지 체크한다. 이 함수는 test 되면 안되기 때문에 소문자 test로 시작한다.
func testIsEqualSlice(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
