package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleAddFootage 함수는 URL에 objectID를 붙여서 /addfootage-item 페이지로 redirect한다.
func handleAddFootage(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	objectID := primitive.NewObjectID().Hex()
	http.Redirect(w, r, fmt.Sprintf("/addfootage-item?objectid=%s", objectID), http.StatusSeeOther)
}

func handleAddFootageItem(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting Adminsetting
		Colorspaces  []Colorspace
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
	ocioConfig, err := loadOCIOConfig(rcp.Adminsetting.OCIOConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Colorspaces = ocioConfig.Colorspaces
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addfootage-item", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadFootageItem 핸들러는 Footage 아이템을 생성한다.
func handleUploadFootageItem(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	item := Item{}
	objectID, err := GetObjectIDfromRequestHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.ID, err = primitive.ObjectIDFromHex(objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.Author = r.FormValue("author")
	item.Description = r.FormValue("description")
	item.InColorspace = r.FormValue("incolorspace")
	item.OutColorspace = r.FormValue("outcolorspace")
	tags := SplitBySpace(r.FormValue("tag"))
	item.Tags = tags
	item.ItemType = "footage"
	attr := make(map[string]string)
	attrNum, err := strconv.Atoi(r.FormValue("attributesNum"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < attrNum; i++ {
		key := r.FormValue(fmt.Sprintf("key%d", i))
		value := r.FormValue(fmt.Sprintf("value%d", i))
		attr[key] = value
	}
	item.Attributes = attr
	item.Status = "ready"
	item.Logs = append(item.Logs, "아이템이 생성되었습니다.")
	currentTime := time.Now()
	item.CreateTime = currentTime.Format("2006-01-02 15:04:05")
	item.ThumbImgUploaded = false
	item.ThumbClipUploaded = false
	item.DataUploaded = false

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
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	objIDpath, err := idToPath(item.ID.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.InputThumbnailImgPath = rootpath + objIDpath + "/originalthumbimg/"
	item.InputThumbnailClipPath = rootpath + objIDpath + "/originalthumbmov/"
	item.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	item.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	item.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	item.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	item.OutputDataPath = rootpath + objIDpath + "/data/"
	err = item.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = AddItem(client, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/addfootage-file?objectid=%s", objectID), http.StatusSeeOther)
}
