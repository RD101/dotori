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

// Processing 함수는 일정시간마다 프로세스를 실행시킨다.
func Processing() {
	for {
		time.Sleep(time.Duration(*flagProcessInterval) * 1000 * time.Millisecond)
		//go ProcessDemo()
		go processingItem()
	}
}

// ProcessDemo 함수는 go 프로세스가 잘 실행되는지 테스트하는 함수이다.
func ProcessDemo() {
	fmt.Println("processing", *flagProcessNum)
	fmt.Println("wait", *flagProcessInterval, "sec")
}

func processingItem() {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	// AdminSetting을 DB에서 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	// 연산해야할 리스트(FileUploaded) 갯수가 몇개인지 구한다.
	fileUploadedItemsNum, err := GetFileUploadedItemsNum(client)
	if err != nil {
		return
	}
	// 연산할 리스트가 없다면 return 한다.
	if fileUploadedItemsNum == 0 {
		return
	}
	// Status가 FileUploaded인 item을 가져온다.
	item, err := GetFileUploadedItem(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetLog(client, item.ItemType, item.ID.Hex(), "GetFileUploadedItem 완료")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	// thumbnail 폴더를 생성한다.
	err = SetStatus(client, item, "creatingthumbdir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetLog(client, item.ItemType, item.ID.Hex(), "genThumbDir 완료")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetStatus(client, item, "createdthumbdir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creatingthumbimg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetLog(client, item.ItemType, item.ID.Hex(), "genThumbImage 완료")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetStatus(client, item, "createdthumbimg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggcontainer")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = genThumbOggContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetLog(client, item.ItemType, item.ID.Hex(), "genThumbOggContainer 완료")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetStatus(client, item, "createdoggcontainer")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovcontainer")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = genThumbMovContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetLog(client, item.ItemType, item.ID.Hex(), "genThumbMovContainer 완료")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetStatus(client, item, "createdmovcontainer")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4container")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = genThumbMp4Container(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetLog(client, item.ItemType, item.ID.Hex(), "genThumbMovContainer 완료")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetLog(client, item.ItemType, item.ID.Hex(), "genThumbMp4Container 완료")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetStatus(client, item, "createdmp4container")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetStatus(client, item, "createdcontainers")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		fmt.Fprintf(os.Stderr, "process: %s", err)
		return
	}
	return
}

//genThumbDir 은 인수로 받은 아이템의 경로에 thumbnail 폴더를 생성한다.
func genThumbDir(adminSetting Adminsetting, item Item) error {
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
	// 변환할 이미지를 가져온다.
	path := item.InputThumbnailImgPath
	target, err := imaging.Open(path)
	if err != nil {
		return err
	}
	// Resize the cropped image to width = 200px preserving the aspect ratio.
	result := imaging.Fill(target, adminSetting.ThumbnailImageWidth, adminSetting.ThumbnailImageHeight, imaging.Center, imaging.Lanczos)
	//생성한 경로에 연산된 이미지 저장
	err = imaging.Save(result, item.OutputThumbnailPngPath)
	if err != nil {
		return err
	}
	return nil
}

// genThumbOggContainer 함수는 인수로 받은 아이템의 .ogg 동영상을 만든다.
func genThumbOggContainer(adminSetting Adminsetting, item Item) error {
	args := []string{
		"-y",
		"-i",
		item.InputThumbnailClipPath,
		"-c:v",
		adminSetting.VideoCodecOgg,
		"-qscale:v",
		"7",
		"-vf",
		fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1",
			adminSetting.ThumbnailContainerWidth,
			adminSetting.ThumbnailContainerHeight,
			adminSetting.ThumbnailContainerWidth,
			adminSetting.ThumbnailContainerHeight,
		),
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

// genThumbMovContainer 함수는 인수로 받은 아이템의 .mov 동영상을 만든다.
func genThumbMovContainer(adminSetting Adminsetting, item Item) error {
	args := []string{
		"-y",
		"-i",
		item.InputThumbnailClipPath,
		"-c:v",
		adminSetting.VideoCodecMov,
		"-qscale:v",
		"7",
		"-vf",
		fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1",
			adminSetting.ThumbnailContainerWidth,
			adminSetting.ThumbnailContainerHeight,
			adminSetting.ThumbnailContainerWidth,
			adminSetting.ThumbnailContainerHeight,
		),
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

// genThumbMp4Container 함수는 인수로 받은 아이템의 .mp4 동영상을 만든다.
func genThumbMp4Container(adminSetting Adminsetting, item Item) error {
	args := []string{
		"-y",
		"-i",
		item.InputThumbnailClipPath,
		"-c:v",
		adminSetting.VideoCodecMp4,
		"-qscale:v",
		"7",
		"-vf",
		fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1",
			adminSetting.ThumbnailContainerWidth,
			adminSetting.ThumbnailContainerHeight,
			adminSetting.ThumbnailContainerWidth,
			adminSetting.ThumbnailContainerHeight,
		),
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
