package main

import (
	"fmt"
	"net/http"
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

// DiskUsage함수는 인자로 넣은 path의 디스크 용량을 확인 하는 함수이다. (syscall 패키지는 linux/unix에만 가능, window 빌드 불가능)
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

// DiskCheck함수는 rootPath의 디스크용량을 확인하는 함수이다.
func DiskCheck(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rootpath, err := GetRootPath(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(rootpath) == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	disk := DiskUsage(rootpath)
	return disk
}