package main

import (
	"context"
	"fmt"
	"time"

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
	fmt.Println(item)
	return nil
}
