package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	fmt.Println("GetReadyItem 완료")
	// 썸네일 이미지를 생성한다.
	err = genThumbImage(item)
	if err != nil {
		return err
	}
	fmt.Println("genThumbImage 완료")
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
	path := item.InputThumbnailImgPath
	// 변환할 이미지를 가져온다.
	target, err := imaging.Open(path)
	if err != nil {
		return err
	}
	// Resize the cropped image to width = 200px preserving the aspect ratio.
	result := imaging.Fill(target, 320, 180, imaging.Center, imaging.Lanczos)
	// 저장할 경로를 생성
	err = os.MkdirAll(filepath.Dir(item.OutputThumbnailPngPath), os.FileMode(0777))
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
