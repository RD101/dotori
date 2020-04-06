package main

import (
	"context"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	var ds DiskStatus

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return ds, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		return ds, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return ds, err
	}

	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		return ds, err
	}

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
