package main

type Storage struct {
	ID   string // 스토리지 ID
	Path string // 스토리지 물리적 경로. 바뀔 수 있어야 한다.
}

type Item struct {
	ID         string   // ID
	Tags       []string // 태그리스트
	Thumbimg   string   // 썸네일 이미지 주소
	Thumbmov   string   // 썸네일 영상 주소
	Inputpath  string   // 최초 등록되는 경로
	Outputpath string   // 저장되는 경로
	Type       string   // maya, source, houdini, blender, nuke ..  같은 형태인가.
	Status     string   // 상태(에러, done, wip)
	Updatetime string   // UTC 타임으로 들어가도록 하기.
	Isrm       bool     // 삭제 판단 값
}
