package main

import (
	"context"
	"net/http"
	"os"
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
		if err == mongo.ErrNoDocuments {
			http.Error(w, id+" ID를 가진 아이템이 존재하지 않습니다", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch typ {
	case "mp4":
		if _, err := os.Stat(item.OutputThumbnailMp4Path); os.IsNotExist(err) {
			http.Error(w, item.OutputThumbnailMp4Path+" 파일이 존재하지 않습니다", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, item.OutputThumbnailMp4Path)
		return
	case "mov":
		if _, err := os.Stat(item.OutputThumbnailMovPath); os.IsNotExist(err) {
			http.Error(w, item.OutputThumbnailMovPath+" 파일이 존재하지 않습니다", http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "video/quicktime")
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, item.OutputThumbnailMovPath)
		return
	case "ogg":
		if _, err := os.Stat(item.OutputThumbnailOggPath); os.IsNotExist(err) {
			http.Error(w, item.OutputThumbnailOggPath+" 파일이 존재하지 않습니다", http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "video/ogg")
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, item.OutputThumbnailOggPath)
		return
	case "png":
		if _, err := os.Stat(item.OutputThumbnailPngPath); os.IsNotExist(err) {
			http.Error(w, item.OutputThumbnailPngPath+" 파일이 존재하지 않습니다", http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "image/png")
		http.ServeFile(w, r, item.OutputThumbnailPngPath)
		return
	default:
		http.Error(w, "지원하지 않는 형식입니다", http.StatusNotFound)
		return
	}
}
