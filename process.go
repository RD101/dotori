package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/sys/unix"
)

func processingItem() error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
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
	// AdminSetting을 DB에서 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
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
	err = SetStatus(client, item, CreatingThumbDir)
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbDir 완료")
	}
	err = SetStatus(client, item, CreatedThumbDir)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, CreatingThumbImg)
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbImage 완료")
	}
	err = SetStatus(client, item, CreatedThumbImg)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, CreatingOggContainer)
	if err != nil {
		return err
	}
	err = getThumbOggContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbOggContainer 완료")
	}
	err = SetStatus(client, item, CreatedOggContainer)
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, CreatingMovContainer)
	if err != nil {
		return err
	}
	err = getThumbMovContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbMovContainer 완료")
	}
	err = SetStatus(client, item, CreatedMovContainer)
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, CreatingMp4Container)
	if err != nil {
		return err
	}
	err = getThumbMp4Container(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	if *flagDebug {
		fmt.Println("genThumbMp4Container 완료")
	}
	err = SetStatus(client, item, CreatedMp4Container)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, CreatedContainers)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, Done)
	if err != nil {
		return err
	}
	return nil
}

//genThumbDir 은 인수로 받은 아이템의 경로에 thumbnail 폴더를 생성한다.
func genThumbDir(adminSetting Adminsetting, item Item) error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
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
	// umask, 권한 셋팅
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
func genThumbImage(adminSetting Adminsetting, item Item) error {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
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
	// 변환할 이미지를 가져온다.
	path := item.InputThumbnailImgPath
	target, err := imaging.Open(path)
	if err != nil {
		return err
	}
	// Resize the cropped image to width = 200px preserving the aspect ratio.
	result := imaging.Fill(target, 320, 180, imaging.Center, imaging.Lanczos)
	//생성한 경로에 연산된 이미지 저장
	err = imaging.Save(result, item.OutputThumbnailPngPath)
	if err != nil {
		return err
	}
	return nil
}

// getThumbOggContainer 함수는 인수로 받은 아이템의 .ogg 동영상을 만든다.
func getThumbOggContainer(adminSetting Adminsetting, item Item) error {
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
	err := exec.Command(adminSetting.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// getThumbMovContainer 함수는 인수로 받은 아이템의 .mov 동영상을 만든다.
func getThumbMovContainer(adminSetting Adminsetting, item Item) error {
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
	err := exec.Command(adminSetting.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// getThumbMp4Container 함수는 인수로 받은 아이템의 .mp4 동영상을 만든다.
func getThumbMp4Container(adminSetting Adminsetting, item Item) error {
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
	err := exec.Command(adminSetting.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}
