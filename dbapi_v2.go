package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddItem 은 데이터베이스에 Item을 넣는 함수이다.
func AddItem(client *mongo.Client, i Item) error {
	collection := client.Database(*flagDBName).Collection(i.ItemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}

// GetItem 은 데이터베이스에 Item을 가지고 오는 함수이다.
func GetItem(client *mongo.Client, itemType, id string) (Item, error) {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result Item
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetAdminSetting 은 관리자 셋팅값을 가지고 온다.
func GetAdminSetting(client *mongo.Client) (Adminsetting, error) {
	//monotonic이 필수적인가? 필수적이라면 대응하는 기능은 무엇인가
	//session.SetMode(mgo.Monotonic, true)
	collection := client.Database(*flagDBName).Collection("setting.admin")
	var result Adminsetting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{"id": "setting.admin"}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNilDocument { // document가 존재하지 않는 경우, 즉 adminsetting이 없는 경우
			return Adminsetting{}, nil
		}
		return Adminsetting{}, err
	}
	return result, nil
}

// RmItem 는 컬렉션 이름과 id를 받아서, 해당 컬렉션에서 id가 일치하는 Item을 삭제한다.
func RmItem(client *mongo.Client, itemType, id string) error {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

// AllItems 는 DB에서 전체 아이템 정보를 가져오는 함수입니다.
func AllItems(client *mongo.Client, itemType string) ([]Item, error) {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var results []Item
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}
