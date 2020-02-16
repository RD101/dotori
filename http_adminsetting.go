package main

import (
	"net/http"
)

// handleAdminSetting 함수는 Admin 설정 페이지로 이동한다.
func handleAdminSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "adminsetting", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
