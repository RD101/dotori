package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2"
)

// searchSeq 함수는 탐색할 경로를 입력받고 dpx, exr, mov 정보를 수집 반환한다.
func searchSeq(searchpath string) ([]Seq, error) {
	// 경로가 존재하는지 체크한다.
	_, err := os.Stat(searchpath)
	if err != nil {
		return nil, err
	}
	paths := make(map[string]Seq)
	err = filepath.Walk(searchpath, func(path string, info os.FileInfo, err error) error {
		// 숨김폴더
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		// 숨김파일
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".mov":
			item := Seq{
				Searchpath: searchpath,
				Dir:        filepath.Dir(path),
				Base:       filepath.Base(path),
				Ext:        ext,
				ConvertExt: ".exr",
			}
			paths[path] = item
		case ".dpx", ".exr":
			key, num, err := Seqnum2Sharp(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				return nil
			}
			if _, has := paths[key]; has {
				// 이미 수집된 경로가 존재할 때 처리되는 코드
				item := paths[key]
				item.Length++
				item.FrameOut = num
				item.RenderOut = num
				paths[key] = item
			} else {
				// 이전에 수집된 경로가 존재하지 않으면 처리되는 코드
				item := Seq{
					Searchpath: searchpath,
					Dir:        filepath.Dir(path),
					Base:       filepath.Base(key),
					Ext:        ext,
					Length:     1,
					FrameIn:    num,
					RenderIn:   num,
					ConvertExt: ext,
				}
				paths[key] = item
			}
		default:
			return nil
		}
		return nil
	})
	var items []Seq
	for _, value := range paths {
		items = append(items, value)
	}
	if len(items) == 0 {
		return nil, errors.New("소스가 존재하지 않습니다")
	}
	return items, nil
}

// Seqnum2Sharp 함수는 경로와 파일명을 받아서 시퀀스부분을 #문자열로 바꾸고 시퀀스의 숫자를 int로 바꾼다.
// "test.0002.jpg" -> "test.####.jpg", 2, nil
func Seqnum2Sharp(filename string) (string, int, error) {
	re, err := regexp.Compile("([0-9]+)(\\.[a-zA-Z]+$)")
	// 이 정보를 통해서 파일명을 구하는 방식으로 바꾼다.
	if err != nil {
		return filename, -1, errors.New("정규 표현식이 잘못되었습니다")
	}
	results := re.FindStringSubmatch(filename)
	if results == nil {
		return filename, -1, errors.New("경로가 시퀀스 형식이 아닙니다")
	}
	seq := results[1]
	ext := results[2]
	header := filename[:strings.LastIndex(filename, seq+ext)]
	seqNum, err := strconv.Atoi(seq)
	if err != nil {
		return filename, -1, err
	}
	return header + strings.Repeat("#", len(seq)) + ext, seqNum, nil
}

// idToPath 함수는 MongoDB ID를 받아서 정한 형식에 맞게 ID를 변경시켜준다.
// 용도 : 몽고디비에서 생성되는 고유아이디로 폴더구조를 생성하여 각 유저마다 해당 에셋에 대한 데이터를 쌓아주기 위함이다.
// 나누는 이유 : 폴더에 저장할 수 있는 파일의 개수는 한정적이기 때문에 파일이 몰리지 않도록 분산해주기 위함이다.
// "54759eb3c090d83494e2d804" -> “54/75/9e/b3/c090d8/3494/e2/d8/04”
func idToPath(id string) (string, error) {
	if len(id) != 24 {
		return id, errors.New("MongoDB ID 형식이 아닙니다")
	}

	// 영문 소문자와 숫자만 허용
	if !regexObjectID.MatchString(id) {
		return id, errors.New("MongoDB ID 형식이 아닙니다")
	}

	// 형식에 맞게 "/" 추가 (2/2/2/2/6/4/2/2/2)
	result := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s/%s/%s", id[0:2], id[2:4], id[4:6], id[6:8], id[8:14], id[14:18], id[18:20], id[20:22], id[22:24])
	return result, nil
}

// GetRootPath 함수는 Admin setting에서 설정한 Rootpath를 가져온다
func GetRootPath(session *mgo.Session) (string, error) {
	rootpath := ""
	//adminSetting에서 rootpath를 가져온다.
	adminsetting, err := GetAdminSetting(session)
	if err != nil {
		return rootpath, err
	}
	rootpath = adminsetting.Rootpath
	// rootpath가 빈문자열이면
	if rootpath == "" {
		return rootpath, errors.New("admin setting에서 rootpath를 설정해주세요")
	}
	// rootpath가 '/'로 끝나지 않으면 끝에 슬래시를 붙여준다.
	if rootpath[len(rootpath)-1] != '/' {
		rootpath = rootpath + "/"
	}
	return rootpath, nil
}
