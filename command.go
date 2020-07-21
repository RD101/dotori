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
