package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleEditMaya(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		ID          primitive.ObjectID `json:"id" bson:"id"`
		ItemType    string             `json:"itemtype" bson:"itemtype"`
		Author      string             `json:"author" bson:"author"`
		Description string             `json:"description" bson:"description"`
		Tags        []string           `json:"tags" bson:"tags"`
		Attributes  map[string]string  `json:"attributes" bson:"attributes"`
		Token
	}
	q := r.URL.Query()
	itemtype := q.Get("itemtype")
	id := q.Get("id")
	if itemtype == "" {
		http.Error(w, "URL에 itemtype을 입력해주세요", http.StatusBadRequest)
		return
	}
	if id == "" {
		http.Error(w, "URL에 id를 입력해주세요", http.StatusBadRequest)
		return
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMonogDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := SearchItem(client, itemtype, id)

	rcp := recipe{
		ID:          item.ID,
		ItemType:    item.ItemType,
		Author:      item.Author,
		Description: item.Description,
		Tags:        item.Tags,
		Attributes:  item.Attributes,
		Token:       token,
	}

	err = TEMPLATES.ExecuteTemplate(w, "editmaya", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//handleEditMayaSubmit 함수는 maya아이템을 수정하는 페이지에서 UPDATE버튼을 누르면 작동하는 함수다.
func handleEditMayaSubmit(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	id := r.FormValue("id")
	itemtype := r.FormValue("itemtype")
	attrNum, err := strconv.Atoi(r.FormValue("attributesNum"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	attr := make(map[string]string)
	for i := 0; i < attrNum; i++ {
		key := r.FormValue(fmt.Sprintf("key%d", i))
		value := r.FormValue(fmt.Sprintf("value%d", i))
		attr[key] = value
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMonogDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item, err := SearchItem(client, itemtype, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.Author = r.FormValue("author")
	item.Description = r.FormValue("description")
	item.Tags = SplitBySpace(r.FormValue("tags"))
	item.Attributes = attr
	err = UpdateItem(client, itemtype, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/editmaya-success", http.StatusSeeOther)
}

func handleEditMayaSuccess(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "editmaya-success", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
