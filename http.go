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
	"LastLog":      LastLog,
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
	http.HandleFunc("/uploadfootage-checkdata", handleUploadFootageCheckData)
	http.HandleFunc("/addfootage-success", handleAddFootageSuccess)
	http.HandleFunc("/editfootage", handleEditFootage)
	http.HandleFunc("/editfootage-submit", handleEditFootageSubmit)
	http.HandleFunc("/editfootage-success", handleEditFootageSuccess)

	// Nuke
	http.HandleFunc("/addnuke", handleAddNuke)
	http.HandleFunc("/addnuke-item", handleAddNukeItem)
	http.HandleFunc("/uploadnuke-item", handlUploadNukeItem)
	http.HandleFunc("/addnuke-file", handleAddNukeFile)
	http.HandleFunc("/uploadnuke-file", handleUploadNukeFile)
	http.HandleFunc("/uploadnuke-checkdata", handleUploadNukeCheckData)
	http.HandleFunc("/addnuke-success", handleAddNukeSuccess)
	http.HandleFunc("/editnuke", handleEditNuke)
	http.HandleFunc("/editnuke-submit", handleEditNukeSubmit)
	http.HandleFunc("/editnuke-success", handleEditNukeSuccess)

	// Houdini
	http.HandleFunc("/addhoudini", handleAddHoudini)
	http.HandleFunc("/addhoudini-item", handleAddHoudiniItem)
	http.HandleFunc("/addhoudini-file", handleAddHoudiniFile)
	http.HandleFunc("/uploadhoudini-item", handleUploadHoudiniItem)
	http.HandleFunc("/uploadhoudini-file", handleUploadHoudiniFile)
	http.HandleFunc("/uploadhoudini-checkdata", handleUploadHoudiniCheckData)
	http.HandleFunc("/edithoudini", handleEditHoudini)
	http.HandleFunc("/edithoudini-submit", handleEditHoudiniSubmit)
	http.HandleFunc("/edithoudini-success", handleEditHoudiniSuccess)

	// Blender
	http.HandleFunc("/addblender", handleAddBlender)
	http.HandleFunc("/addblender-item", handleAddBlenderItem)
	http.HandleFunc("/addblender-file", handleAddBlenderFile)
	http.HandleFunc("/uploadblender-item", handleUploadBlenderItem)
	http.HandleFunc("/uploadblender-file", handleUploadBlenderFile)
	http.HandleFunc("/uploadblender-checkdata", handleUploadBlenderCheckData)
	http.HandleFunc("/addblender-success", handleAddBlenderSuccess)
	http.HandleFunc("/editblender", handleEditBlender)
	http.HandleFunc("/editblender-submit", handleEditBlenderSubmit)
	http.HandleFunc("/editblender-success", handleEditBlenderSuccess)

	// Alembic
	http.HandleFunc("/addalembic", handleAddAlembic)
	http.HandleFunc("/addalembic-item", handleAddAlembicItem)
	http.HandleFunc("/addalembic-file", handleAddAlembicFile)
	http.HandleFunc("/uploadalembic-item", handleUploadAlembicItem)
	http.HandleFunc("/uploadalembic-file", handleUploadAlembicFile)
	http.HandleFunc("/uploadalembic-checkdata", handleUploadAlembicCheckData)
	http.HandleFunc("/addalembic-success", handleAddAlembicSuccess)
	http.HandleFunc("/editalembic", handleEditAlembic)
	http.HandleFunc("/editalembic-submit", handleEditAlembicSubmit)
	http.HandleFunc("/editalembic-success", handleEditAlembicSuccess)

	// PixarUSD
	http.HandleFunc("/addusd", handleAddUSD)
	http.HandleFunc("/addusd-item", handleAddUSDItem)
	http.HandleFunc("/uploadusd-item", handleUploadUSDItem)
	http.HandleFunc("/addusd-file", handleAddUSDFile)
	http.HandleFunc("/uploadusd-file", handleUploadUSDFile)
	http.HandleFunc("/uploadusd-checkdata", handleUploadUSDCheckData)
	http.HandleFunc("/addusd-success", handleAddUSDSuccess)
	http.HandleFunc("/editusd", handleEditUSD)
	http.HandleFunc("/editusd-submit", handleEditUSDSubmit)
	http.HandleFunc("/editusd-success", handleEditUSDSuccess)

	// Sound
	http.HandleFunc("/addsound", handleAddSound)
	http.HandleFunc("/addsound-item", handleAddSoundItem)
	http.HandleFunc("/addsound-file", handleAddSoundFile)
	http.HandleFunc("/uploadsound-item", handleUploadSoundItem)
	http.HandleFunc("/uploadsound-file", handleUploadSoundFile)
	http.HandleFunc("/uploadsound-checkdata", handleUploadSoundCheckData)
	http.HandleFunc("/addsound-success", handleAddSoundSuccess)
	http.HandleFunc("/editsound", handleEditSound)
	http.HandleFunc("/editsound-submit", handleEditSoundSubmit)
	http.HandleFunc("/editsound-success", handleEditSoundSuccess)

	// pdf
	http.HandleFunc("/addpdf", handleAddPdf)
	http.HandleFunc("/addpdf-item", handleAddPdfItem)
	http.HandleFunc("/addpdf-file", handleAddPdfFile)
	http.HandleFunc("/uploadpdf-item", handleUploadPdfItem)
	http.HandleFunc("/uploadpdf-file", handleUploadPdfFile)
	http.HandleFunc("/uploadpdf-checkdata", handleUploadPdfCheckData)
	http.HandleFunc("/addpdf-success", handleAddPdfSuccess)
	http.HandleFunc("/editpdf", handleEditPdf)
	http.HandleFunc("/editpdf-submit", handleEditPdfSubmit)
	http.HandleFunc("/editpdf-success", handleEditPdfSuccess)

	// hwp
	http.HandleFunc("/addhwp", handleAddHwp)
	http.HandleFunc("/addhwp-item", handleAddHwpItem)
	http.HandleFunc("/addhwp-file", handleAddHwpFile)
	http.HandleFunc("/uploadhwp-item", handleUploadHwpItem)
	http.HandleFunc("/uploadhwp-file", handleUploadHwpFile)
	http.HandleFunc("/uploadhwp-checkdata", handleUploadHwpCheckData)
	http.HandleFunc("/addhwp-success", handleAddHwpSuccess)
	http.HandleFunc("/edithwp", handleEditHwp)
	http.HandleFunc("/edithwp-submit", handleEditHwpSubmit)
	http.HandleFunc("/edithwp-success", handleEditHwpSuccess)

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
	if *flagCertFullchain != "" || *flagCertPrivkey != "" {
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
