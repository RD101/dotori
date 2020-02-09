package main

import (
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func handleEditMaya(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		ID          bson.ObjectId     `json:"id" bson:"id"`
		ItemType    string            `json:"itemtype" bson:"itemtype"`
		Author      string            `json:"author" bson:"author"`
		Description string            `json:"description" bson:"description"`
		Tags        []string          `json:"tags" bson:"tags"`
		Attributes  map[string]string `json:"attributes" bson:"attributes"`
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	item, err := SearchItem(session, itemtype, id)

	rcp := recipe{
		ID:          item.ID,
		ItemType:    item.ItemType,
		Author:      item.Author,
		Description: item.Description,
		Tags:        item.Tags,
		Attributes:  item.Attributes,
	}

	err = TEMPLATES.ExecuteTemplate(w, "editmaya", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEditMayaSubmit(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	item, err := SearchItem(session, itemtype, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.Author = r.FormValue("author")
	item.Description = r.FormValue("description")
	item.Tags = SplitBySpace(r.FormValue("tags"))
	item.Attributes = attr
	err = UpdateItem(session, itemtype, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/editmaya-success", http.StatusSeeOther)
}

func handleEditMayaSuccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "editmaya-success", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
