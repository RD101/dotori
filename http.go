package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/html/vfstemplate"
	"gopkg.in/mgo.v2"
)

// LoadTemplates 함수는 템플릿을 로딩합니다.
func LoadTemplates() (*template.Template, error) {
	t := template.New("")
	t, err := vfstemplate.ParseGlob(assets, t, "/template/*.html")
	return t, err
}

func webserver() {
	// 템플릿 로딩을 위해서 vfs(가상파일시스템)을 로딩합니다.
	vfsTemplate, err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}
	TEMPLATES = vfsTemplate
	// 웹주소 설정
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/search", handleSearch)
	// 웹서버 실행
	err = http.ListenAndServe(*flagHTTPPort, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("add page"))
}

// handleSearch는 URL을 통해 query를 할 수 있게 해주는 함수입니다.
func handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	itemType := q.Get("itemtype")
	if itemType == "" {
		http.Error(w, "URL에 itemtype을 입력해주세요", http.StatusBadRequest)
		return
	}

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
