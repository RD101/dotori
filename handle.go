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
		log.Println(*flagMongoDBURI + " 에 mongoDB가 실행되고 있는지 체크해주세요")
		log.Fatal(err)
	}
	// 웹주소 설정
	http.HandleFunc("/", handleInit)
	http.HandleFunc("/mediadata", handleMediaData)
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

	// 3dsmax
	http.HandleFunc("/addmax", handleAddMax)
	http.HandleFunc("/addmax-item", handleAddMaxItem)
	http.HandleFunc("/addmax-file", handleAddMaxFile)
	http.HandleFunc("/uploadmax-item", handleUploadMaxItem)
	http.HandleFunc("/uploadmax-file", handleUploadMaxFile)
	http.HandleFunc("/uploadmax-checkdata", handleUploadMaxCheckData)
	http.HandleFunc("/addmax-success", handleAddMaxSuccess)
	http.HandleFunc("/editmax", handleEditMax)
	http.HandleFunc("/editmax-submit", handleEditMaxSubmit)
	http.HandleFunc("/editmax-success", handleEditMaxSuccess)

	// Fusion360
	http.HandleFunc("/addfusion360", handleAddFusion360)
	http.HandleFunc("/addfusion360-item", handleAddFusion360Item)
	http.HandleFunc("/addfusion360-file", handleAddFusion360File)
	http.HandleFunc("/uploadfusion360-item", handleUploadFusion360Item)
	http.HandleFunc("/uploadfusion360-file", handleUploadFusion360File)
	http.HandleFunc("/uploadfusion360-checkdata", handleUploadFusion360CheckData)
	http.HandleFunc("/addfusion360-success", handleAddFusion360Success)
	http.HandleFunc("/editfusion360", handleEditFusion360)
	http.HandleFunc("/editfusion360-submit", handleEditFusion360Submit)
	http.HandleFunc("/editfusion360-success", handleEditFusion360Success)

	// Footage
	http.HandleFunc("/addfootage", handleAddFootage)
	http.HandleFunc("/addfootage-item", handleAddFootageItem)
	http.HandleFunc("/addfootage-items", handleAddFootageItems)
	http.HandleFunc("/uploadfootage-item", handleUploadFootageItem)
	http.HandleFunc("/addfootage-file", handleAddFootageFile)
	http.HandleFunc("/uploadfootage-file", handleUploadFootageFile)
	http.HandleFunc("/uploadfootage-checkdata", handleUploadFootageCheckData)
	http.HandleFunc("/addfootage-success", handleAddFootageSuccess)
	http.HandleFunc("/editfootage", handleEditFootage)
	http.HandleFunc("/editfootage-submit", handleEditFootageSubmit)
	http.HandleFunc("/editfootage-success", handleEditFootageSuccess)
	http.HandleFunc("/api/searchfootages", handleAPISearchFootageAndClip)
	http.HandleFunc("/api/addfootage", handleAPIAddFootage)

	// Clip
	http.HandleFunc("/addclip", handleAddClip)
	http.HandleFunc("/addclip-item", handleAddClipItem)
	http.HandleFunc("/addclip-items", handleAddClipItems)
	http.HandleFunc("/uploadclip-item", handleUploadClipItem)
	http.HandleFunc("/addclip-file", handleAddClipFile)
	http.HandleFunc("/uploadclip-file", handleUploadClipFile)
	http.HandleFunc("/uploadclip-checkdata", handleUploadClipCheckData)
	http.HandleFunc("/addclip-success", handleAddClipSuccess)
	http.HandleFunc("/editclip", handleEditClip)
	http.HandleFunc("/editclip-submit", handleEditClipSubmit)
	http.HandleFunc("/editclip-success", handleEditClipSuccess)
	http.HandleFunc("/api/searchclips", handleAPISearchClips)
	http.HandleFunc("/api/addclip", handleAPIAddClip)

	// HDRI
	http.HandleFunc("/addhdri", handleAddHDRI)
	http.HandleFunc("/addhdri-item", handleAddHDRIItem)
	http.HandleFunc("/uploadhdri-item", handleUploadHDRIItem)
	http.HandleFunc("/addhdri-file", handleAddHDRIFile)
	http.HandleFunc("/uploadhdri-file", handleUploadHDRIFile)
	http.HandleFunc("/uploadhdri-checkdata", handleUploadHDRICheckData)
	http.HandleFunc("/addhdri-success", handleAddHDRISuccess)
	http.HandleFunc("/edithdri", handleEditHDRI)
	http.HandleFunc("/edithdri-submit", handleEditHDRISubmit)
	http.HandleFunc("/edithdri-success", handleEditHDRISuccess)

	// Texture
	http.HandleFunc("/addtexture", handleAddTexture)
	http.HandleFunc("/addtexture-item", handleAddTextureItem)
	http.HandleFunc("/uploadtexture-item", handleUploadTextureItem)
	http.HandleFunc("/addtexture-file", handleAddTextureFile)
	http.HandleFunc("/uploadtexture-file", handleUploadTextureFile)
	http.HandleFunc("/uploadtexture-checkdata", handleUploadTextureCheckData)
	http.HandleFunc("/addtexture-success", handleAddTextureSuccess)
	http.HandleFunc("/edittexture", handleEditTexture)
	http.HandleFunc("/edittexture-submit", handleEditTextureSubmit)
	http.HandleFunc("/edittexture-success", handleEditTextureSuccess)

	// Nuke
	http.HandleFunc("/addnuke", handleAddNuke)
	http.HandleFunc("/addnuke-item", handleAddNukeItem)
	http.HandleFunc("/uploadnuke-item", handleUploadNukeItem)
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
	http.HandleFunc("/addhoudini-success", handleAddHoudiniSuccess)
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

	// OpenVDB
	http.HandleFunc("/addopenvdb", handleAddOpenVDB)
	http.HandleFunc("/addopenvdb-item", handleAddOpenVDBItem)
	http.HandleFunc("/uploadopenvdb-item", handleUploadOpenVDBItem)
	http.HandleFunc("/addopenvdb-file", handleAddOpenVDBFile)
	http.HandleFunc("/uploadopenvdb-file", handleUploadOpenVDBFile)
	http.HandleFunc("/uploadopenvdb-checkdata", handleUploadOpenVDBCheckData)
	http.HandleFunc("/addopenvdb-success", handleAddOpenVDBSuccess)
	http.HandleFunc("/editopenvdb", handleEditOpenVDB)
	http.HandleFunc("/editopenvdb-submit", handleEditOpenVDBSubmit)
	http.HandleFunc("/editopenvdb-success", handleEditOpenVDBSuccess)

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

	// ppt
	http.HandleFunc("/addppt", handleAddPpt)
	http.HandleFunc("/addppt-item", handleAddPptItem)
	http.HandleFunc("/addppt-file", handleAddPptFile)
	http.HandleFunc("/uploadppt-item", handleUploadPptItem)
	http.HandleFunc("/uploadppt-file", handleUploadPptFile)
	http.HandleFunc("/uploadppt-checkdata", handleUploadPptCheckData)
	http.HandleFunc("/addppt-success", handleAddPptSuccess)
	http.HandleFunc("/editppt", handleEditPpt)
	http.HandleFunc("/editppt-submit", handleEditPptSubmit)
	http.HandleFunc("/editppt-success", handleEditPptSuccess)

	// unreal
	http.HandleFunc("/addunreal", handleAddUnreal)
	http.HandleFunc("/addunreal-item", handleAddUnrealItem)
	http.HandleFunc("/addunreal-file", handleAddUnrealFile)
	http.HandleFunc("/uploadunreal-item", handleUploadUnrealItem)
	http.HandleFunc("/uploadunreal-file", handleUploadUnrealFile)
	http.HandleFunc("/uploadunreal-checkdata", handleUploadUnrealCheckData)
	http.HandleFunc("/addunreal-success", handleAddUnrealSuccess)
	http.HandleFunc("/editunreal", handleEditUnreal)
	http.HandleFunc("/editunreal-submit", handleEditUnrealSubmit)
	http.HandleFunc("/editunreal-success", handleEditUnrealSuccess)

	// ies
	http.HandleFunc("/addies", handleAddIes)
	http.HandleFunc("/addies-item", handleAddIesItem)
	http.HandleFunc("/addies-file", handleAddIesFile)
	http.HandleFunc("/uploadies-item", handleUploadIesItem)
	http.HandleFunc("/uploadies-file", handleUploadIesFile)
	http.HandleFunc("/uploadies-checkdata", handleUploadIesCheckData)
	http.HandleFunc("/addies-success", handleAddIesSuccess)
	http.HandleFunc("/edities", handleEditIes)
	http.HandleFunc("/edities-submit", handleEditIesSubmit)
	http.HandleFunc("/edities-success", handleEditIesSuccess)

	// Modo
	http.HandleFunc("/addmodo", handleAddModo)
	http.HandleFunc("/addmodo-item", handleAddModoItem)
	http.HandleFunc("/addmodo-file", handleAddModoFile)
	http.HandleFunc("/uploadmodo-item", handleUploadModoItem)
	http.HandleFunc("/uploadmodo-file", handleUploadModoFile)
	http.HandleFunc("/uploadmodo-checkdata", handleUploadModoCheckData)
	http.HandleFunc("/addmodo-success", handleAddModoSuccess)
	http.HandleFunc("/editmodo", handleEditModo)
	http.HandleFunc("/editmodo-submit", handleEditModoSubmit)
	http.HandleFunc("/editmodo-success", handleEditModoSuccess)

	// Katana
	http.HandleFunc("/addkatana", handleAddKatana)
	http.HandleFunc("/addkatana-item", handleAddKatanaItem)
	http.HandleFunc("/addkatana-file", handleAddKatanaFile)
	http.HandleFunc("/uploadkatana-item", handleUploadKatanaItem)
	http.HandleFunc("/uploadkatana-file", handleUploadKatanaFile)
	http.HandleFunc("/uploadkatana-checkdata", handleUploadKatanaCheckData)
	http.HandleFunc("/addkatana-success", handleAddKatanaSuccess)
	http.HandleFunc("/editkatana", handleEditKatana)
	http.HandleFunc("/editkatana-submit", handleEditKatanaSubmit)
	http.HandleFunc("/editkatana-success", handleEditKatanaSuccess)

	// lut
	http.HandleFunc("/addlut", handleAddLut)
	http.HandleFunc("/addlut-item", handleAddLutItem)
	http.HandleFunc("/uploadlut-item", handleUploadLutItem)
	http.HandleFunc("/addlut-file", handleAddLutFile)
	http.HandleFunc("/uploadlut-file", handleUploadLutFile)
	http.HandleFunc("/uploadlut-checkdata", handleUploadLutCheckData)
	http.HandleFunc("/addlut-success", handleAddLutSuccess)
	http.HandleFunc("/editlut", handleEditLut)
	http.HandleFunc("/editlut-submit", handleEditLutSubmit)
	http.HandleFunc("/editlut-success", handleEditLutSuccess)

	// Matte
	http.HandleFunc("/addmatte", handleAddMatte)
	http.HandleFunc("/addmatte-item", handleAddMatteItem)
	http.HandleFunc("/uploadmatte-item", handleUploadMatteItem)
	http.HandleFunc("/addmatte-file", handleAddMatteFile)
	http.HandleFunc("/uploadmatte-file", handleUploadMatteFile)
	http.HandleFunc("/uploadmatte-checkdata", handleUploadMatteCheckData)
	http.HandleFunc("/addmatte-success", handleAddMatteSuccess)
	http.HandleFunc("/editmatte", handleEditMatte)
	http.HandleFunc("/editmatte-submit", handleEditMatteSubmit)
	http.HandleFunc("/editmatte-success", handleEditMatteSuccess)

	// Admin
	http.HandleFunc("/users", handleUsers)
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
	http.HandleFunc("/favoriteassets", handleFavoriteAssets)
	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/signup-submit", handleSignupSubmit)
	http.HandleFunc("/signup-success", handleSignupSuccess)
	http.HandleFunc("/signin", handleSignin)
	http.HandleFunc("/signin-submit", handleSigninSubmit)
	http.HandleFunc("/signout", handleSignOut)
	http.HandleFunc("/invalidaccess", handleInvalidAccess)

	// REST API
	http.HandleFunc("/api/item", handleAPIItem)
	http.HandleFunc("/api/search", handleAPISearch)
	http.HandleFunc("/api/adminsetting", handleAPIAdminSetting)
	http.HandleFunc("/api/usingrate", handleAPIUsingRate)
	http.HandleFunc("/api/recentitem", handleAPIRecentItem)
	http.HandleFunc("/api/topusingitem", handleAPITopUsingItem)
	http.HandleFunc("/api/favoriteasset", handleAPIFavoriteAsset)
	http.HandleFunc("/api/initpassword", handleAPIInitPassword)
	http.HandleFunc("/api/downloadzipfile", handleAPIDownloadZipfile)
	http.HandleFunc("/api/user", handleAPIUser)

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
