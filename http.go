package main

import (
	"context"
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

// 이 파일에서는 최대한 핸들러만 선언합니다.

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
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
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
	http.HandleFunc("/addmaya", handleAddMaya)
	http.HandleFunc("/addmaya-item", handleAddMayaItem)
	http.HandleFunc("/addmaya-file", handleAddMayaFile)
	http.HandleFunc("/uploadmaya-item", handleUploadMayaItem)
	http.HandleFunc("/uploadmaya-file", handleUploadMayaFile)
	http.HandleFunc("/uploadmaya-checkdata", handleUploadMayaCheckData)
	http.HandleFunc("/addmaya-success", handleAddMayaSuccess)
	http.HandleFunc("/editmaya", handleEditMaya)
	http.HandleFunc("/editmaya-submit", handleEditMayaSubmit)
	http.HandleFunc("/editmaya-success", handleEditMayaSuccess)

	// Footage
	http.HandleFunc("/addfootage", handleAddFootage)
	http.HandleFunc("/addfootage-item", handleAddFootageItem)
	http.HandleFunc("/uploadfootage-item", handleUploadFootageItem)
	http.HandleFunc("/addfootage-file", handleAddFootageFile)
	http.HandleFunc("/uploadfootage-file", handleUploadFootageFile)
	http.HandleFunc("/editfootage", handleEditFootage)
	http.HandleFunc("/editfootage-submit", handleEditFootageSubmit)

	// Nuke
	http.HandleFunc("/addnuke", handleAddNuke)
	http.HandleFunc("/addnuke-item", handleAddNukeItem)
	http.HandleFunc("/uploadnuke-item", handlUploadNukeItem)
	http.HandleFunc("/addnuke-file", handleAddNukeFile)
	http.HandleFunc("/uploadnuke-file", handleUploadNukeFile)
	http.HandleFunc("/uploadnuke-checkdata", handleUploadNukeCheckData)
	http.HandleFunc("/addnuke-success", handleAddNukeSuccess)

	// Houdini
	http.HandleFunc("/addhoudini", handleAddHoudini)
	http.HandleFunc("/addhoudini-process", handleAddHoudiniProcess)
	http.HandleFunc("/upload-houdini", handleUploadHoudini)

	// Alembic
	http.HandleFunc("/addalembic", handleAddAlembic)
	http.HandleFunc("/addalembic-process", handleAddAlembicProcess)
	http.HandleFunc("/upload-alembic", handleUploadAlembic)

	// Blender
	http.HandleFunc("/addblender", handleAddBlender)

	// PixarUSD
	http.HandleFunc("/addusd", handleAddUSD)

	// Admin
	http.HandleFunc("/adminsetting", handleAdminSetting)
	http.HandleFunc("/adminsetting-submit", handleAdminSettingSubmit)
	http.HandleFunc("/adminsetting-success", handleAdminSettingSuccess)

	// Process
	http.HandleFunc("/item-process", handleItemProcess)
	http.HandleFunc("/cleanup-db", handleCleanUpDB)
	http.HandleFunc("/cleanup-db-submit", handleCleanUpDBSubmit)
	http.HandleFunc("/cleanup-db-success", handleCleanUpDBSuccess)

	// Download
	http.HandleFunc("/download-item", handleDownloadItem)

	// Help
	http.HandleFunc("/help", handleHelp)

	// Error
	http.HandleFunc("/error-ocio", handleErrorOCIO)

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
