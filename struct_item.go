package main

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Item 은 라이브러리의 에셋 자료구조이다.
type Item struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`                              // ID
	Author                 string             `json:"author" bson:"author"`                                 // 에셋을 제작한 사람
	Title                  string             `json:"title" bson:"title"`                                   // 에셋 타이틀
	Tags                   []string           `json:"tags" bson:"tags"`                                     // 태그리스트
	Description            string             `json:"description" bson:"description"`                       // 에셋에 대한 추가 정보. 에셋의 제약, 사용전 알아야 할 특징
	ItemType               string             `json:"itemtype" bson:"itemtype"`                             // maya, source, houdini, blender, nuke ..  같은 형태인가.
	Status                 string             `json:"status" bson:"status"`                                 // 상태: "error", "done", "wip" 등등
	Logs                   []string           `json:"logs" bson:"logs"`                                     // 데이터를 처리할 때 생성되는 로그
	CreateTime             string             `json:"createtime" bson:"createtime"`                         // Item 생성 시간
	Updatetime             string             `json:"updatetime" bson:"updatetime"`                         // UTC 타임으로 들어가도록 하기.
	UsingRate              int64              `json:"usingrate" bson:"usingrate"`                           // 사용 빈도 수
	Attributes             map[string]string  `json:"attributes" bson:"attributes"`                         // 해상도, 속성, 메타데이터 등의 파일정보
	InputThumbnailImgPath  string             `json:"inputthumbnailimgpath" bson:"inputthumbnailimgpath"`   // 사용자가 업로드한 썸네일 이미지의 업로드 경로
	InputThumbnailClipPath string             `json:"inputthumbnailclippath" bson:"inputthumbnailclippath"` // 사용자가 업로드한 클립 파일의 업로드 경로
	OutputThumbnailPngPath string             `json:"outputthumbnailpngpath" bson:"outputthumbnailpngpath"` // 생성된 썸네일 이미지를 저장하는 경로
	OutputThumbnailMp4Path string             `json:"outputthumbnailmp4path" bson:"outputthumbnailmp4path"` // 생성된 mp4형식의 썸네일 클립을 저장하는 경로
	OutputThumbnailOggPath string             `json:"outputthumbnailoggpath" bson:"outputthumbnailoggpath"` // 생성된 ogg형식의 썸네일 클립을 저장하는 경로
	OutputThumbnailMovPath string             `json:"outputthumbnailmovpath" bson:"outputthumbnailmovpath"` // 생성된 mov형식의 썸네일 클립을 저장하는 경로
	OutputProxyImgPath     string             `json:"outputproxyimgpath" bson:"outputproxyimgpath"`         // 프록시 이미지를 저장하는 경로
	OutputDataPath         string             `json:"outputdatapath" bson:"outputdatapath"`                 // 사용자가 업로드한 파일 중 썸네일 이미지와 클립을 제외한 나머지 파일을 저장하는 경로
	ThumbImgUploaded       bool               `json:"thumbimguploaded" bson:"thumbimguploaded"`             // 썸네일 이미지의 업로드 여부 체크
	ThumbClipUploaded      bool               `json:"thumbclipuploaded" bson:"thumbclipuploaded"`           // 썸네일 클립의 업로드 여부 체크
	DataUploaded           bool               `json:"datauploaded" bson:"datauploaded"`                     // 데이터의 업로드 여부 체크
	InColorspace           string             `json:"incolorspace" bson:"incolorspace"`                     // InColorspace
	OutColorspace          string             `json:"outcolorspace" bson:"outcolorspace"`                   // OutColorspace
	Fps                    string             `json:"fps" bson:"fps"`                                       // fps값 ffmpeg 연산에 사용되는 값이기 때문에 문자열로 처리함.
	Premultiply            bool               `json:"premultiply" bson:"premultiply"`                       // proxy sequence to video 연산 과정에서 Premultiply 적용 여부 체크
	KindOfUSD              string             `json:"kindofusd" bson:"kindofusd"`                           // Kind Of USD
	RequireCopyInProcess   bool               `json:"requirecopyinprocess" bson:"requirecopyinprocess"`     // Process 단계에서 InputData로 부터 데이터 카피에 대한 필요 여부(예)인트라넷 Footage)
	RequireMkdirInProcess  bool               `json:"requiremkdirinprocess" bson:"requiremkdirinprocess"`   // Process 단계에서 데이터 복사시 해당 id의 폴더를 생성할지 여부
	InputData              InputData          // 최초 소스 정보 (예)인트라넷 Footage)
	Categories             []string           `json:"categories" bson:"categories"` // 카테고리 리스트. []string{"rootcategory"}, []string{"rootcategory","subcategory"} 형태로 저장되는 구조이다.
}

type InputData struct {
	Dir      string `json:"dir" bson:"dir"`           // 시퀀스 디렉토리
	Base     string `json:"base" bson:"base"`         // 파일명(시퀀스 숫자 제외)
	FrameIn  int    `json:"framein" bson:"framein"`   // 시작프레임
	FrameOut int    `json:"frameout" bson:"frameout"` // 끝프레임
}

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
	// if !regexPath.MatchString(i.InputThumbnailImgPath) {
	// 	return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	// }
	// if !regexPath.MatchString(i.InputThumbnailClipPath) {
	// 	return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	// }
	// if !regexPath.MatchString(i.OutputThumbnailPngPath) {
	// 	return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	// }
	// if !regexPath.MatchString(i.OutputThumbnailMp4Path) {
	// 	return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	// }
	// if !regexPath.MatchString(i.OutputThumbnailOggPath) {
	// 	return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	// }
	// if !regexPath.MatchString(i.OutputThumbnailMovPath) {
	// 	return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	// }
	// if !regexPath.MatchString(i.OutputDataPath) {
	// 	return errors.New("asset 저장 경로가 /test/test 형식의 문자열이 아닙니다")
	// }
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
	if i.Title == "" {
		return errors.New("title을 입력해주세요")
	}
	if len(i.Title) == 1 {
		return errors.New("title에는 한 자리의 단어를 사용할 수 없습니다")
	}
	if !regexTitle.MatchString(i.Title) {
		return errors.New("title에는 특수문자를 사용할 수 없습니다")
	}

	return nil
}

// ItemsTagsDeduplication 함수는 아이템들의 태그들을 중복제거한 리스트를 반환한다.
func ItemsTagsDeduplication(items []Item) []string {
	keys := make(map[string]bool)
	filteredTag := []string{}
	for itemIndex := range items {
		for tagIndex := range items[itemIndex].Tags {
			tagValue := items[itemIndex].Tags[tagIndex]
			if _, saveValue := keys[tagValue]; !saveValue {
				keys[tagValue] = true
				filteredTag = append(filteredTag, items[itemIndex].Tags[tagIndex])
			}
		}
	}
	return filteredTag
}
