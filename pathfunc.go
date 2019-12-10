package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if strings.HasPrefix(info.Name(), ".") { // 숨김폴더는 Skip 한다.
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
