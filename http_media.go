package main

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleMediaData 함수는 영상데이터를 전송한다.
func handleMediaData(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	if id == "" {
		http.Error(w, "id 값이 빈 문자열 입니다", http.StatusInternalServerError)
		return
	}
	typ := q.Get("type")
	if !(typ == "mp4" || typ == "ogg" || typ == "mov" || typ == "png") {
		http.Error(w, "type 값은 mp4, ogg, mov, png 값만 지원합니다", http.StatusInternalServerError)
		return
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
	item, err := GetItem(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch typ {
	case "mp4":
		http.ServeFile(w, r, item.OutputThumbnailMp4Path)
		return
	case "mov":
		http.ServeFile(w, r, item.OutputThumbnailMovPath)
		return
	case "ogg":
		http.ServeFile(w, r, item.OutputThumbnailOggPath)
		return
	case "png":
		http.ServeFile(w, r, item.OutputThumbnailPngPath)
		return
	default:
		http.ServeFile(w, r, item.OutputThumbnailPngPath)
		return
	}
}
