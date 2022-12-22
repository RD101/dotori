package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	"Tags2str":       Tags2str,
	"add":            add,
	"sub":            sub,
	"mod":            mod,
	"divCeil":        divCeil,
	"PreviousPage":   PreviousPage,
	"NextPage":       NextPage,
	"RmRootpath":     RmRootpath,
	"LastLog":        LastLog,
	"SplitTimeData":  SplitTimeData,
	"ItemListLength": ItemListLength,
	"IntToSlice":     IntToSlice,
}

func helpMethodOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
}

func webserver() {
	// 템플릿 로딩을 위해서 vfs(가상파일시스템)을 로딩합니다.
	vfsTemplate, err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}
	TEMPLATES = vfsTemplate

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
		log.Println(*flagMongoDBURI + " 에 mongoDB가 실행되고 있는지 체크해주세요")
		log.Fatal(err)
	}
	// 웹주소 설정
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assets)))
	r.HandleFunc("/", handleInit)
	r.HandleFunc("/mediadata", handleMediaData)
	r.HandleFunc("/search", handleSearch)
	r.HandleFunc("/search-submit", handleSearchSubmit)
	r.HandleFunc("/tags", handleTags)

	// Rename
	r.HandleFunc("/rename", handleRename).Methods(http.MethodGet)
	r.HandleFunc("/rename/{id}", handleRename).Methods(http.MethodGet)
	r.HandleFunc("/api/searchfile", handleAPISearchFile)
	r.HandleFunc("/api/rename", handleAPIRename).Methods(http.MethodPost, http.MethodOptions)

	// Maya
	r.HandleFunc("/addmaya", handleAddMaya)
	r.HandleFunc("/addmaya-item", handleAddMayaItem)
	r.HandleFunc("/addmaya-file", handleAddMayaFile)
	r.HandleFunc("/uploadmaya-item", handleUploadMayaItem)
	r.HandleFunc("/uploadmaya-file", handleUploadMayaFile)
	r.HandleFunc("/uploadmaya-checkdata", handleUploadMayaCheckData)
	r.HandleFunc("/addmaya-success", handleAddMayaSuccess)
	r.HandleFunc("/editmaya", handleEditMaya)
	r.HandleFunc("/editmaya-submit", handleEditMayaSubmit)
	r.HandleFunc("/editmaya-success", handleEditMayaSuccess)

	// 3dsmax
	r.HandleFunc("/addmax", handleAddMax)
	r.HandleFunc("/addmax-item", handleAddMaxItem)
	r.HandleFunc("/addmax-file", handleAddMaxFile)
	r.HandleFunc("/uploadmax-item", handleUploadMaxItem)
	r.HandleFunc("/uploadmax-file", handleUploadMaxFile)
	r.HandleFunc("/uploadmax-checkdata", handleUploadMaxCheckData)
	r.HandleFunc("/addmax-success", handleAddMaxSuccess)
	r.HandleFunc("/editmax", handleEditMax)
	r.HandleFunc("/editmax-submit", handleEditMaxSubmit)
	r.HandleFunc("/editmax-success", handleEditMaxSuccess)

	// Fusion360
	r.HandleFunc("/addfusion360", handleAddFusion360)
	r.HandleFunc("/addfusion360-item", handleAddFusion360Item)
	r.HandleFunc("/addfusion360-file", handleAddFusion360File)
	r.HandleFunc("/uploadfusion360-item", handleUploadFusion360Item)
	r.HandleFunc("/uploadfusion360-file", handleUploadFusion360File)
	r.HandleFunc("/uploadfusion360-checkdata", handleUploadFusion360CheckData)
	r.HandleFunc("/addfusion360-success", handleAddFusion360Success)
	r.HandleFunc("/editfusion360", handleEditFusion360)
	r.HandleFunc("/editfusion360-submit", handleEditFusion360Submit)
	r.HandleFunc("/editfusion360-success", handleEditFusion360Success)

	// Footage
	r.HandleFunc("/addfootage", handleAddFootage)
	r.HandleFunc("/addfootage-item", handleAddFootageItem)
	r.HandleFunc("/addfootage-items", handleAddFootageItems)
	r.HandleFunc("/uploadfootage-item", handleUploadFootageItem)
	r.HandleFunc("/addfootage-file", handleAddFootageFile)
	r.HandleFunc("/uploadfootage-file", handleUploadFootageFile)
	r.HandleFunc("/uploadfootage-checkdata", handleUploadFootageCheckData)
	r.HandleFunc("/addfootage-success", handleAddFootageSuccess)
	r.HandleFunc("/editfootage", handleEditFootage)
	r.HandleFunc("/editfootage-submit", handleEditFootageSubmit)
	r.HandleFunc("/editfootage-success", handleEditFootageSuccess)
	r.HandleFunc("/api/searchfootages", handleAPISearchFootages)
	r.HandleFunc("/api/addfootage", handleAPIAddFootage)

	// Clip
	r.HandleFunc("/addclip", handleAddClip)
	r.HandleFunc("/addclip-item", handleAddClipItem)
	r.HandleFunc("/addclip-items", handleAddClipItems)
	r.HandleFunc("/uploadclip-item", handleUploadClipItem)
	r.HandleFunc("/addclip-file", handleAddClipFile)
	r.HandleFunc("/uploadclip-file", handleUploadClipFile)
	r.HandleFunc("/uploadclip-checkdata", handleUploadClipCheckData)
	r.HandleFunc("/addclip-success", handleAddClipSuccess)
	r.HandleFunc("/editclip", handleEditClip)
	r.HandleFunc("/editclip-submit", handleEditClipSubmit)
	r.HandleFunc("/editclip-success", handleEditClipSuccess)
	r.HandleFunc("/api/searchclips", handleAPISearchClips)
	r.HandleFunc("/api/addclip", handleAPIAddClip)

	// HDRI
	r.HandleFunc("/addhdri", handleAddHDRI)
	r.HandleFunc("/addhdri-item", handleAddHDRIItem)
	r.HandleFunc("/uploadhdri-item", handleUploadHDRIItem)
	r.HandleFunc("/addhdri-file", handleAddHDRIFile)
	r.HandleFunc("/uploadhdri-file", handleUploadHDRIFile)
	r.HandleFunc("/uploadhdri-checkdata", handleUploadHDRICheckData)
	r.HandleFunc("/addhdri-success", handleAddHDRISuccess)
	r.HandleFunc("/edithdri", handleEditHDRI)
	r.HandleFunc("/edithdri-submit", handleEditHDRISubmit)
	r.HandleFunc("/edithdri-success", handleEditHDRISuccess)

	// Texture
	r.HandleFunc("/addtexture", handleAddTexture)
	r.HandleFunc("/addtexture-item", handleAddTextureItem)
	r.HandleFunc("/uploadtexture-item", handleUploadTextureItem)
	r.HandleFunc("/addtexture-file", handleAddTextureFile)
	r.HandleFunc("/uploadtexture-file", handleUploadTextureFile)
	r.HandleFunc("/uploadtexture-checkdata", handleUploadTextureCheckData)
	r.HandleFunc("/addtexture-success", handleAddTextureSuccess)
	r.HandleFunc("/edittexture", handleEditTexture)
	r.HandleFunc("/edittexture-submit", handleEditTextureSubmit)
	r.HandleFunc("/edittexture-success", handleEditTextureSuccess)

	// Nuke
	r.HandleFunc("/addnuke", handleAddNuke)
	r.HandleFunc("/addnuke-item", handleAddNukeItem)
	r.HandleFunc("/uploadnuke-item", handleUploadNukeItem)
	r.HandleFunc("/addnuke-file", handleAddNukeFile)
	r.HandleFunc("/uploadnuke-file", handleUploadNukeFile)
	r.HandleFunc("/uploadnuke-checkdata", handleUploadNukeCheckData)
	r.HandleFunc("/addnuke-success", handleAddNukeSuccess)
	r.HandleFunc("/editnuke", handleEditNuke)
	r.HandleFunc("/editnuke-submit", handleEditNukeSubmit)
	r.HandleFunc("/editnuke-success", handleEditNukeSuccess)
	r.HandleFunc("/api/nukepath/{id}", getNukePathHandler).Methods("GET")

	// Houdini
	r.HandleFunc("/addhoudini", handleAddHoudini)
	r.HandleFunc("/addhoudini-item", handleAddHoudiniItem)
	r.HandleFunc("/addhoudini-file", handleAddHoudiniFile)
	r.HandleFunc("/uploadhoudini-item", handleUploadHoudiniItem)
	r.HandleFunc("/uploadhoudini-file", handleUploadHoudiniFile)
	r.HandleFunc("/uploadhoudini-checkdata", handleUploadHoudiniCheckData)
	r.HandleFunc("/addhoudini-success", handleAddHoudiniSuccess)
	r.HandleFunc("/edithoudini", handleEditHoudini)
	r.HandleFunc("/edithoudini-submit", handleEditHoudiniSubmit)
	r.HandleFunc("/edithoudini-success", handleEditHoudiniSuccess)

	// Blender
	r.HandleFunc("/addblender", handleAddBlender)
	r.HandleFunc("/addblender-item", handleAddBlenderItem)
	r.HandleFunc("/addblender-file", handleAddBlenderFile)
	r.HandleFunc("/uploadblender-item", handleUploadBlenderItem)
	r.HandleFunc("/uploadblender-file", handleUploadBlenderFile)
	r.HandleFunc("/uploadblender-checkdata", handleUploadBlenderCheckData)
	r.HandleFunc("/addblender-success", handleAddBlenderSuccess)
	r.HandleFunc("/editblender", handleEditBlender)
	r.HandleFunc("/editblender-submit", handleEditBlenderSubmit)
	r.HandleFunc("/editblender-success", handleEditBlenderSuccess)

	// Alembic
	r.HandleFunc("/addalembic", handleAddAlembic)
	r.HandleFunc("/addalembic-item", handleAddAlembicItem)
	r.HandleFunc("/addalembic-file", handleAddAlembicFile)
	r.HandleFunc("/uploadalembic-item", handleUploadAlembicItem)
	r.HandleFunc("/uploadalembic-file", handleUploadAlembicFile)
	r.HandleFunc("/uploadalembic-checkdata", handleUploadAlembicCheckData)
	r.HandleFunc("/addalembic-success", handleAddAlembicSuccess)
	r.HandleFunc("/editalembic", handleEditAlembic)
	r.HandleFunc("/editalembic-submit", handleEditAlembicSubmit)
	r.HandleFunc("/editalembic-success", handleEditAlembicSuccess)

	// PixarUSD
	r.HandleFunc("/addusd", handleAddUSD)
	r.HandleFunc("/addusd-item", handleAddUSDItem)
	r.HandleFunc("/uploadusd-item", handleUploadUSDItem)
	r.HandleFunc("/addusd-file", handleAddUSDFile)
	r.HandleFunc("/uploadusd-file", handleUploadUSDFile)
	r.HandleFunc("/uploadusd-checkdata", handleUploadUSDCheckData)
	r.HandleFunc("/addusd-success", handleAddUSDSuccess)
	r.HandleFunc("/editusd", handleEditUSD)
	r.HandleFunc("/editusd-submit", handleEditUSDSubmit)
	r.HandleFunc("/editusd-success", handleEditUSDSuccess)

	// OpenVDB
	r.HandleFunc("/addopenvdb", handleAddOpenVDB)
	r.HandleFunc("/addopenvdb-item", handleAddOpenVDBItem)
	r.HandleFunc("/uploadopenvdb-item", handleUploadOpenVDBItem)
	r.HandleFunc("/addopenvdb-file", handleAddOpenVDBFile)
	r.HandleFunc("/uploadopenvdb-file", handleUploadOpenVDBFile)
	r.HandleFunc("/uploadopenvdb-checkdata", handleUploadOpenVDBCheckData)
	r.HandleFunc("/addopenvdb-success", handleAddOpenVDBSuccess)
	r.HandleFunc("/editopenvdb", handleEditOpenVDB)
	r.HandleFunc("/editopenvdb-submit", handleEditOpenVDBSubmit)
	r.HandleFunc("/editopenvdb-success", handleEditOpenVDBSuccess)

	// Sound
	r.HandleFunc("/addsound", handleAddSound)
	r.HandleFunc("/addsound-item", handleAddSoundItem)
	r.HandleFunc("/addsound-file", handleAddSoundFile)
	r.HandleFunc("/uploadsound-item", handleUploadSoundItem)
	r.HandleFunc("/uploadsound-file", handleUploadSoundFile)
	r.HandleFunc("/uploadsound-checkdata", handleUploadSoundCheckData)
	r.HandleFunc("/addsound-success", handleAddSoundSuccess)
	r.HandleFunc("/editsound", handleEditSound)
	r.HandleFunc("/editsound-submit", handleEditSoundSubmit)
	r.HandleFunc("/editsound-success", handleEditSoundSuccess)

	// pdf
	r.HandleFunc("/addpdf", handleAddPdf)
	r.HandleFunc("/addpdf-item", handleAddPdfItem)
	r.HandleFunc("/addpdf-file", handleAddPdfFile)
	r.HandleFunc("/uploadpdf-item", handleUploadPdfItem)
	r.HandleFunc("/uploadpdf-file", handleUploadPdfFile)
	r.HandleFunc("/uploadpdf-checkdata", handleUploadPdfCheckData)
	r.HandleFunc("/addpdf-success", handleAddPdfSuccess)
	r.HandleFunc("/editpdf", handleEditPdf)
	r.HandleFunc("/editpdf-submit", handleEditPdfSubmit)
	r.HandleFunc("/editpdf-success", handleEditPdfSuccess)

	// hwp
	r.HandleFunc("/addhwp", handleAddHwp)
	r.HandleFunc("/addhwp-item", handleAddHwpItem)
	r.HandleFunc("/addhwp-file", handleAddHwpFile)
	r.HandleFunc("/uploadhwp-item", handleUploadHwpItem)
	r.HandleFunc("/uploadhwp-file", handleUploadHwpFile)
	r.HandleFunc("/uploadhwp-checkdata", handleUploadHwpCheckData)
	r.HandleFunc("/addhwp-success", handleAddHwpSuccess)
	r.HandleFunc("/edithwp", handleEditHwp)
	r.HandleFunc("/edithwp-submit", handleEditHwpSubmit)
	r.HandleFunc("/edithwp-success", handleEditHwpSuccess)

	// ppt
	r.HandleFunc("/addppt", handleAddPpt)
	r.HandleFunc("/addppt-item", handleAddPptItem)
	r.HandleFunc("/addppt-file", handleAddPptFile)
	r.HandleFunc("/uploadppt-item", handleUploadPptItem)
	r.HandleFunc("/uploadppt-file", handleUploadPptFile)
	r.HandleFunc("/uploadppt-checkdata", handleUploadPptCheckData)
	r.HandleFunc("/addppt-success", handleAddPptSuccess)
	r.HandleFunc("/editppt", handleEditPpt)
	r.HandleFunc("/editppt-submit", handleEditPptSubmit)
	r.HandleFunc("/editppt-success", handleEditPptSuccess)

	// unreal
	r.HandleFunc("/addunreal", handleAddUnreal)
	r.HandleFunc("/addunreal-item", handleAddUnrealItem)
	r.HandleFunc("/addunreal-file", handleAddUnrealFile)
	r.HandleFunc("/uploadunreal-item", handleUploadUnrealItem)
	r.HandleFunc("/uploadunreal-file", handleUploadUnrealFile)
	r.HandleFunc("/uploadunreal-checkdata", handleUploadUnrealCheckData)
	r.HandleFunc("/addunreal-success", handleAddUnrealSuccess)
	r.HandleFunc("/editunreal", handleEditUnreal)
	r.HandleFunc("/editunreal-submit", handleEditUnrealSubmit)
	r.HandleFunc("/editunreal-success", handleEditUnrealSuccess)

	// ies
	r.HandleFunc("/addies", handleAddIes)
	r.HandleFunc("/addies-item", handleAddIesItem)
	r.HandleFunc("/addies-file", handleAddIesFile)
	r.HandleFunc("/uploadies-item", handleUploadIesItem)
	r.HandleFunc("/uploadies-file", handleUploadIesFile)
	r.HandleFunc("/uploadies-checkdata", handleUploadIesCheckData)
	r.HandleFunc("/addies-success", handleAddIesSuccess)
	r.HandleFunc("/edities", handleEditIes)
	r.HandleFunc("/edities-submit", handleEditIesSubmit)
	r.HandleFunc("/edities-success", handleEditIesSuccess)

	// Modo
	r.HandleFunc("/addmodo", handleAddModo)
	r.HandleFunc("/addmodo-item", handleAddModoItem)
	r.HandleFunc("/addmodo-file", handleAddModoFile)
	r.HandleFunc("/uploadmodo-item", handleUploadModoItem)
	r.HandleFunc("/uploadmodo-file", handleUploadModoFile)
	r.HandleFunc("/uploadmodo-checkdata", handleUploadModoCheckData)
	r.HandleFunc("/addmodo-success", handleAddModoSuccess)
	r.HandleFunc("/editmodo", handleEditModo)
	r.HandleFunc("/editmodo-submit", handleEditModoSubmit)
	r.HandleFunc("/editmodo-success", handleEditModoSuccess)

	// Katana
	r.HandleFunc("/addkatana", handleAddKatana)
	r.HandleFunc("/addkatana-item", handleAddKatanaItem)
	r.HandleFunc("/addkatana-file", handleAddKatanaFile)
	r.HandleFunc("/uploadkatana-item", handleUploadKatanaItem)
	r.HandleFunc("/uploadkatana-file", handleUploadKatanaFile)
	r.HandleFunc("/uploadkatana-checkdata", handleUploadKatanaCheckData)
	r.HandleFunc("/addkatana-success", handleAddKatanaSuccess)
	r.HandleFunc("/editkatana", handleEditKatana)
	r.HandleFunc("/editkatana-submit", handleEditKatanaSubmit)
	r.HandleFunc("/editkatana-success", handleEditKatanaSuccess)

	// lut
	r.HandleFunc("/addlut", handleAddLut)
	r.HandleFunc("/addlut-item", handleAddLutItem)
	r.HandleFunc("/uploadlut-item", handleUploadLutItem)
	r.HandleFunc("/addlut-file", handleAddLutFile)
	r.HandleFunc("/uploadlut-file", handleUploadLutFile)
	r.HandleFunc("/uploadlut-checkdata", handleUploadLutCheckData)
	r.HandleFunc("/addlut-success", handleAddLutSuccess)
	r.HandleFunc("/editlut", handleEditLut)
	r.HandleFunc("/editlut-submit", handleEditLutSubmit)
	r.HandleFunc("/editlut-success", handleEditLutSuccess)

	// Matte
	r.HandleFunc("/addmatte", handleAddMatte)
	r.HandleFunc("/addmatte-item", handleAddMatteItem)
	r.HandleFunc("/uploadmatte-item", handleUploadMatteItem)
	r.HandleFunc("/addmatte-file", handleAddMatteFile)
	r.HandleFunc("/uploadmatte-file", handleUploadMatteFile)
	r.HandleFunc("/uploadmatte-checkdata", handleUploadMatteCheckData)
	r.HandleFunc("/addmatte-success", handleAddMatteSuccess)
	r.HandleFunc("/editmatte", handleEditMatte)
	r.HandleFunc("/editmatte-submit", handleEditMatteSubmit)
	r.HandleFunc("/editmatte-success", handleEditMatteSuccess)

	// Admin
	r.HandleFunc("/users", handleUsers)
	r.HandleFunc("/adminsetting", handleAdminSetting)
	r.HandleFunc("/adminsetting-submit", handleAdminSettingSubmit)
	r.HandleFunc("/adminsetting-success", handleAdminSettingSuccess)

	// Process
	r.HandleFunc("/item-process", handleItemProcess)
	r.HandleFunc("/cleanup-db", handleCleanUpDB)
	r.HandleFunc("/cleanup-db-submit", handleCleanUpDBSubmit)
	r.HandleFunc("/cleanup-db-success", handleCleanUpDBSuccess)

	// Download
	r.HandleFunc("/download-item", handleDownloadItem)

	// Help
	r.HandleFunc("/help", handleHelp)

	// Error
	r.HandleFunc("/error-ocio", handleErrorOCIO)

	// User
	r.HandleFunc("/profile", handleProfile)
	r.HandleFunc("/favoriteassets", handleFavoriteAssets)
	r.HandleFunc("/signup", handleSignup)
	r.HandleFunc("/signup-submit", handleSignupSubmit)
	r.HandleFunc("/signup-success", handleSignupSuccess)
	r.HandleFunc("/signin", handleSignin)
	r.HandleFunc("/signin-submit", handleSigninSubmit)
	r.HandleFunc("/signout", handleSignOut)
	r.HandleFunc("/invalidaccess", handleInvalidAccess)

	// RestAPI
	r.HandleFunc("/api/item", handleAPIGetItem).Methods("GET")
	r.HandleFunc("/api/item/{id}", handleAPIPutItem).Methods("PUT")
	r.HandleFunc("/api/item", handleAPIPostItem).Methods("POST")
	r.HandleFunc("/api/item", handleAPIDeleteItem).Methods("DELETE")
	r.HandleFunc("/api/search", handleAPISearch)
	r.HandleFunc("/api/adminsetting", handleAPIAdminSetting)
	r.HandleFunc("/api/dbbackup", postDBBackupHandler).Methods("POST")
	r.HandleFunc("/api/usingrate", handleAPIUsingRate)
	r.HandleFunc("/api/recentitem", handleAPIRecentItem)
	r.HandleFunc("/api/topusingitem", handleAPITopUsingItem)
	r.HandleFunc("/api/favoriteasset", handleAPIFavoriteAsset)
	r.HandleFunc("/api/initpassword", handleAPIInitPassword).Methods("POST")
	r.HandleFunc("/api/downloadzipfile", handleAPIDownloadZipfile)
	r.HandleFunc("/api/user/autoplay", handleAPIUserAutoplay).Methods("PUT")
	r.HandleFunc("/api/user/newsnum", handleAPIUserNewsNum).Methods("PUT")
	r.HandleFunc("/api/user/topnum", handleAPIUserTopNum).Methods("PUT")
	r.HandleFunc("/api/user/accesslevel", handleAPIUserAccessLevel).Methods("POST")

	// RestAPI Tags for an item
	r.HandleFunc("/api/tags", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/tags/{id}", getTagsHandler).Methods("GET")
	r.HandleFunc("/api/tags/{id}", putTagsHandler).Methods("PUT")

	// REST API Tag list for all item
	r.HandleFunc("/api/taglist", helpMethodOptionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/api/taglist", getTaglistHandler).Methods("GET")

	// REST API Category managing
	r.HandleFunc("/category", handleCategory)
	r.HandleFunc("/api/category", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/category", postCategoryHandler).Methods("POST")
	r.HandleFunc("/api/category/{id}", getCategoryHandler).Methods("GET")
	r.HandleFunc("/api/category/{id}", putCategoryHandler).Methods("PUT")
	r.HandleFunc("/api/category/{id}", deleteCategoryHandler).Methods("DELETE")
	r.HandleFunc("/api/rootcategories", getRootCategoriesHandler).Methods("GET")          // 메인 카테고리를 가지고 온다.
	r.HandleFunc("/api/subcategories/{parentid}", getSubCategoriesHandler).Methods("GET") // Sub 카테고리를 가지고 온다.

	r.Use(mux.CORSMethodMiddleware(r))
	http.Handle("/", r)
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
