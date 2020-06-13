package main

import (
	"context"
	"net/http"
	"strconv"
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
	type recipe struct {
		Adminsetting
		Token Token
	}
	rcp := recipe{}
	setting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = setting
	rcp.Token = token
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
	a.Umask = r.FormValue("umask")
	a.FolderPermission = r.FormValue("folderpermission")
	a.FilePermission = r.FormValue("filepermission")
	a.UID = r.FormValue("uid")
	a.GID = r.FormValue("gid")
	a.FFmpeg = r.FormValue("ffmpeg")
	a.OCIOConfig = r.FormValue("ocioconfig")
	a.OpenImageIO = r.FormValue("openimageio")
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
	containerWidth, err := strconv.Atoi(r.FormValue("containerwidth"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.ContainerWidth = containerWidth
	containerHeight, err := strconv.Atoi(r.FormValue("containerheight"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.ContainerHeight = containerHeight
	err = a.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.VideoCodecOgg = r.FormValue("videocodecogg")
	a.VideoCodecMp4 = r.FormValue("videocodecmp4")
	a.VideoCodecMov = r.FormValue("videocodecmov")
	a.AudioCodec = r.FormValue("audiocodec")
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
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "adminsetting-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
