package main

import (
	"context"
	"fmt"
	"log"
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
		go processingItem()
	}
}

func processingItem() {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
		return
	}
	// 연산 갯수를 체크한다.
	getProcessNum, err := GetProcessingItemNum(client)
	if err != nil {
		log.Println(err)
		return
	}
	if *flagProcessNum < getProcessNum {
		return
	}
	// AdminSetting을 DB에서 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		log.Println(err)
		return
	}
	// Status가 FileUploaded인 item을 가져온다.
	item, err := GetFileUploadedItem(client)
	if err != nil {
		// 가지고 올 문서가 없다면 그냥 return 한다.
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Println(err)
		return
	}
	// ItemType별로 연산한다.
	switch item.ItemType {
	case "maya":
		err = ProcessMayaItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	case "footage": // Footage 소스, 시퀀스
		return
	case "nuke": // 뉴크파일
		return
	case "usd": // Pixar USD
		err = ProcessUSDItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	case "alembic": // Alembic
		return
	case "houdini": // 후디니
		return
	case "openvdb": // 볼륨데이터
		return
	case "pdf": // 문서
		return
	case "ies": // 조명파일
		return
	case "hdri": // HDRI 이미지, 환경맵
		return
	case "blender": // 블렌더 파일
		return
	case "texture": // 텍스쳐
		return
	case "psd": // 포토샵 파일
		return
	case "modo": // 모도
		return
	case "lut", "3dl", "blut", "cms", "csp", "cub", "cube", "vf", "vfz": // LUT 파일들
		return
	case "sound":
		err = ProcessSoundItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	default:
		log.Println("약속된 type이 아닙니다")
		return
	}
}

// ProcessMayaItem 함수는 maya 아이템을 연산한다.
func ProcessMayaItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creatingthumbdir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdthumbdir")
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creatingthumbimg")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdthumbimg")
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggcontainer")
	if err != nil {
		return err
	}
	err = genThumbOggContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdoggcontainer")
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovcontainer")
	if err != nil {
		return err
	}
	err = genThumbMovContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdmovcontainer")
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4container")
	if err != nil {
		return err
	}
	err = genThumbMp4Container(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdmp4container")
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessUSDItem 함수는 USD 아이템을 연산한다.
func ProcessUSDItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creatingthumbdir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdthumbdir")
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creatingthumbimg")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdthumbimg")
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggcontainer")
	if err != nil {
		return err
	}
	err = genThumbOggContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdoggcontainer")
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovcontainer")
	if err != nil {
		return err
	}
	err = genThumbMovContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdmovcontainer")
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4container")
	if err != nil {
		return err
	}
	err = genThumbMp4Container(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdmp4container")
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessSoundItem 함수는 sound 아이템을 연산한다.
func ProcessSoundItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creatingthumbdir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdthumbdir")
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creatingthumbimg")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdthumbimg")
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggcontainer")
	if err != nil {
		return err
	}
	err = genThumbOggContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdoggcontainer")
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovcontainer")
	if err != nil {
		return err
	}
	err = genThumbMovContainer(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdmovcontainer")
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4container")
	if err != nil {
		return err
	}
	err = genThumbMp4Container(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "createdmp4container")
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
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
