package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

func webserver() {
	// 웹주소 설정
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/search", handleSearch)
	// 웹서버 실행
	http.ListenAndServe(*flagHTTPPort, nil)

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("dotori"))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("add page"))
}

// handleSearch는 URL을 통해 query를 할 수 있게 해주는 함수입니다.
func handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	itemType := q.Get("itemtype")

	log.Println(itemType)

	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	items, err := allItems(session, itemType)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
