package main

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

// Adminsetting 자료구조
type Adminsetting struct {
	ID                  string `json:"id" bson:"id"`                                   // DB에서 값을 가지고 오기 위한 id
	Rootpath            string `json:"rootpath" bson:"rootpath"`                       // 에셋 라이브러러리 물리경로
	LinuxProtocolPath   string `json:"linuxprotocolpath" bson:"linuxprotocolpath"`     // 웹에서 클릭시 사용하는 Linux 경로
	WindowsProtocolPath string `json:"windowsprotocolpath" bson:"windowsprotocolpath"` // 웹에서 클릭시 사용하는 Windows 경로
	MacosProtocolPath   string `json:"macosprotocolpath" bson:"macosprotocolpath"`     // 웹에서 클릭시 사용하는 macOS 경로
	Umask               string `json:"umask" bson:"umask"`                             // Umask 값
	FolderPermission    string `json:"folderpermission" bson:"folderpermission"`       // 폴더 생성시 사용하는 권한
	FilePermission      string `json:"filepermission" bson:"filepermission"`           // 파일 생성시 사용하는 권한
	UID                 string `json:"uid" bson:"uid"`                                 // 유저 ID
	GID                 string `json:"gid" bson:"gid"`                                 // 그룹 ID
	FFmpeg              string `json:"ffmpeg" bson:"ffmpeg"`                           // FFmpeg 명령어 경로
	OCIOConfig          string `json:"ocioconfig" bson:"ocioconfig"`                   // ocio.config 경로
	OpenImageIO         string `json:"openimageio" bson:"openimageio"`                 // OpenImageIO 명령어 경로
}

// Item 은 라이브러리의 에셋 자료구조이다.
type Item struct {
	ID          bson.ObjectId     `json:"id" bson:"_id,omitempty"`        // ID
	Author      string            `json:"author" bson:"author"`           // 에셋을 제작한 사람
	Tags        []string          `json:"tags" bson:"tags"`               // 태그리스트
	Description string            `json:"description" bson:"description"` // 에셋에 대한 추가 정보. 에셋의 제약, 사용전 알아야 할 특징
	Thumbimg    string            `json:"thumbimg" bson:"thumbimg"`       // 썸네일 이미지 주소
	Thumbmov    string            `json:"thumbmov" bson:"thumbmov"`       // 썸네일 영상 주소
	Inputpath   string            `json:"inputpath" bson:"inputpath"`     // 최초 등록되는 경로
	Outputpath  string            `json:"outputpath" bson:"outputpath"`   // 저장되는 경로
	ItemType    string            `json:"itemtype" bson:"itemtype"`       // maya, source, houdini, blender, nuke ..  같은 형태인가.
	Status      ItemStatus        `json:"status" bson:"status"`           // 상태(에러, done, wip)
	Log         string            `json:"log" bson:"log"`                 // 데이터를 처리할 때 생성되는 로그
	CreateTime  string            `json:"createtime" bson:"createtime"`   // Item 생성 시간
	Updatetime  string            `json:"updatetime" bson:"updatetime"`   // UTC 타임으로 들어가도록 하기.
	UsingRate   int64             `json:"usingrate" bson:"usingrate"`     // 사용 빈도 수
	Attributes  map[string]string `json:"attributes" bson:"attributes"`   // 해상도, 속성, 메타데이터 등의 파일정보
}

// ItemStatus 는 숫자이다.
type ItemStatus int

// item의 상태
const (
	Ready             = ItemStatus(iota) // 0 복사전
	Copying                              // 1 복사중
	Copied                               // 2 복사 완료
	CreatingThumbnail                    // 3 썸네일 생성중
	CreatedThumbnail                     // 4 썸네일 생성완료
	CreatingContainer                    // 5 썸네일 동영상 생성중
	CreatedContainer                     // 6 썸네일 동영상 생성완료
	Done                                 // 7 등록 완료
)

// CheckError 는 Item 자료구조에 값이 정확히 들어갔는지 확인하는 메소드이다.
func (i Item) CheckError() error {
	// 테스트를 하기 위해 임시로 넣어놓은 값. 나중에 제거해야 한다.
	i.CreateTime = "2019-09-09T14:43:34+09:00"
	i.Updatetime = "2019-09-09T14:43:34+09:00"

	if i.ItemType == "" {
		return errors.New("itemtype을 입력해주세요")
	}
	if !regexRFC3339Time.MatchString(i.CreateTime) {
		return errors.New("생성시간이 2019-09-09T14:43:34+09:00 형식의 문자열이 아닙니다")
	}
	if !regexRFC3339Time.MatchString(i.Updatetime) {
		return errors.New("업데이트 시간이 2019-09-09T14:43:34+09:00 형식의 문자열이 아닙니다")
	}
	if !regexPath.MatchString(i.Inputpath) {
		return errors.New("최초 등록 경로가 /test/test 형식의 문자열이 아닙니다")
	}
	if !regexPath.MatchString(i.Outputpath) {
		return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	}
	if !regexLower.MatchString(i.ItemType) {
		return errors.New("type이 소문자가 아닙니다")
	}
	for _, tag := range i.Tags {
		if !regexTag.MatchString(tag) {
			return errors.New("tag에는 특수문자를 사용할 수 없습니다")
		}
		if len(tag) == 1 {
			return errors.New("tag에는 한자리의 단어를 사용할 수 없습니다")
		}
	}
	return nil
}
