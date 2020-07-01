package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	if *flagMaxProcessNum < getProcessNum {
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
		err = ProcessFootageItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	case "nuke": // 뉴크파일
		err = ProcessNukeItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
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
		err = ProcessAlembicItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	case "houdini": // 후디니
		err = ProcessHoudiniItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	case "openvdb": // 볼륨데이터
		return
	case "pdf": // 문서
		return
	case "ies": // 조명파일
		return
	case "hdri": // HDRI 이미지, 환경맵
		return
	case "blender": // 블렌더 파일
		err = ProcessBlenderItem(client, adminSetting, item)
		if err != nil {
			log.Println(err)
			err = SetLog(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
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
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
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
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdthumbimg")
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggmedia")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdoggmedia")
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdmovmedia")
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdmp4media")
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessHoudiniItem 함수는 houdini 아이템을 연산한다.
func ProcessHoudiniItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creatingthumbdir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creatingthumbimg")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggmedia")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessFootageItem 함수는 footage 아이템을 연산한다.
func ProcessFootageItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail dir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}

	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbFootage(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// 썸네일을 생성하였다. 썸네일이 업로드 되었다고 체크한다.
	err = SetThumbImgUploaded(client, item, true)
	if err != nil {
		return err
	}

	err = SetStatus(client, item, "creating proxy dir")
	if err != nil {
		return err
	}

	err = genProxyDir(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}

	// Proxy 이미지를 생성한다.
	err = SetStatus(client, item, "creating proxy sequence")
	if err != nil {
		return err
	}

	err = genProxySequence(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}

	// 썸네일 동영상 생성

	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genProxyToOggMedia(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}

	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genProxyToMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}

	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genProxyToMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}

	// Proxy 제거
	err = SetStatus(client, item, "removing proxy")
	if err != nil {
		return err
	}

	err = os.RemoveAll(item.OutputProxyImgPath)
	if err != nil {
		return err
	}

	// 완료
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessNukeItem 함수는 nuke 아이템을 연산한다.
func ProcessNukeItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creatingthumbdir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creatingthumbimg")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggmedia")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
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
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdthumbimg")
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggmedia")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdoggmedia")
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdmovmedia")
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdmp4media")
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessAlembicItem 함수는 alembic 아이템을 연산한다.
func ProcessAlembicItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creatingthumbdir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
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
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggmedia")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
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
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessBlenderItem 함수는 blender 아이템을 연산한다.
func ProcessBlenderItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
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
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdthumbimg")
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingoggmedia")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdoggmedia")
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdmovmedia")
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmp4media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		// 상태를 error로 바꾼다.
		err = SetStatus(client, item, "error")
		if err != nil {
			return err
		}
		// 에러 내용을 로그로 남긴다.
		err = SetLog(client, item.ID.Hex(), err.Error())
		if err != nil {
			return err
		}
		return err
	}
	err = SetStatus(client, item, "createdmp4media")
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

//genProxyDir 은 인수로 받은 아이템의 경로에 thumbnail 폴더를 생성한다.
func genProxyDir(adminSetting Adminsetting, item Item) error {
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
	path := path.Dir(item.OutputProxyImgPath)
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

// genThumbFootage 함수는 Footage 아이템정보를 이용하여 oiiotool 명령어로 썸네일 이미지를 만든다.
func genThumbFootage(adminSetting Adminsetting, item Item) error {
	// 변환할 이미지를 한프레임 가지고온다.
	path := item.OutputDataPath
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	input := ""
	// 연산할 파일을 하나 구한다.
	switch len(files) {
	case 0: // 파일이 없을 때
		return errors.New("no files")
	case 1: // 파일이 한개일 때
		input = files[0].Name()
	default: // 중간프레임 파일명을 구한다.
		input = files[len(files)/2].Name()
	}
	args := []string{
		path + "/" + input,
		"--colorconvert",
		item.InColorspace,
		item.OutColorspace,
		"--fit",
		fmt.Sprintf("%dx%d", adminSetting.ThumbnailImageWidth, adminSetting.ThumbnailImageHeight),
		"-o",
		item.OutputThumbnailPngPath,
	}
	// OIIO 이미지 연산을 위해 OCIO 환경변수를 설정한다.
	_, err = os.Stat(adminSetting.OCIOConfig)
	if os.IsNotExist(err) {
		return errors.New("admin 셋팅에서 OCIOConfig 값을 설정해주세요")
	}
	os.Setenv("OCIO", adminSetting.OCIOConfig)
	err = exec.Command(adminSetting.OpenImageIO, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// genProxySequence 함수는 Footage 아이템정보를 이용하여 oiiotool 명령어로 Proxy 이미지를 만든다.
func genProxySequence(adminSetting Adminsetting, item Item) error {
	// 변환할 이미지를 한프레임 가지고온다.
	path := item.OutputDataPath
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	// OIIO 이미지 연산을 위해 OCIO 환경변수를 설정한다.
	_, err = os.Stat(adminSetting.OCIOConfig)
	if os.IsNotExist(err) {
		return errors.New("admin 셋팅에서 OCIOConfig 값을 설정해주세요")
	}
	os.Setenv("OCIO", adminSetting.OCIOConfig)

	// 각 파일을 돌면서 연산을 진행한다.
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		rmExt := strings.TrimSuffix(file.Name(), ext)
		args := []string{
			path + "/" + file.Name(),
			"--colorconvert",
			item.InColorspace,
			item.OutColorspace,
			"-o",
			item.OutputProxyImgPath + rmExt + ".png",
		}
		err = exec.Command(adminSetting.OpenImageIO, args...).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// genThumbOggMedia 함수는 인수로 받은 아이템의 .ogg 동영상을 만든다.
func genThumbOggMedia(adminSetting Adminsetting, item Item) error {
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
			adminSetting.MediaWidth,
			adminSetting.MediaHeight,
			adminSetting.MediaWidth,
			adminSetting.MediaHeight,
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

// genThumbMovMedia 함수는 인수로 받은 아이템의 .mov 동영상을 만든다.
func genThumbMovMedia(adminSetting Adminsetting, item Item) error {
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
			adminSetting.MediaWidth,
			adminSetting.MediaHeight,
			adminSetting.MediaWidth,
			adminSetting.MediaHeight,
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

// genThumbMp4Media 함수는 인수로 받은 아이템의 .mp4 동영상을 만든다.
func genThumbMp4Media(adminSetting Adminsetting, item Item) error {
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
			adminSetting.MediaWidth,
			adminSetting.MediaHeight,
			adminSetting.MediaWidth,
			adminSetting.MediaHeight,
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

// genProxyToOggMedia 함수는 인수로 받은 아이템의 .ogg 동영상을 만든다.
func genProxyToOggMedia(adminSetting Adminsetting, item Item) error {
	seqs, err := searchSeq(item.OutputProxyImgPath)
	if err != nil {
		return err
	}
	for _, seq := range seqs {
		args := []string{
			"-f",
			"image2",
			"-start_number",
			strconv.Itoa(seq.FrameIn),
			"-r",
			"24",
			"-y",
			"-i",
			fmt.Sprintf("%s/%s", seq.Dir, seq.Base),
			"-pix_fmt",
			"yuv420p",
			"-c:v",
			adminSetting.VideoCodecOgg,
			"-qscale:v",
			"7",
			"-vf",
			fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1",
				adminSetting.MediaWidth,
				adminSetting.MediaHeight,
				adminSetting.MediaWidth,
				adminSetting.MediaHeight,
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
	}
	return nil
}

// genProxyToMovMedia 함수는 인수로 받은 아이템의 .mov 동영상을 만든다.
func genProxyToMovMedia(adminSetting Adminsetting, item Item) error {
	seqs, err := searchSeq(item.OutputProxyImgPath)
	if err != nil {
		return err
	}
	for _, seq := range seqs {
		args := []string{
			"-f",
			"image2",
			"-start_number",
			strconv.Itoa(seq.FrameIn),
			"-r",
			"24",
			"-y",
			"-i",
			fmt.Sprintf("%s/%s", seq.Dir, seq.Base),
			"-pix_fmt",
			"yuv420p",
			"-c:v",
			adminSetting.VideoCodecMov,
			"-qscale:v",
			"7",
			"-vf",
			fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1",
				adminSetting.MediaWidth,
				adminSetting.MediaHeight,
				adminSetting.MediaWidth,
				adminSetting.MediaHeight,
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
	}
	return nil
}

// genProxyToMp4Media 함수는 인수로 받은 아이템의 .mp4 동영상을 만든다.
func genProxyToMp4Media(adminSetting Adminsetting, item Item) error {
	seqs, err := searchSeq(item.OutputProxyImgPath)
	if err != nil {
		return err
	}
	for _, seq := range seqs {
		args := []string{
			"-f",
			"image2",
			"-start_number",
			strconv.Itoa(seq.FrameIn),
			"-r",
			"24",
			"-y",
			"-i",
			fmt.Sprintf("%s/%s", seq.Dir, seq.Base),
			"-pix_fmt",
			"yuv420p",
			"-c:v",
			adminSetting.VideoCodecMp4,
			"-qscale:v",
			"7",
			"-vf",
			fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1",
				adminSetting.MediaWidth,
				adminSetting.MediaHeight,
				adminSetting.MediaWidth,
				adminSetting.MediaHeight,
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
	}
	return nil
}
