package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleSearch는 URL을 통해 query를 할 수 있게 해주는 함수입니다.
func handleSearch(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	q := r.URL.Query()
	itemType := q.Get("itemtype")
	searchword := q.Get("searchword")
	page := PageToString(q.Get("page"))
	if page == "" {
		page = "1"
	}
	type recipe struct {
		Items       []Item
		Searchword  string
		ItemType    string
		TotalNum    int64
		CurrentPage int64
		TotalPage   int64
		Pages       []int64
		Token
		User         User
		Adminsetting Adminsetting
	}
	rcp := recipe{}
	rcp.Searchword = searchword
	rcp.ItemType = itemType
	rcp.Token = token
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
	rcp.CurrentPage = PageToInt(page)
	totalPage, totalNum, items, err := SearchPage(client, itemType, searchword, rcp.CurrentPage, *flagPagenum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Items = items
	rcp.User, err = GetUser(client, token.ID) // user 정보 가져옴
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rcp.TotalNum = totalNum
	rcp.TotalPage = totalPage
	// Pages를 설정한다.
	rcp.Pages = make([]int64, totalPage) // page에 필요한 메모리를 미리 설정한다.
	for i := range rcp.Pages {
		rcp.Pages[i] = int64(i) + 1
	}
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting
	err = TEMPLATES.ExecuteTemplate(w, "index", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSearchSubmit(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	itemType := r.FormValue("itemtype")
	searchword := r.FormValue("searchword")
	page := PageToString(r.FormValue("page"))
	http.Redirect(w, r, fmt.Sprintf("/search?itemtype=%s&searchword=%s&page=%s", itemType, searchword, page), http.StatusSeeOther)
}
