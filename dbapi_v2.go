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
