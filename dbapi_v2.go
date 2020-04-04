package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// AddItem 은 데이터베이스에 Item을 넣는 함수이다.
func AddItem(client *mongo.Client, i Item) error {
	collection := client.Database("dotori").Collection("maya")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}
