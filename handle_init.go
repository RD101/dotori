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
	type recipe struct {
		RecentlyCreateItems []Item
		TopUsingItems       []Item
		RecentTags          []string // 최근 등록된 태그 리스트
		Token
		Adminsetting         Adminsetting
		Searchword           string
		ItemType             string
		TotalNum             int64
		AllItemCount         string
		RecentlyTotalItemNum int64
		TopUsingTotalItemNum int64
		User                 User
		RootCategoryID       string
		SubCategoryID        string
		RootCategories       []Category
		SubCategories        []Category
	}
	rcp := recipe{}
	rcp.Token = token
	rcp.ItemType = "" // search sortList all
	rcp.TotalNum = 0  // search button totalNum

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

	num, err := GetAllItemsNum(client) // 전체아이템의 개수를 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.AllItemCount = humanize.Comma(num) // 숫자를 1000단위마다 comma를 찍음(string형으로 변경)

	rcp.RecentlyTotalItemNum = int64(rcp.User.NewsNum)
	rcp.TopUsingTotalItemNum = int64(rcp.User.TopNum)

	RecentlyTagItems, err := GetRecentlyCreatedItems(client, 20, 1) // 최근생성된 20개의 아이템들을 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.RecentTags = ItemsTagsDeduplication(RecentlyTagItems) // 중복 태그를 정리함

	RecentlyCreateItems, err := GetRecentlyCreatedItems(client, int64(rcp.User.NewsNum), 1) // 최근생성된 100개의 아이템들을 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.RecentlyCreateItems = RecentlyCreateItems

	TopUsingItems, err := GetTopUsingItems(client, int64(rcp.User.TopNum), 1) // 사용률이 높은 20개의 아이템들을 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TopUsingItems = TopUsingItems
	rcp.RootCategories, err = GetRootCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = TEMPLATES.ExecuteTemplate(w, "initPage", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
