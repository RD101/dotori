package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shurcooL/httpfs/html/vfstemplate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// LoadTemplates 함수는 템플릿을 로딩합니다.
func LoadTemplates() (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	t, err := vfstemplate.ParseGlob(assets, t, "/template/*.html")
	return t, err
}

var funcMap = template.FuncMap{
	"Tags2str":     Tags2str,
	"add":          add,
	"PreviousPage": PreviousPage,
	"NextPage":     NextPage,
	"RmRootpath":   RmRootpath,
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
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/storage/", http.StripPrefix("/storage/", http.FileServer(http.Dir(adminsetting.Rootpath))))
	// 웹주소 설정
	http.HandleFunc("/", handleSearch)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/search-submit", handleSearchSubmit)

	// Maya
	http.HandleFunc("/addmaya-submit", handleAddMayaSubmit)
	http.HandleFunc("/addmaya-item", handleAddMayaItem)
	http.HandleFunc("/addmaya-file", handleAddMayaFile)
	http.HandleFunc("/uploadmaya-item", handleUploadMayaItem)
	http.HandleFunc("/uploadmaya-file", handleUploadMayaFile)
	http.HandleFunc("/uploadmaya-checkdata", handleUploadMayaCheckData)
	http.HandleFunc("/addmaya-success", handleAddMayaSuccess)
	http.HandleFunc("/editmaya", handleEditMaya)
	http.HandleFunc("/editmaya-submit", handleEditMayaSubmit)
	http.HandleFunc("/editmaya-success", handleEditMayaSuccess)

	// source
	http.HandleFunc("/addsource", handleAddSource)
	http.HandleFunc("/addsource-item", handleAddSourceItem)

	// nuke
	http.HandleFunc("/addnuke", handleAddNuke)
	http.HandleFunc("/addnuke-process", handleAddNukeProcess)
	http.HandleFunc("/upload-nuke", handleUploadNuke)

	// Houdini
	http.HandleFunc("/addhoudini", handleAddHoudini)
	http.HandleFunc("/addhoudini-process", handleAddHoudiniProcess)
	http.HandleFunc("/upload-houdini", handleUploadHoudini)

	// Alembic
	http.HandleFunc("/addalembic", handleAddAlembic)
	http.HandleFunc("/addalembic-process", handleAddAlembicProcess)
	http.HandleFunc("/upload-alembic", handleUploadAlembic)

	// 앞으로 정리할 것
	http.HandleFunc("/addblender", handleAddBlender)
	http.HandleFunc("/addusd", handleAddUSD)

	// Admin
	http.HandleFunc("/adminsetting", handleAdminSetting)
	http.HandleFunc("/adminsetting-submit", handleAdminSettingSubmit)
	http.HandleFunc("/adminsetting-success", handleAdminSettingSuccess)

	// Process
	http.HandleFunc("/item-process", handleItemProcess)
	http.HandleFunc("/cleanup-db", handleCleanUpDB)

	// Download
	http.HandleFunc("/download-item", handleDownloadItem)

	// Help
	http.HandleFunc("/help", handleHelp)

	// User
	http.HandleFunc("/profile", handleProfile)
	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/signup-submit", handleSignupSubmit)
	http.HandleFunc("/signup-success", handleSignupSuccess)
	http.HandleFunc("/signin", handleSignin)
	http.HandleFunc("/signin-submit", handleSigninSubmit)
	http.HandleFunc("/signout", handleSignOut)

	// REST API
	http.HandleFunc("/api/item", handleAPIItem)
	http.HandleFunc("/api/search", handleAPISearch)
	http.HandleFunc("/api/adminsetting", handleAPIAdminSetting)
	http.HandleFunc("/api/usingrate", handleAPIUsingRate)
	http.HandleFunc("/api/rmitem", handleAPIRmItem)

	// 웹서버 실행
	if *flagHTTPPort == ":443" { // https ports
		if *flagCertFullchain == "" {
			log.Fatal("CertFullchanin 인증서 설정이 필요합니다.")
		}
		if *flagCertPrivkey == "" {
			log.Fatal("CertPrivkey 인증서 설정이 필요합니다.")
		}
		if _, err := os.Stat(*flagCertFullchain); os.IsNotExist(err) {
			log.Fatal(*flagCertFullchain + " 경로에 인증서 파일이 존재하지 않습니다")
		}
		if _, err := os.Stat(*flagCertFullchain); os.IsNotExist(err) {
			log.Fatal(*flagCertPrivkey + " 경로에 인증서 파일이 존재하지 않습니다")
		}
		err := http.ListenAndServeTLS(*flagHTTPPort, *flagCertFullchain, *flagCertPrivkey, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = http.ListenAndServe(*flagHTTPPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

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
	if itemType == "" {
		itemType = "maya"
	}
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

func handleHelp(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting Adminsetting
	}
	rcp := recipe{}
	rcp.Token = token
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
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
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "help", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleItemProcess(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		Items []Item
		Token
		Adminsetting     Adminsetting
		StorageClassName string
		StorageTitle     string
		StoragePercent   int64
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
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
	rcp := recipe{}
	// 완료되지 않은 아이템을 가져온다
	rcp.Items, err = GetOngoingProcess(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Token = token
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting
	ds, err := DiskCheck()
	if err != nil {
		rcp.StorageTitle = "Storage Usage (Please set RootPath)"
		rcp.StoragePercent = 0
		rcp.StorageClassName = "progress-bar bg-success"
	} else {
		rcp.StorageTitle = "Storage Usage"
		rcp.StoragePercent = int64((float64(ds.Used) / float64(ds.All)) * 100)
		num := rcp.StoragePercent / 10
		switch num {
		case 10:
		case 9:
			rcp.StorageClassName = "progress-bar bg-danger"
			break
		case 8:
			rcp.StorageClassName = "progress-bar bg-warning"
			break
		case 7:
			rcp.StorageClassName = "progress-bar bg-info"
			break
		default:
			rcp.StorageClassName = "progress-bar bg-success"
			break
		}
	}
	err = TEMPLATES.ExecuteTemplate(w, "item-process", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCleanUpDB(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		Items []Item
		Token
		Adminsetting Adminsetting
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
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
	rcp := recipe{}
	// 데이터가 모두 업로드되지 않은 아이템을 가져온다
	rcp.Items, err = GetIncompleteItems(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Token = token
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting
	err = TEMPLATES.ExecuteTemplate(w, "cleanup-db", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
