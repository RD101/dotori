package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/sys/unix"
)

// searchSeq 함수는 탐색할 경로를 입력받고 dpx, exr, png, mov 정보를 수집 반환한다.
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
		case ".mov", ".mp4":
			item := Seq{
				Searchpath: searchpath,
				Dir:        filepath.Dir(path),
				Base:       filepath.Base(path),
				Ext:        ext,
				ConvertExt: ".exr",
			}
			paths[path] = item
		case ".dpx", ".exr", ".png":
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
	return header + "%0" + strconv.Itoa(len(seq)) + "d" + ext, seqNum, nil
}

// idToPath 함수는 MongoDB ID를 받아서 정한 형식에 맞게 ID를 변경시켜준다.
// 용도 : 몽고디비에서 생성되는 고유아이디로 폴더구조를 생성하여 각 유저마다 해당 에셋에 대한 데이터를 쌓아주기 위함이다.
// 나누는 이유 : 폴더에 저장할 수 있는 파일의 개수는 한정적이기 때문에 파일이 몰리지 않도록 분산해주기 위함이다.
// "54759eb3c090d83494e2d804" -> “/54/75/9e/b3/c090d8/3494/e2/d8/04”
func idToPath(id string) (string, error) {
	if len(id) != 24 {
		return id, errors.New("MongoDB ID 형식이 아닙니다")
	}

	// 영문 소문자와 숫자만 허용
	if !regexObjectID.MatchString(id) {
		return id, errors.New("MongoDB ID 형식이 아닙니다")
	}

	// 형식에 맞게 "/" 추가 (2/2/2/2/6/4/2/2/2)
	result := fmt.Sprintf("/%s/%s/%s/%s/%s/%s/%s/%s/%s", id[0:2], id[2:4], id[4:6], id[6:8], id[8:14], id[14:18], id[18:20], id[20:22], id[22:24])
	return result, nil
}

// GetRootPath 함수는 Admin setting에서 설정한 Rootpath를 가져온다
func GetRootPath(client *mongo.Client) (string, error) {
	rootpath := ""
	//adminSetting에서 rootpath를 가져온다.
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		return rootpath, err
	}
	rootpath = adminsetting.Rootpath
	// rootpath가 빈문자열이면
	if rootpath == "" {
		return rootpath, errors.New("admin setting에서 rootpath를 설정해주세요")
	}
	// rootpath가 '/'로 시작하지 않으면 앞에 슬래시를 붙여준다.
	if rootpath[0] != '/' {
		rootpath = "/" + rootpath
	}
	return rootpath, nil
}

// RmData 함수는 받아온 item id에 해당하는 데이터를 폴더 트리에서 삭제한다
func RmData(client *mongo.Client, id string) error {
	// get path
	rootpath, err := GetRootPath(client)
	if err != nil {
		return errors.New("admin setting에서 rootpath를 가져오지 못했습니다")
	}
	idpath, err := idToPath(id)
	if err != nil {
		return errors.New("id를 경로 형식으로 변환하지 못했습니다")
	}
	rmpath := rootpath + idpath
	splitpath, _ := path.Split(rmpath)
	if _, err := os.Stat(splitpath); os.IsNotExist(err) {
		return nil // 삭제할 경로가 존재하지 않는것은 에러가 아니다.
	}
	// 데이터가 존재한다. 데이터를 삭제한다.
	err = os.RemoveAll(rmpath) // idpath와 정확히 일치하는 최하단 경로만 강제로 삭제
	if err != nil {
		return err
	}
	for {
		splitpath, _ = path.Split(splitpath)
		splitpath = strings.TrimSuffix(splitpath, "/")
		c, err := ioutil.ReadDir(splitpath)
		if err != nil {
			return err
		}
		// 하위 폴더 없으면 삭제
		if len(c) == 0 {
			err := os.Remove(splitpath)
			if err != nil {
				return err
			}
		} else {
			break
		}
		if splitpath == rootpath {
			break
		}
	}
	return nil
}

// 참고한 코드: https://stackoverflow.com/a/21067803
// copyFile 함수는 inputpath경로의 파일을 outputpath로 복사한다.
func copyFile(inputpath, outputpath string) error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	// adminsetting을 가져온다,
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		return err
	}
	// adminsetting에서 폴더 생성에 필요한 값들을 가져온다.
	// umask, 권한 셋팅
	umask, err := strconv.Atoi(adminsetting.Umask)
	if err != nil {
		return err
	}
	unix.Umask(umask)
	folderP := adminsetting.FolderPermission
	folderPerm, err := strconv.ParseInt(folderP, 8, 64)
	if err != nil {
		return err
	}
	u := adminsetting.UID
	uid, err := strconv.Atoi(u)
	if err != nil {
		return err
	}
	g := adminsetting.GID
	gid, err := strconv.Atoi(g)
	if err != nil {
		return err
	}

	// input경로 검사
	src, err := os.Stat(inputpath)
	if err != nil {
		return err
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	if !src.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: 폴더, 심볼릭 링크 등은 복사할 수 없습니다. non-regular source file %s (%q)", src.Name(), src.Mode().String())
	}

	// output경로 검사.
	dst, err := os.Stat(outputpath)
	// 경로가 존재하지 않으면 새로 만든다.
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputpath, os.FileMode(folderPerm))
		if err != nil {
			return err
		}
		err = os.Chown(outputpath, uid, gid)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	_, filename := path.Split(inputpath)
	outputpath = outputpath + filename
	// src경로와 dst경로가 같으면 옮길 필요가 없다.
	if os.SameFile(src, dst) {
		return nil
	}
	err = copyFileContents(inputpath, outputpath)
	if err != nil {
		return err
	}
	return nil
}

func copyFileContents(inputpath, outputpath string) error {
	in, err := os.Open(inputpath)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(outputpath)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}
	return nil
}

// HasWildcard 함수는 경로를 받아서 Wildcard를 포함한다면 true, Wildcard를 포함하지 않는다면 false를 반환한다.
func HasWildcard(path string) bool {
	if strings.Contains(path, "*") || strings.Contains(path, "?") {
		return true
	}
	return false
}

// getFilesFromPath 함수는 입력받은 경로에 존재하는 파일들의 이름을 반환한다.
func getFilesFromPath(rootpath string) ([]string, error) {
	// 경로가 존재하는지 체크한다.
	_, err := os.Stat(rootpath)
	if err != nil {
		return nil, err
	}
	var result []string
	files, err := ioutil.ReadDir(rootpath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// 숨김파일
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		result = append(result, file.Name())
	}

	if len(result) == 0 {
		return nil, errors.New("파일이 존재하지 않습니다")
	}
	return result, nil
}
