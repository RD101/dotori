package main

import (
	"net/http"
)

func webserver() {
	// 웹주소 설정
	http.HandleFunc("/", handleIndex)
	// 웹서버 실행
	http.ListenAndServe(*flagHTTPPort, nil)

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("dotori"))
}
