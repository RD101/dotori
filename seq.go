package main

// Seq 자료구조는 시퀀스를 분석할 때 사용하는 자료구조이다.
type Seq struct {
	// 사용자 입력값
	InColorspace  string `json:"incolorspace" bson:"incolorspace"`   // 시퀀스의 IN 컬러스페이스
	OutColorspace string `json:"outcolorspace" bson:"outcolorspace"` // 시퀀스의 OUT 컬러스페이스
	FrameInPoint  int    `json:"frameinpoint" bson:"frameinpoint"`   // 사용자 시작 IN 프레임
	FrameOutPoint int    `json:"frameoutpoint" bson:"frameoutpoint"` // 사용자 끝 OUT 프레임

	// 분석을 통해서 구할 수 있는 것
	Dir         string  `json:"dir" bson:"dir"`                 // 시퀀스 디렉토리
	Base        string  `json:"base" bson:"base"`               // 파일명(시퀀스 숫자 제외)
	Ext         string  `json:"ext" bson:"ext"`                 // 확장자
	Digitnum    int     `json:"digitnum" bson:"digitnum"`       // 시퀀스 자릿수
	StartFrame  int     `json:"startframe" bson:"startframe"`   // 시작프레임
	EndFrame    int     `json:"endframe" bson:"endframe"`       // 끝프레임
	Width       int     `json:"width" bson:"width"`             // 가로길이
	Height      int     `json:"height" bson:"height"`           // 세로길이
	TimecodeIn  string  `json:"timecodein" bson:"timecodein"`   // 시작 타임코드
	TimecodeOut string  `json:"timecodeout" bson:"timecodeout"` // 마지막 타임코드
	Length      int     `json:"length" bson:"length"`           // 소스 길이
	InputCodec  string  `json:"inputcodec" bson:"inputcodec"`   // 소스 코덱
	OutputCodec string  `json:"outputcodec" bson:"outputcodec"` // 설정하는 아웃풋 코덱. 웹이라서 일반적으로 H.264를 사용한다.
	Fps         float64 `json:"fps" bson:"fps"`                 // FPS
	Rollmedia   string  `json:"rollmedia" bson:"rollmedia"`     // 촬영한 데이터라면, 카메라에서 생성된 데이터 이름

	// 자동설정
	Error string // 에러기록
}
