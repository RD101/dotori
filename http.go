package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/html/vfstemplate"
	"gopkg.in/mgo.v2"
)

// LoadTemplates 함수는 템플릿을 로딩합니다.
func LoadTemplates() (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	t, err := vfstemplate.ParseGlob(assets, t, "/template/*.html")
	return t, err
}

var funcMap = template.FuncMap{
	"Tags2str": Tags2str,
}

func webserver() {
	// 템플릿 로딩을 위해서 vfs(가상파일시스템)을 로딩합니다.
	vfsTemplate, err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}
	TEMPLATES = vfsTemplate
	// 리소스 로딩
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets)))

	// 웹주소 설정
	http.HandleFunc("/", handleSearch)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/search-submit", handleSearchSubmit)

	// Maya
	http.HandleFunc("/addmaya", handleAddMaya)
	http.HandleFunc("/addmaya-process", handleAddMayaProcess)
	http.HandleFunc("/upload-maya", handleUploadMaya)
	http.HandleFunc("/editmaya", handleEditMaya)
	http.HandleFunc("/editmaya-submit", handleEditMayaSubmit)
	http.HandleFunc("/editmaya-success", handleEditMayaSuccess)

	// 앞으로 정리할 것
	http.HandleFunc("/addhoudini", handleAddHoudini)
	http.HandleFunc("/addblender", handleAddBlender)
	http.HandleFunc("/addabc", handleAddABC)
	http.HandleFunc("/addusd", handleAddUSD)

	// Admin
	http.HandleFunc("/setlibrarypath", handleSetLibraryPath)
	// Help
	http.HandleFunc("/help", handleHelp)

	// REST API
	http.HandleFunc("/api/item", handleAPIItem)
	// 웹서버 실행
	err = http.ListenAndServe(*flagHTTPPort, nil)
	if err != nil {
		log.Fatal(err)
	}

}

// handleSearch는 URL을 통해 query를 할 수 있게 해주는 함수입니다.
func handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	itemType := q.Get("itemtype")
	searchword := q.Get("searchword")
	if itemType == "" {
		itemType = "maya"
	}
	type recipe struct {
		Items      []Item
		Searchword string
		ItemType   string
		TotalNum   int
	}
	rcp := recipe{}
	rcp.Searchword = searchword
	rcp.ItemType = itemType
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	totalNum, items, err := SearchPage(session, itemType, searchword, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Items = items
	rcp.TotalNum = totalNum
	err = TEMPLATES.ExecuteTemplate(w, "index", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSearchSubmit(w http.ResponseWriter, r *http.Request) {
	itemType := r.FormValue("itemtype")
	searchword := r.FormValue("searchword")
	http.Redirect(w, r, fmt.Sprintf("/search?itemtype=%s&searchword=%s", itemType, searchword), http.StatusSeeOther)
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "help", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
