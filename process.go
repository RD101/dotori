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

// ProcessMain 함수는 프로세스 전체 흐름을 만드는 함수
func ProcessMain() {
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
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
		return
	}
	// AdminSetting을 DB에서 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		log.Println(err)
		return
	}
	client.Disconnect(ctx)
	cancel()

	// 버퍼 채널을 만든다.
	jobs := make(chan Item, adminSetting.ProcessBufferSize)

	// worker 프로세스를 지정한 개수만큼 실행시킨다.
	for w := 1; w <= *flagMaxProcessNum; w++ {
		go worker(jobs)
	}

	// queueingItem을 실행시킨다.
	go queueingItem(jobs)

	select {}
}

// 실제 연산을 하는 worker
func worker(jobs <-chan Item) {
	for j := range jobs {
		processingItem(j)
	}
}

func processingItem(item Item) {
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
	// AdminSetting을 DB에서 가지고 온다.
	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		log.Println(err)
		return
	}
	// LD_LIBRARY_PATH 환경변수를 세팅한다.
	adminLdLibPath := adminSetting.LDLibraryPath
	if adminLdLibPath != "" {
		os.Setenv("LD_LIBRARY_PATH", adminLdLibPath)
	}
	// ItemType별로 연산한다.
	switch item.ItemType {
	case "maya":
		err = ProcessMayaItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "max":
		err = ProcessMaxItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "fusion360":
		err = ProcessFusion360Item(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "clip": // Clip 소스
		err = ProcessClipItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "footage": // Footage 소스, 시퀀스
		err = ProcessFootageItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "nuke": // 뉴크파일
		err = ProcessNukeItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "usd": // Pixar USD
		err = ProcessUSDItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "alembic": // Alembic
		err = ProcessAlembicItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "houdini": // 후디니
		err = ProcessHoudiniItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "openvdb": // 볼륨데이터
		err = ProcessOpenVDBItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "pdf": // 문서
		err = ProcessPdfItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "ies": // 조명파일
		err = ProcessIesItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "hdri": // HDRI 이미지, 환경맵
		err = ProcessHDRIItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "texture": // Texture 파일
		err = ProcessTextureItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "blender": // 블렌더 파일
		err = ProcessBlenderItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "modo": // 모도
		err = ProcessModoItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "katana": //katana
		err = ProcessKatanaItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "lut": // LUT 파일들
		err = ProcessLutItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "sound":
		err = ProcessSoundItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "unreal":
		err = ProcessUnrealItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "hwp":
		err = ProcessHwpItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "ppt":
		err = ProcessPptItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	case "matte":
		err = ProcessMatteItem(client, adminSetting, item)
		if err != nil {
			err = SetErrStatus(client, item.ID.Hex(), err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		return
	default:
		log.Println("약속된 type이 아닙니다")
		return
	}
}

// 아이템을 가져와서 버퍼 채널에 채우는 함수
func queueingItem(jobs chan<- Item) {
	for {
		item, err := GetFileUploadedItem()
		if err != nil {
			// 가지고 올 문서가 없다면 10초 기다렸다가 continue
			if err == mongo.ErrNoDocuments {
				time.Sleep(time.Second * 10)
				continue
			}
			// DB에서 아이템을 가지고 오는 과정에서 에러가 발생하면 로그 출력후 10초를 기다리고 다시 진행
			log.Println(err)
			time.Sleep(time.Second * 10)
			continue
		}
		updatedItem, err := SetStatusAndGetItem(item, "queued for processing")
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 10)
			continue
		}
		jobs <- updatedItem
		// 10초후 다시 queueing 한다.
		time.Sleep(time.Second * 10)
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
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
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
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessMaxItem 함수는 max 아이템을 연산한다.
func ProcessMaxItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbn dir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumb img")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessFusion360Item 함수는 maya 아이템을 연산한다.
func ProcessFusion360Item(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumb dir")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumb img")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessClipItem 함수는 clip 아이템을 연산한다.
func ProcessClipItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// Item의 Data가 복사될 경로를 생성한다.
	if item.RequireMkdirInProcess {
		err := SetStatus(client, item, "creating item data directory")
		if err != nil {
			return err
		}
		err = genOutputDataPath(adminSetting, item)
		if err != nil {
			return err
		}
	}

	// InputData의 파일을 복사한다.
	if item.RequireCopyInProcess {
		err := SetStatus(client, item, "copy input data")
		if err != nil {
			return err
		}
		err = copyInputDataToOutputDataPathClip(adminSetting, item)
		if err != nil {
			return err
		}
	}

	// Thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating ogg media")
	if err != nil {
		return err
	}
	err = genClipToOggMedia(adminSetting, item)
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating mov media")
	if err != nil {
		return err
	}
	err = genClipToMovMedia(adminSetting, item)
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating mp4 media")
	if err != nil {
		return err
	}
	err = genClipToMp4Media(adminSetting, item)
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
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
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
	// Item의 Data가 복사될 경로를 생성한다.
	if item.RequireMkdirInProcess {
		err := SetStatus(client, item, "creating item data directory")
		if err != nil {
			return err
		}
		err = genOutputDataPath(adminSetting, item)
		if err != nil {
			return err
		}
	}

	// InputData의 파일을 복사한다.
	if item.RequireCopyInProcess {
		err := SetStatus(client, item, "copy input data")
		if err != nil {
			return err
		}
		err = copyInputDataToOutputDataPathFootage(adminSetting, item)
		if err != nil {
			return err
		}
	}

	// Thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}

	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbFootage(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetThumbImgUploaded(client, item, true) // 썸네일이 생성되었으니 썸네일이 업로드 되었다고 체크한다.
	if err != nil {
		return err
	}

	// 프록시 폴더를 생성
	err = SetStatus(client, item, "creating proxy dir")
	if err != nil {
		return err
	}
	err = genProxyDir(adminSetting, item)
	if err != nil {
		return err
	}

	// Proxy 이미지를 생성한다.
	err = SetStatus(client, item, "creating proxy sequence")
	if err != nil {
		return err
	}
	err = genProxySequence(adminSetting, item)
	if err != nil {
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
		return err
	}

	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genProxyToMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}

	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genProxyToMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetThumbClipUploaded(client, item, true) // 모든 썸네일 동영상을 생성하였다. 썸네일 동영상이 업로드 되었다고 체크한다.
	if err != nil {
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

// ProcessHDRIItem 함수는 HDRI 아이템을 연산한다.
func ProcessHDRIItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// Thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}

	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbHDRI(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetThumbImgUploaded(client, item, true) // 썸네일이 생성되었으니 썸네일이 업로드 되었다고 체크한다.
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

// ProcessTextureItem 함수는 Texture 아이템을 연산한다.
func ProcessTextureItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// Thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}

	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbTexture(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetThumbImgUploaded(client, item, true) // 썸네일이 생성되었으니 썸네일이 업로드 되었다고 체크한다.
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
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
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
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
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
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessOpenVDBItem 함수는 OpenVDB 아이템을 연산한다.
func ProcessOpenVDBItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
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
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creatingmovmedia")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
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
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
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
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessLutItem 함수는 LUT 아이템을 연산한다.
func ProcessLutItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
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
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessMatteItem 함수는 Matte 아이템을 연산한다.
func ProcessMatteItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// Thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}

	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbTexture(adminSetting, item)
	if err != nil {
		return err
	}
	err = SetThumbImgUploaded(client, item, true) // 썸네일이 생성되었으니 썸네일이 업로드 되었다고 체크한다.
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

// copyInputDataToOuputDataPathFootage 함수는 인풋데이터를 아웃풋 데이터 경로에 복사한다.
func copyInputDataToOutputDataPathFootage(adminSetting Adminsetting, item Item) error {
	// 복사할 파일의 권한을 불러온다.
	fileP := adminSetting.FilePermission
	filePerm, err := strconv.ParseInt(fileP, 8, 64)
	if err != nil {
		return err
	}

	for i := item.InputData.FrameIn; i <= item.InputData.FrameOut; i++ {
		src := fmt.Sprintf(item.InputData.Dir+"/"+item.InputData.Base, i)
		dest := fmt.Sprintf(item.OutputDataPath+"/"+item.InputData.Base, i)
		// file을 복사한다.
		bytesRead, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(dest, bytesRead, os.FileMode(filePerm))
		if err != nil {
			return err
		}
		// 복사된 파일의 권한을 변경한다.
		err = setFilePermission(dest, adminSetting)
		if err != nil {
			return err
		}
	}
	return nil
}

func setFilePermission(path string, adminSetting Adminsetting) error {
	umask, err := strconv.Atoi(adminSetting.Umask)
	if err != nil {
		return err
	}
	unix.Umask(umask)
	uid, err := strconv.Atoi(adminSetting.UID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(adminSetting.GID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

// copyInputDataToOuputDataPathClip 함수는 인풋데이터를 아웃풋 데이터 경로에 복사한다.
func copyInputDataToOutputDataPathClip(adminSetting Adminsetting, item Item) error {
	// 복사할 파일의 권한을 불러온다.
	fileP := adminSetting.FilePermission
	filePerm, err := strconv.ParseInt(fileP, 8, 64)
	if err != nil {
		return err
	}
	src := fmt.Sprintf(item.InputData.Dir + "/" + item.InputData.Base)
	dest := fmt.Sprintf(item.OutputDataPath + "/" + item.InputData.Base)
	// file을 복사한다.
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dest, bytesRead, os.FileMode(filePerm))
	if err != nil {
		return err
	}
	// 복사된 파일의 권한을 변경한다.
	err = setFilePermission(dest, adminSetting)
	if err != nil {
		return err
	}
	return nil
}

//genOutputDataPath Item의 data 경로를 생성한다.
func genOutputDataPath(adminSetting Adminsetting, item Item) error {
	// umask, 권한 셋팅
	umask, err := strconv.Atoi(adminSetting.Umask)
	if err != nil {
		return err
	}
	unix.Umask(umask)
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(adminSetting.FolderPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	path := path.Dir(item.OutputDataPath)
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(adminSetting.UID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(adminSetting.GID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
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
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(adminSetting.UID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(adminSetting.GID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
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
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(adminSetting.UID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(adminSetting.GID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
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
	if *flagDebug {
		fmt.Println(adminSetting.OpenImageIO, strings.Join(args, " "))
	}
	err = exec.Command(adminSetting.OpenImageIO, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// genThumbHDRI 함수는 HDRI 아이템정보를 이용하여 oiiotool 명령어로 썸네일 이미지를 만든다.
func genThumbHDRI(adminSetting Adminsetting, item Item) error {
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
	default: // 파일이 2개 이상인 경우 에러처리 한다.
		return errors.New("파일이 여러 개가 존재합니다")
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
	if *flagDebug {
		fmt.Println(adminSetting.OpenImageIO, strings.Join(args, " "))
	}
	err = exec.Command(adminSetting.OpenImageIO, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// genThumbTexture 함수는 Texture 아이템정보를 이용하여 oiiotool 명령어로 썸네일 이미지를 만든다.
func genThumbTexture(adminSetting Adminsetting, item Item) error {
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
		return errors.New("파일이 여러개가 존재합니다")
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
	if *flagDebug {
		fmt.Println(adminSetting.OpenImageIO, strings.Join(args, " "))
	}
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
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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
	seqs, err := searchSeqAndClip(item.OutputProxyImgPath)
	if err != nil {
		return err
	}
	if len(seqs) != 1 {
		return errors.New("해당 경로에 여러 소스가 존재합니다")
	}
	target := seqs[0]
	if target.Length != target.FrameOut-target.FrameIn+1 {
		return errors.New("중간에 빈 프레임이 존재합니다")
	}
	args := []string{
		"-f",
		"image2",
		"-start_number",
		strconv.Itoa(target.FrameIn),
		"-r",
		item.Fps,
		"-y",
		"-i",
		fmt.Sprintf("%s/%s", target.Dir, target.Base),
		"-pix_fmt",
		"yuv420p",
		"-c:v",
		adminSetting.VideoCodecOgg,
		"-qscale:v",
		"7",
	}
	if item.Premultiply {
		args = append(args, []string{"-vf", "premultiply=inplace=1"}...)
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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

// genClipToOggMedia 함수는 인수로 받은 아이템의 .ogg 동영상을 만든다.
func genClipToOggMedia(adminSetting Adminsetting, item Item) error {
	files, err := searchSeqAndClip(item.OutputDataPath)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("해당 경로에 여러 소스가 존재합니다")
	}
	target := files[0]
	args := []string{
		"-y",
		"-i",
		fmt.Sprintf("%s/%s", target.Dir, target.Base),
		"-pix_fmt",
		"yuv420p",
		"-c:v",
		adminSetting.VideoCodecOgg,
		"-qscale:v",
		"7",
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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

// genProxyToMovMedia 함수는 인수로 받은 아이템의 .mov 동영상을 만든다.
func genProxyToMovMedia(adminSetting Adminsetting, item Item) error {
	seqs, err := searchSeqAndClip(item.OutputProxyImgPath)
	if err != nil {
		return err
	}
	if len(seqs) != 1 {
		return errors.New("해당 경로에 여러 소스가 존재합니다")
	}
	target := seqs[0]
	if target.Length != target.FrameOut-target.FrameIn+1 {
		return errors.New("중간에 빈 프레임이 존재합니다")
	}

	args := []string{
		"-f",
		"image2",
		"-start_number",
		strconv.Itoa(target.FrameIn),
		"-r",
		item.Fps,
		"-y",
		"-i",
		fmt.Sprintf("%s/%s", target.Dir, target.Base),
		"-pix_fmt",
		"yuv420p",
		"-c:v",
		adminSetting.VideoCodecMov,
		"-qscale:v",
		"7",
	}
	if item.Premultiply {
		args = append(args, []string{"-vf", "premultiply=inplace=1"}...)
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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

// genClipToMovMedia 함수는 인수로 받은 아이템의 .mov 동영상을 만든다.
func genClipToMovMedia(adminSetting Adminsetting, item Item) error {
	files, err := searchSeqAndClip(item.OutputDataPath)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("해당 경로에 여러 소스가 존재합니다")
	}
	target := files[0]
	args := []string{
		"-y",
		"-i",
		fmt.Sprintf("%s/%s", target.Dir, target.Base),
		"-pix_fmt",
		"yuv420p",
		"-c:v",
		adminSetting.VideoCodecMov,
		"-qscale:v",
		"7",
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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

// genProxyToMp4Media 함수는 인수로 받은 아이템의 .mp4 동영상을 만든다.
func genProxyToMp4Media(adminSetting Adminsetting, item Item) error {
	seqs, err := searchSeqAndClip(item.OutputProxyImgPath)
	if err != nil {
		return err
	}
	if len(seqs) != 1 {
		return errors.New("해당 경로에 여러 소스가 존재합니다")
	}
	target := seqs[0]
	if target.Length != target.FrameOut-target.FrameIn+1 {
		return errors.New("중간에 빈 프레임이 존재합니다")
	}
	args := []string{
		"-f",
		"image2",
		"-start_number",
		strconv.Itoa(target.FrameIn),
		"-r",
		item.Fps,
		"-y",
		"-i",
		fmt.Sprintf("%s/%s", target.Dir, target.Base),
		"-pix_fmt",
		"yuv420p",
		"-c:v",
		adminSetting.VideoCodecMp4,
		"-qscale:v",
		"7",
	}
	if item.Premultiply {
		args = append(args, []string{"-vf", "premultiply=inplace=1"}...)
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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

// genClipToMp4Media 함수는 인수로 받은 아이템의 .mp4 동영상을 만든다.
func genClipToMp4Media(adminSetting Adminsetting, item Item) error {
	files, err := searchSeqAndClip(item.OutputDataPath)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("해당 경로에 여러 소스가 존재합니다")
	}
	target := files[0]
	args := []string{
		"-y",
		"-i",
		fmt.Sprintf("%s/%s", target.Dir, target.Base),
		"-pix_fmt",
		"yuv420p",
		"-c:v",
		adminSetting.VideoCodecMp4,
		"-qscale:v",
		"7",
	}
	if adminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, adminSetting.AudioCodec)
	}
	// 영상의 세로 픽셀이 홀수일 때 연산되지 않는다. -vf 옵션이 마지막으로 한번 붙어야 한다.
	args = append(args, []string{"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2"}...)
	args = append(args, []string{"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:0:0,setsar=1", adminSetting.MediaWidth, adminSetting.MediaHeight, adminSetting.MediaWidth, adminSetting.MediaHeight)}...)
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

// ProcessIesItem 함수는 ies 아이템을 연산한다.
func ProcessIesItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	err := SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessPdfItem 함수는 PDF 아이템을 연산한다.
func ProcessPdfItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// 아무 프로세스는 없지만 "done" 처리 해야한다. 그래야 프로세싱하지 않는다.
	err := SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessHwpItem 함수는 Hwp 아이템을 연산한다.
func ProcessHwpItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// 아무 프로세스는 없지만 "done" 처리 해야한다. 그래야 프로세싱하지 않는다.
	err := SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessPptItem 함수는 PPT 아이템을 연산한다.
func ProcessPptItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// 아무 프로세스는 없지만 "done" 처리 해야한다. 그래야 프로세싱하지 않는다.
	err := SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessUnrealItem 함수는 unreal 아이템을 연산한다.
func ProcessUnrealItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessModoItem 함수는 modo 아이템을 연산한다.
func ProcessModoItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}

// ProcessKatanaItem 함수는 katana 아이템을 연산한다.
func ProcessKatanaItem(client *mongo.Client, adminSetting Adminsetting, item Item) error {
	// thumbnail 폴더를 생성한다.
	err := SetStatus(client, item, "creating thumbnail directory")
	if err != nil {
		return err
	}
	err = genThumbDir(adminSetting, item)
	if err != nil {
		return err
	}
	// 썸네일 이미지를 생성한다.
	err = SetStatus(client, item, "creating thumbnail image")
	if err != nil {
		return err
	}
	err = genThumbImage(adminSetting, item)
	if err != nil {
		return err
	}
	// .ogg 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .ogg media")
	if err != nil {
		return err
	}
	err = genThumbOggMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mov 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mov media")
	if err != nil {
		return err
	}
	err = genThumbMovMedia(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	// .mp4 썸네일 동영상을 생성한다.
	err = SetStatus(client, item, "creating .mp4 media")
	if err != nil {
		return err
	}
	err = genThumbMp4Media(adminSetting, item) // FFmpeg는 확장자에 따라 옵션이 다양하거나 호환되지 않는다. 포멧별로 분리한다.
	if err != nil {
		return err
	}
	err = SetStatus(client, item, "done")
	if err != nil {
		return err
	}
	return nil
}
