package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleAdminSetting 함수는 Admin 설정 페이지로 이동한다.
func handleAdminSetting(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Access level 체크
	if token.AccessLevel != "admin" {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Adminsetting
		Token Token
		User  User
	}
	rcp := recipe{}
	setting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = setting
	rcp.Token = token
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "adminsetting", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAdminSettingSubmit 함수는 관리자 설정을 저장한다.
func handleAdminSettingSubmit(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	a := Adminsetting{}
	a.Appname = r.FormValue("appname")
	a.Rootpath = r.FormValue("rootpath")
	a.LinuxProtocolPath = r.FormValue("linuxprotocolpath")
	a.WindowsProtocolPath = r.FormValue("windowsprotocolpath")
	a.MacosProtocolPath = r.FormValue("macosprotocolpath")
	a.WindowsUNCPrefix = r.FormValue("windowsuncprefix")
	a.Umask = r.FormValue("umask")
	a.FolderPermission = r.FormValue("folderpermission")
	a.FilePermission = r.FormValue("filepermission")
	a.UID = r.FormValue("uid")
	a.GID = r.FormValue("gid")
	pbsize, err := strconv.Atoi(r.FormValue("processbuffersize"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.ProcessBufferSize = pbsize
	a.FFmpeg = r.FormValue("ffmpeg")
	a.OCIOConfig = r.FormValue("ocioconfig")
	a.OpenImageIO = r.FormValue("openimageio")
	a.Mongodump = r.FormValue("mongodump")
	a.Backuppath = r.FormValue("backuppath")
	a.LDLibraryPath = r.FormValue("ldlibrarypath")
	bsize, err := strconv.Atoi(r.FormValue("multipartformbuffersize"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.MultipartFormBufferSize = bsize
	imageWidth, err := strconv.Atoi(r.FormValue("thumbnailimagewidth"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.ThumbnailImageWidth = imageWidth
	imageHeight, err := strconv.Atoi(r.FormValue("thumbnailimageheight"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.ThumbnailImageHeight = imageHeight
	mediaWidth, err := strconv.Atoi(r.FormValue("mediawidth"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaWidth < 1 {
		http.Error(w, "mediaWidth 값은 1보다 커야합니다", http.StatusBadRequest)
		return
	}
	a.MediaWidth = mediaWidth
	mediaHeight, err := strconv.Atoi(r.FormValue("mediaheight"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaHeight < 1 {
		http.Error(w, "mediaHeight 값은 1보다 커야합니다", http.StatusBadRequest)
		return
	}
	a.MediaHeight = mediaHeight
	err = a.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.VideoCodecOgg = r.FormValue("videocodecogg")
	a.VideoCodecMp4 = r.FormValue("videocodecmp4")
	a.VideoCodecMov = r.FormValue("videocodecmov")
	a.AudioCodec = r.FormValue("audiocodec")
	a.InitPassword = r.FormValue("initpassword")
	a.EnableRVLink = str2bool(r.FormValue("enablervlink"))
	a.EnableCategory = str2bool(r.FormValue("enablecategory"))
	// 지원 포멧 셋팅
	a.Maya = str2bool(r.FormValue("maya"))
	a.Max = str2bool(r.FormValue("max"))
	a.Houdini = str2bool(r.FormValue("houdini"))
	a.Blender = str2bool(r.FormValue("blender"))
	a.Alembic = str2bool(r.FormValue("alembic"))
	a.USD = str2bool(r.FormValue("usd"))
	a.Unreal = str2bool(r.FormValue("unreal"))
	a.OpenVDB = str2bool(r.FormValue("openvdb"))
	a.Modo = str2bool(r.FormValue("modo"))
	a.Fusion360 = str2bool(r.FormValue("fusion360"))
	a.Katana = str2bool(r.FormValue("katana"))
	a.HWP = str2bool(r.FormValue("hwp"))
	a.PDF = str2bool(r.FormValue("pdf"))
	a.PPT = str2bool(r.FormValue("ppt"))
	a.LUT = str2bool(r.FormValue("lut"))
	a.HDRI = str2bool(r.FormValue("hdri"))
	a.IES = str2bool(r.FormValue("ies"))
	a.Texture = str2bool(r.FormValue("texture"))
	a.Matte = str2bool(r.FormValue("matte"))
	a.Footage = str2bool(r.FormValue("footage"))
	a.MultipleFootage = str2bool(r.FormValue("multiplefootage"))
	a.Clip = str2bool(r.FormValue("clip"))
	a.MultipleClip = str2bool(r.FormValue("multipleclip"))
	a.Nuke = str2bool(r.FormValue("nuke"))
	a.Sound = str2bool(r.FormValue("sound"))
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = SetAdminSetting(client, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/adminsetting-success", http.StatusSeeOther)
}

func handleAdminSettingSuccess(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting Adminsetting
		User         User
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
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
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
	rcp.User, err = GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "adminsetting-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func postDBBackupHandler(w http.ResponseWriter, r *http.Request) {
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accesslevel, err := GetAccessLevelFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accesslevel != "default" && accesslevel != "manager" && accesslevel != "admin" {
		http.Error(w, "등록 권한이 없는 계정입니다", http.StatusUnauthorized)
		return
	}
	type Option struct {
		Date       string `json:"date"`
		Mongodump  string `json:"mongodump"`
		Backuppath string `json:"backuppath"`
	}
	opt := Option{}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&opt)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// adminsetting을 불러오기
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := os.Stat(adminsetting.Mongodump); errors.Is(err, os.ErrNotExist) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err := os.Stat(adminsetting.Backuppath); errors.Is(err, os.ErrNotExist) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	opt.Mongodump = adminsetting.Mongodump
	opt.Backuppath = adminsetting.Backuppath
	dbHostName := strings.ReplaceAll(*flagMongoDBURI, "mongodb://", "")
	args := []string{
		"-h",
		dbHostName,
		"-o",
		adminsetting.Backuppath + "/" + opt.Date,
	}
	err = exec.Command(adminsetting.Mongodump, args...).Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data, err := json.Marshal(opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
