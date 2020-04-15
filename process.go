package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
	// 썸네일 이미지를 생성한다.
	err = genThumbImage(item)
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbImage 완료")
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = getThumbOggContainer(item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = getThumbMovContainer(item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = getThumbMp4Container(item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
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

// getThumbOggContainer 함수는 인수로 받은 아이템의 .ogg 동영상을 만든다.
func getThumbOggContainer(item Item) error {
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
	args := []string{
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
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	args = append(args, item.OutputThumbnailOggPath)
	if *flagDebug {
		fmt.Println(adminSetting.FFmpeg, strings.Join(args, " "))
	}
	err = exec.Command(adminSetting.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// getThumbMovContainer 함수는 인수로 받은 아이템의 .mov 동영상을 만든다.
func getThumbMovContainer(item Item) error {
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
	args := []string{
		"-i",
		item.InputThumbnailClipPath,
		"-c:v",
		adminSetting.VideoCodecMov,
		"-qscale:v",
		"7",
		"-s",
		fmt.Sprintf("%dx%d", adminSetting.ThumbnailContainerWidth, adminSetting.ThumbnailContainerHeight),
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	args = append(args, item.OutputThumbnailMovPath)
	if *flagDebug {
		fmt.Println(adminSetting.FFmpeg, strings.Join(args, " "))
	}
	err = exec.Command(adminSetting.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// getThumbMp4Container 함수는 인수로 받은 아이템의 .mp4 동영상을 만든다.
func getThumbMp4Container(item Item) error {
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
	args := []string{
		"-i",
		item.InputThumbnailClipPath,
		"-c:v",
		adminSetting.VideoCodecMp4,
		"-qscale:v",
		"7",
		"-s",
		fmt.Sprintf("%dx%d", adminSetting.ThumbnailContainerWidth, adminSetting.ThumbnailContainerHeight),
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	args = append(args, item.OutputThumbnailMp4Path)
	if *flagDebug {
		fmt.Println(adminSetting.FFmpeg, strings.Join(args, " "))
	}
	err = exec.Command(adminSetting.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}
