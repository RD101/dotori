package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2"
)

// handleAdminSetting 함수는 Admin 설정 페이지로 이동한다.
func handleAdminSetting(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
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
	bsize, err := strconv.Atoi(r.FormValue("multipartformbuffersize"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.MultipartFormBufferSize = bsize
	err = a.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	err = SetAdminSetting(session, a)
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
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "adminsetting-success", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
