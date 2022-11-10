package main

import (
	"errors"
	"os"
)

// Adminsetting 자료구조
type Adminsetting struct {
	ID                      string `json:"id" bson:"id"`                                           // DB에서 값을 가지고 오기 위한 id
	Appname                 string `json:"appname" bson:"appname"`                                 // 어플리케이션 표기이름
	Rootpath                string `json:"rootpath" bson:"rootpath"`                               // 에셋 라이브러러리 물리경로
	LinuxProtocolPath       string `json:"linuxprotocolpath" bson:"linuxprotocolpath"`             // 웹에서 클릭시 사용하는 Linux 경로
	WindowsProtocolPath     string `json:"windowsprotocolpath" bson:"windowsprotocolpath"`         // 웹에서 클릭시 사용하는 Windows 경로
	MacosProtocolPath       string `json:"macosprotocolpath" bson:"macosprotocolpath"`             // 웹에서 클릭시 사용하는 macOS 경로
	WindowsUNCPrefix        string `json:"windowsuncprefix" bson:"windowsuncprefix"`               // CopyPath를 실행할 때 윈도우즈의 경우 경로 앞에 붙는 문자열
	Umask                   string `json:"umask" bson:"umask"`                                     // Umask 값
	FolderPermission        string `json:"folderpermission" bson:"folderpermission"`               // 폴더 생성시 사용하는 권한
	FilePermission          string `json:"filepermission" bson:"filepermission"`                   // 파일 생성시 사용하는 권한
	UID                     string `json:"uid" bson:"uid"`                                         // 유저 ID
	GID                     string `json:"gid" bson:"gid"`                                         // 그룹 ID
	ProcessBufferSize       int    `json:"processbuffersize" bson:"processbuffersize"`             // 프로세스할 아이템을 담아둘 버퍼의 사이즈. 회사 인원 수 만큼을 추천
	FFmpeg                  string `json:"ffmpeg" bson:"ffmpeg"`                                   // FFmpeg 명령어 경로
	OCIOConfig              string `json:"ocioconfig" bson:"ocioconfig"`                           // ocio.config 경로
	OpenImageIO             string `json:"openimageio" bson:"openimageio"`                         // OpenImageIO 명령어 경로
	Mongodump               string `json:"mongodump" bson:"mongodump"`                             // mongodump 경로
	Backuppath              string `json:"backuppath" bson:"backuppath"`                           // 백업경로
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
	InitPassword            string `json:"initpassword" bson:"initpassword"`                       // 사용자가 패스워드를 잃어버렸을 때 초기화하는 패스워드
	EnableRVLink            bool   `json:"enablervlink" bson:"enablervlink"`                       // RVLink 활성화
	EnableCategory          bool   `json:"enablercategory" bson:"enablecategory"`                  // 카테고리 활성화
	// Support Add menu
	Maya            bool `json:"maya" bson:"maya"`                       // Maya
	Max             bool `json:"max" bson:"max"`                         // Max
	Nuke            bool `json:"nuke" bson:"nuke"`                       // Nuke
	Houdini         bool `json:"houdini" bson:"houdini"`                 // Houdini
	Blender         bool `json:"blender" bson:"blender"`                 // Blender
	Footage         bool `json:"footage" bson:"footage"`                 // Footage
	MultipleFootage bool `json:"multiplefootage" bson:"multiplefootage"` // Multiple Footage
	Alembic         bool `json:"alembic" bson:"alembic"`                 // Alembic
	USD             bool `json:"usd" bson:"usd"`                         // Usd
	Unreal          bool `json:"unreal" bson:"unreal"`                   // Unreal
	OpenVDB         bool `json:"openvdb" bson:"openvdb"`                 // OpenVDB
	Sound           bool `json:"sound" bson:"sound"`                     // Sound
	Modo            bool `json:"modo" bson:"modo"`                       // Modo
	Katana          bool `json:"katana" bson:"katana"`                   // Katana
	HWP             bool `json:"hwp" bson:"hwp"`                         // HWP
	PDF             bool `json:"pdf" bson:"pdf"`                         // PDF
	PPT             bool `json:"ppt" bson:"ppt"`                         // PPT
	IES             bool `json:"ies" bson:"ies"`                         // IES
	LUT             bool `json:"lut" bson:"lut"`                         // LUT
	HDRI            bool `json:"hdri" bson:"hdri"`                       // HDRI
	Texture         bool `json:"texture" bson:"texture"`                 // Texture
	Clip            bool `json:"clip" bson:"clip"`                       // Clip
	MultipleClip    bool `json:"multipleclip" bson:"multipleclip"`       // Multiple Clip
	Fusion360       bool `json:"fusion360" bson:"fusion360"`             // Fusion360
	Matte           bool `json:"matte" bson:"matte"`                     // Matte
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
		return errors.New("umask가 형식에 맞지 않습니다")
	}
	if a.FFmpeg == "" {
		return errors.New("FFmpeg 경로를 설정해주세요")
	}
	if _, err := os.Stat(a.FFmpeg); os.IsNotExist(err) {
		return errors.New(a.FFmpeg + " 경로에 FFmpeg 명령어가 존재하지 않습니다")
	}
	return nil
}
