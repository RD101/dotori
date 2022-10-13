package main

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleItemProcess(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Items []Item
		Token
		Adminsetting   Adminsetting
		StorageTitle   string
		StoragePercent int64
		StorageLevel   int64
		User           User
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp := recipe{}
	// 완료되지 않은 아이템을 가져온다
	rcp.Items, err = GetUndoneItem(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Token = token
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting

	rcp.User, err = GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ds, err := DiskCheck()
	if err != nil {
		rcp.StorageTitle = "Storage Usage (Please set RootPath)"
		rcp.StoragePercent = 0
		rcp.StorageLevel = 0
	} else {
		rcp.StorageTitle = "Storage Usage"
		rcp.StoragePercent = int64((float64(ds.Used) / float64(ds.All)) * 100)
		rcp.StorageLevel = rcp.StoragePercent / 10
	}
	err = TEMPLATES.ExecuteTemplate(w, "item-process", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
