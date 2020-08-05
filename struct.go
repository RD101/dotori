package main

import (
	"errors"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Adminsetting 자료구조
type Adminsetting struct {
	ID                      string `json:"id" bson:"id"`                                           // DB에서 값을 가지고 오기 위한 id
	Appname                 string `json:"appname" bson:"appname"`                                 // 어플리케이션 표기이름
	Rootpath                string `json:"rootpath" bson:"rootpath"`                               // 에셋 라이브러러리 물리경로
	LinuxProtocolPath       string `json:"linuxprotocolpath" bson:"linuxprotocolpath"`             // 웹에서 클릭시 사용하는 Linux 경로
	WindowsProtocolPath     string `json:"windowsprotocolpath" bson:"windowsprotocolpath"`         // 웹에서 클릭시 사용하는 Windows 경로
	MacosProtocolPath       string `json:"macosprotocolpath" bson:"macosprotocolpath"`             // 웹에서 클릭시 사용하는 macOS 경로
	Umask                   string `json:"umask" bson:"umask"`                                     // Umask 값
	FolderPermission        string `json:"folderpermission" bson:"folderpermission"`               // 폴더 생성시 사용하는 권한
	FilePermission          string `json:"filepermission" bson:"filepermission"`                   // 파일 생성시 사용하는 권한
	UID                     string `json:"uid" bson:"uid"`                                         // 유저 ID
	GID                     string `json:"gid" bson:"gid"`                                         // 그룹 ID
	ProcessBufferSize       int    `json:"processbuffersize" bson:"processbuffersize"`             // 프로세스할 아이템을 담아둘 버퍼의 사이즈. 회사 인원 수 만큼을 추천
	FFmpeg                  string `json:"ffmpeg" bson:"ffmpeg"`                                   // FFmpeg 명령어 경로
	OCIOConfig              string `json:"ocioconfig" bson:"ocioconfig"`                           // ocio.config 경로
	OpenImageIO             string `json:"openimageio" bson:"openimageio"`                         // OpenImageIO 명령어 경로
	LDLibraryPath           string `json:"ldlibrarypath"`                                          // LD_LIBRARY_PATH
	MultipartFormBufferSize int    `json:"multipartformbuffersize" bson:"multipartformbuffersize"` // MultipartForm Buffersize
	ThumbnailImageWidth     int    `json:"thumbnailimagewidth" bson:"thumbnailimagewidth"`         // 썸네일 이미지 가로 픽셀 사이즈
	ThumbnailImageHeight    int    `json:"thumbnailimageheight" bson:"thumbnailimageheight"`       // 썸네일 이미지 세로 픽셀 사이즈
	MediaWidth              int    `json:"mediawidth" bson:"mediawidth"`                           // 동영상 가로 픽셀 사이즈
	MediaHeight             int    `json:"mediaheight" bson:"mediaheight"`                         // 동영상 세로 픽셀 사이즈
	VideoCodecOgg           string `json:"videocodecogg" bson:"videocodecogg"`                     // 비디오 코덱
	VideoCodecMp4           string `json:"videocodecmp4" bson:"videocodecmp4"`                     // 비디오 코덱
	VideoCodecMov           string `json:"videocodecmov" bson:"videocodecmov"`                     // 비디오 코덱
	AudioCodec              string `json:"audiocodec" bson:"audiocodec"`                           // 오디오 코덱
}

// Token 자료구조. JWT 방식을 사용한다. restAPI 사용시 보안체크를 위해 http 헤더에 들어간다.
type Token struct {
	ID          string `json:"id" bson:"id"`                   // 사용자 ID
	AccessLevel string `json:"accesslevel" bson:"accesslevel"` // admin, manager, default
	jwt.StandardClaims
}

// User 는 사용자 자료구조이다.
type User struct {
	ID          string `json:"id" bson:"id"`                   // 사용자 ID
	Password    string `json:"password" bson:"password"`       // 암호화된 암호
	Token       string `json:"token" bson:"token"`             // JWT 토큰
	SignKey     string `json:"signkey" bson:"signkey"`         // JWT 토큰을 만들 때 사용하는 SignKey
	AccessLevel string `json:"accesslevel" bson:"accesslevel"` // admin, manager, default
}

// Item 은 라이브러리의 에셋 자료구조이다.
type Item struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`        // ID
	Author      string             `json:"author" bson:"author"`           // 에셋을 제작한 사람
	Title       string             `json:"title" bson:"title"`             // 에셋 타이틀
	Tags        []string           `json:"tags" bson:"tags"`               // 태그리스트
	Description string             `json:"description" bson:"description"` // 에셋에 대한 추가 정보. 에셋의 제약, 사용전 알아야 할 특징
	ItemType    string             `json:"itemtype" bson:"itemtype"`       // maya, source, houdini, blender, nuke ..  같은 형태인가.
	Status      string             `json:"status" bson:"status"`           // 상태: "error", "done", "wip" 등등
	Logs        []string           `json:"logs" bson:"logs"`               // 데이터를 처리할 때 생성되는 로그
	CreateTime  string             `json:"createtime" bson:"createtime"`   // Item 생성 시간
	Updatetime  string             `json:"updatetime" bson:"updatetime"`   // UTC 타임으로 들어가도록 하기.
	UsingRate   int64              `json:"usingrate" bson:"usingrate"`     // 사용 빈도 수
	Attributes  map[string]string  `json:"attributes" bson:"attributes"`   // 해상도, 속성, 메타데이터 등의 파일정보

	InputThumbnailImgPath  string `json:"inputthumbnailimgpath" bson:"inputthumbnailimgpath"`   // 사용자가 업로드한 썸네일 이미지의 업로드 경로
	InputThumbnailClipPath string `json:"inputthumbnailclippath" bson:"inputthumbnailclippath"` // 사용자가 업로드한 클립 파일의 업로드 경로
	OutputThumbnailPngPath string `json:"outputthumbnailpngpath" bson:"outputthumbnailpngpath"` // 생성된 썸네일 이미지를 저장하는 경로
	OutputThumbnailMp4Path string `json:"outputthumbnailmp4path" bson:"outputthumbnailmp4path"` // 생성된 mp4형식의 썸네일 클립을 저장하는 경로
	OutputThumbnailOggPath string `json:"outputthumbnailoggpath" bson:"outputthumbnailoggpath"` // 생성된 ogg형식의 썸네일 클립을 저장하는 경로
	OutputThumbnailMovPath string `json:"outputthumbnailmovpath" bson:"outputthumbnailmovpath"` // 생성된 mov형식의 썸네일 클립을 저장하는 경로
	OutputProxyImgPath     string `json:"outputproxyimgpath" bson:"outputproxyimgpath"`         // 프록시 이미지를 저장하는 경로
	OutputDataPath         string `json:"outputdatapath" bson:"outputdatapath"`                 // 사용자가 업로드한 파일 중 썸네일 이미지와 클립을 제외한 나머지 파일을 저장하는 경로

	ThumbImgUploaded  bool `json:"thumbimguploaded" bson:"thumbimguploaded"`   // 썸네일 이미지의 업로드 여부 체크
	ThumbClipUploaded bool `json:"thumbclipuploaded" bson:"thumbclipuploaded"` // 썸네일 클립의 업로드 여부 체크
	DataUploaded      bool `json:"datauploaded" bson:"datauploaded"`           // 데이터의 업로드 여부 체크

	InColorspace  string `json:"incolorspace" bson:"incolorspace"`   // InColorspace
	OutColorspace string `json:"outcolorspace" bson:"outcolorspace"` // OutColorspace
	Fps           string `json:"fps" bson:"fps"`                     // fps값 ffmpeg 연산에 사용되는 값이기 때문에 문자열로 처리함.

	KindOfUSD string `json:"kindofusd" bson:"kindofusd"` // Kind Of USD
}

// ItemStatus 는 숫자이다.
type ItemStatus int

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

// CheckError 는 Adminsetting 자료구조에 값이 정확히 들어갔는지 확인하는 메소드이다.
func (a Adminsetting) CheckError() error {
	if !regexPermission.MatchString(a.FolderPermission) {
		return errors.New("FolderPermission이 형식에 맞지 않습니다")
	}
	if !regexPermission.MatchString(a.FilePermission) {
		return errors.New("FilePermission이 형식에 맞지 않습니다")
	}
	if !regexPermission.MatchString(a.Umask) {
		return errors.New("Umask가 형식에 맞지 않습니다")
	}
	if a.FFmpeg == "" {
		return errors.New("FFmpeg 경로를 설정해주세요")
	}
	if _, err := os.Stat(a.FFmpeg); os.IsNotExist(err) {
		return errors.New(a.FFmpeg + " 경로에 FFmpeg 명령어가 존재하지 않습니다")
	}
	return nil
}

// CreateToken 메소드는 토큰을 생성합니다.
func (u *User) CreateToken() error {
	if u.ID == "" {
		return errors.New("ID is an empty string")
	}
	if u.Password == "" {
		return errors.New("Password is an empty string")
	}
	if u.AccessLevel == "" {
		return errors.New("AccessLevel is an empty string")
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{
		ID:          u.ID,
		AccessLevel: u.AccessLevel,
	})
	signKey, err := Encrypt(u.Password)
	if err != nil {
		return err
	}
	u.SignKey = signKey
	tokenString, err := token.SignedString([]byte(signKey))
	if err != nil {
		return err
	}
	u.Token = tokenString
	return nil
}

// ItemListLength 함수는 Item형 리스트 전체의 개수를 반환한다.
func ItemListLength(items []Item) int {
	return len(items)
}
