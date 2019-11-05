package main

import (
	"net/http"
)

func webserver() {
	// 웹주소 설정
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
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
