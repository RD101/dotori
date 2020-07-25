package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func addMayaItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = SplitBySpace(*flagTag)
	i.Attributes = StringToMap(*flagAttributes)
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	currentTime := time.Now()
	i.CreateTime = currentTime.Format("2006-01-02 15:04:05")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크. 유효한 파일인지 체크.
	// 있으면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크. 유효한 파일인지 체크.
	// 있으면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. 데이터
	// 데이터 경로에 실재 파일이 존재하는지 체크. 유효한 파일인지 체크.
	// 있으면 OutputData 경로로 복사하기
	err = copyFile(*flagInputDataPath, i.OutputDataPath)
	if err != nil {
		log.Fatal(err)
	}
	// DataUploaded true로 바꾸기
	i.DataUploaded = true

	// 다 잘 업로드 됐으면 status바꾸기
	i.Status = "fileuploaded"

	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Print(err)
	}
}

func rmItemCmd() {
	if *flagItemID == "" {
		log.Fatal("id가 빈 문자열 입니다")
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	err = RmItem(client, *flagItemID)
	if err != nil {
		log.Print(err)
	}
}
