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
	rootCategoryID := q.Get("rootcategoryid")
	subCategoryID := q.Get("subcategoryid")
	page := PageToString(q.Get("page"))
	if page == "" {
		page = "1"
	}
	type recipe struct {
		Items          []Item
		Searchword     string
		ItemType       string
		RootCategoryID string
		SubCategoryID  string
		TotalNum       int64
		CurrentPage    int64
		TotalPage      int64
		Pages          []int64
		Token
		User           User
		Adminsetting   Adminsetting
		RootCategories []Category
		SubCategories  []Category
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

	rcp.RootCategories, err = GetRootCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rootCategoryID != "" {
		rcp.RootCategoryID = rootCategoryID
		c, err := GetCategory(client, rootCategoryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		searchword += " categories:" + c.Name

		rcp.SubCategories, err = GetSubCategories(client, rootCategoryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if subCategoryID != "" {
		rcp.SubCategoryID = subCategoryID
		c, err := GetCategory(client, subCategoryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		searchword += " categories:" + c.Name
	}
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
	// 10 페이지씩 보이도록 Pages를 설정한다.
	pageCount := (rcp.CurrentPage - 1) / 10
	for i := 0; i < 10; i++ {
		rcp.Pages = append(rcp.Pages, int64(i)+pageCount*10+1)
		if rcp.Pages[i] == rcp.TotalPage {
			break
		}
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
	rootCategoryID := r.FormValue("searchbox-rootcategory-id")
	subCategoryID := r.FormValue("searchbox-subcategory-id")
	page := PageToString(r.FormValue("page"))
	http.Redirect(w, r, fmt.Sprintf("/search?itemtype=%s&searchword=%s&page=%s&rootcategoryid=%s&subcategoryid=%s", itemType, searchword, page, rootCategoryID, subCategoryID), http.StatusSeeOther)
}
