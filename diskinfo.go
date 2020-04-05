package main

import (
	"log"
	"syscall"

	"gopkg.in/mgo.v2"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// DiskCheck함수는 rootPath의 디스크용량을 확인하는 함수이다.
func DiskCheck() (DiskStatus, error) {

	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(session)
	if err != nil {
		log.Fatal(err)
	}

	var ds DiskStatus

	// rootpath경로의 디스크 용량을 확인한다.
	fs := syscall.Statfs_t{}
	err = syscall.Statfs(rootpath, &fs)
	if err != nil {
		return ds, err
	}

	ds.All = fs.Blocks * uint64(fs.Bsize)
	ds.Free = fs.Bfree * uint64(fs.Bsize)
	ds.Used = ds.All - ds.Free
	return ds, nil
}
