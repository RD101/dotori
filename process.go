package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/sys/unix"
)

func processingItem() error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMonogDBURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	// Status가 Ready인 item을 가져온다.
	item, err := GetReadyItem(client)
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("GetReadyItem 완료")
	}
	// thumbnail 폴더를 생성한다.
	err = genThumbDir(item)
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbDir 완료")
	}
	// 썸네일 이미지를 생성한다.
	err = genThumbImage(item)
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbImage 완료")
	}
	// 썸네일 동영상을 생성한다.
	err = getThumbContainers(item)
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbContainers 완료")
	}
	return nil
}

//genThumbDir 은 인수로 받은 아이템의 경로에 thumbnail 폴더를 생성한다.
func genThumbDir(item Item) error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMonogDBURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	// 연산전에 Admin셋팅을 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		return err
	}
	umask, err := strconv.Atoi(adminSetting.Umask)
	if err != nil {
		return err
	}
	unix.Umask(umask)
	per, err := strconv.ParseInt(adminSetting.FolderPermission, 8, 64)
	if err != nil {
		return err
	}
	// 생성할 경로를 가져온다.
	path := path.Dir(item.OutputThumbnailPngPath)
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	return nil
}

// genThumbImage 함수는 인수로 받은 아이템의 썸네일 이미지를 만든다.
func genThumbImage(item Item) error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMonogDBURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	collection := client.Database(*flagDBName).Collection(item.ItemType)
	// Status를 CreatingThumbnail로 바꾼다.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": bson.M{"status": CreatingThumbnail}})
	if err != nil {
		return err
	}
	// 연산전에 Admin셋팅을 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		return err
	}
	umask, err := strconv.Atoi(adminSetting.Umask)
	if err != nil {
		return err
	}
	unix.Umask(umask)
	// 변환할 이미지를 가져온다.
	path := item.InputThumbnailImgPath
	target, err := imaging.Open(path)
	if err != nil {
		return err
	}
	// Resize the cropped image to width = 200px preserving the aspect ratio.
	result := imaging.Fill(target, 320, 180, imaging.Center, imaging.Lanczos)
	// 저장할 경로를 생성
	per, err := strconv.ParseInt(adminSetting.FolderPermission, 8, 64)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(item.OutputThumbnailPngPath), os.FileMode(per))
	if err != nil {
		return err
	}
	//생성한 경로에 연산된 이미지 저장
	err = imaging.Save(result, item.OutputThumbnailPngPath)
	if err != nil {
		return err
	}
	// Status를 CreatedThumbnail로 바꾼다.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": bson.M{"status": CreatedThumbnail}})
	if err != nil {
		return err
	}
	return nil
}

// getThumbContainers 함수는 인수로 받은 아이템의 동영상을 만든다.
func getThumbContainers(item Item) error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMonogDBURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	collection := client.Database(*flagDBName).Collection(item.ItemType)
	// Status를 StartContainers로 바꾼다.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": bson.M{"status": StartContainers}})
	if err != nil {
		return err
	}
	// 연산전에 Admin셋팅을 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		return err
	}
	umask, err := strconv.Atoi(adminSetting.Umask)
	if err != nil {
		return err
	}
	unix.Umask(umask)
	// ogg 생성
	_, err = collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": bson.M{"status": CreatingOggContainer}})
	if err != nil {
		return err
	}
	per, err := strconv.ParseInt(adminSetting.FolderPermission, 8, 64)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(item.OutputThumbnailOggPath), os.FileMode(per))
	if err != nil {
		return err
	}
	argsOgg := []string{
		"-i",
		item.InputThumbnailClipPath,
		"-c:v",
		adminSetting.VideoCodecOgg,
		"-qscale:v",
		"7",
		"-s",
		fmt.Sprintf("%dx%d", adminSetting.ThumbnailContainerWidth, adminSetting.ThumbnailContainerHeight),
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		argsOgg = append(argsOgg, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		argsOgg = append(argsOgg, "-c:a")
		argsOgg = append(argsOgg, adminSetting.AudioCodec)
	}
	argsOgg = append(argsOgg, item.OutputThumbnailOggPath)
	if *flagDebug {
		fmt.Println(adminSetting.FFmpeg, strings.Join(argsOgg, " "))
	}
	err = exec.Command(adminSetting.FFmpeg, argsOgg...).Run()
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": bson.M{"status": CreatedOggContainer}})
	if err != nil {
		return err
	}
	// mp4 생성 - 작성예정
	// mov 생성 - 작성예정
	// 종료상태로 변경
	_, err = collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": bson.M{"status": CreatedContainers}})
	if err != nil {
		return err
	}
	return nil
}
