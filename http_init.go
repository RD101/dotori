package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleInit는 URL을 통해 query를 할 수 있게 해주는 함수입니다.
func handleInit(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	q := r.URL.Query()
	itemType := q.Get("itemtype")
	searchword := q.Get("searchword")
	type recipe struct {
		RecentlyCreateItems []Item
		TopUsingItems       []Item
		RecentlyTagItems    []Item
		Token
		Adminsetting Adminsetting
		Searchword   string
		ItemType     string
		TotalNum     int64
		AllItemCount string
		User         User
	}
	rcp := recipe{}
	rcp.Token = token
	rcp.Searchword = searchword
	rcp.ItemType = itemType
	rcp.TotalNum = 0

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

	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting

	rcp.User, err = GetUser(client, token.ID) // user 정보 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	num, err := AllItemsCount(client) // 전체아이템의 개수를 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.AllItemCount = humanize.Comma(num) // 숫자를 1000단위마다 comma를 찍음(string형으로 변경)

	RecentlyTagItems, err := RecentlyCreateItems(client, 20) // 최근생성된 20개의 아이템들을 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.RecentlyTagItems = RecentlyTagItems

	RecentlyCreateItems, err := RecentlyCreateItems(client, 100) // 최근생성된 100개의 아이템들을 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.RecentlyCreateItems = RecentlyCreateItems

	TopUsingItems, err := TopUsingItems(client, 20) // 사용률이 높은 20개의 아이템들을 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TopUsingItems = TopUsingItems

	err = TEMPLATES.ExecuteTemplate(w, "initPage", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
