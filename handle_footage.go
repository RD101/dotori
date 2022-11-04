package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/sys/unix"
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
		Adminsetting   Adminsetting
		Colorspaces    []Colorspace
		User           User
		RootCategories []Category
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
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	rcp.RootCategories, err = GetRootCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ocioConfig, err := loadOCIOConfig(rcp.Adminsetting.OCIOConfig)
	if err != nil {
		http.Redirect(w, r, "/error-ocio", http.StatusSeeOther)
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

func handleAddFootageItems(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting   Adminsetting
		RootCategories []Category
		Colorspaces    []Colorspace
		User           User
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
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	categories, err := GetRootCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.RootCategories = categories
	ocioConfig, err := loadOCIOConfig(rcp.Adminsetting.OCIOConfig)
	if err != nil {
		http.Redirect(w, r, "/error-ocio", http.StatusSeeOther)
		return
	}
	rcp.Colorspaces = ocioConfig.Colorspaces
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addfootage-items", rcp)
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
	item.Title = r.FormValue("title")
	item.Author = r.FormValue("author")
	item.Description = r.FormValue("description")
	item.InColorspace = r.FormValue("incolorspace")
	item.OutColorspace = r.FormValue("outcolorspace")
	item.Fps = r.FormValue("fps")
	tags := Str2List(r.FormValue("tags"))
	item.Tags = tags
	item.ItemType = "footage"
	if r.FormValue("premult") == "" {
		item.Premultiply = false
	} else {
		item.Premultiply = true
	}
	attr := make(map[string]string)
	attrNum, err := strconv.Atoi(r.FormValue("attributesNum"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < attrNum; i++ {
		key := r.FormValue(fmt.Sprintf("key%d", i))
		value := r.FormValue(fmt.Sprintf("value%d", i))
		if key == "" || value == "" {
			continue
		}
		attr[key] = value
	}
	item.Attributes = attr
	item.Status = "ready"
	item.Logs = append(item.Logs, "아이템이 생성되었습니다.")
	item.ThumbImgUploaded = false
	item.ThumbClipUploaded = false
	item.DataUploaded = false
	item.Categories = []string{r.FormValue("rootcategory-string"), r.FormValue("subcategory-string")}

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
	// admin setting에서 rootpath를 가져와서 경로를 생성한다.
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
	item.OutputProxyImgPath = rootpath + objIDpath + "/proxy/"
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

// handleAddFootageFile 함수는 Footage 파일을 추가하는 페이지 이다.
func handleAddFootageFile(w http.ResponseWriter, r *http.Request) {
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
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addfootage-file", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadFootageFile 함수는 Footage 파일을 DB에 업로드하는 페이지를 연다. dropzone에 파일을 올릴 경우 실행된다.
func handleUploadFootageFile(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	objectID, err := GetObjectIDfromRequestHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uploadFootageFile(w, r, objectID)
}

// uploadFootageFile 함수는 Footage 파일 정보를 DB에 업로드하고 파일을 storage에 복사한다.
func uploadFootageFile(w http.ResponseWriter, r *http.Request, objectID string) {
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
	item, err := GetItem(client, objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//admin setting에서 폴더권한에 관련된 옵션값을 가져온다
	um := adminsetting.Umask
	umask, err := strconv.Atoi(um)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	folderP := adminsetting.FolderPermission
	folderPerm, err := strconv.ParseInt(folderP, 8, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileP := adminsetting.FilePermission
	filePerm, err := strconv.ParseInt(fileP, 8, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u := adminsetting.UID
	uid, err := strconv.Atoi(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	g := adminsetting.GID
	gid, err := strconv.Atoi(g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer := adminsetting.MultipartFormBufferSize
	err = r.ParseMultipartForm(int64(buffer)) // grab the multipart form
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, files := range r.MultipartForm.File {
		for _, f := range files {
			if f.Size == 0 {
				http.Error(w, "file size is 0 byte", http.StatusInternalServerError)
				return
			}
			file, err := f.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				continue
			}
			defer file.Close()
			unix.Umask(umask)
			mimeType := f.Header.Get("Content-Type")
			switch mimeType {
			case "image/jpeg", "image/png", "image/x-exr", "application/octet-stream":
				ext := strings.ToLower(filepath.Ext(f.Filename))
				if ext != ".jpg" && ext != ".png" && ext != ".dpx" && ext != ".exr" { // .jpg .png .dpx .exr 외에는 허용하지 않는다.
					http.Error(w, "Not support file format", http.StatusBadRequest)
					return
				}
				data, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				path := item.OutputDataPath
				err = os.MkdirAll(path, os.FileMode(folderPerm))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = os.Chown(path, uid, gid)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = ioutil.WriteFile(path+"/"+f.Filename, data, os.FileMode(filePerm))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				item.DataUploaded = true
			default:
				//허용하지 않는 파일 포맷입니다.
				http.Error(w, "Not support file format", http.StatusBadRequest)
				return
			}
		}
	}
	if item.DataUploaded {
		item.Status = "fileuploaded"
	}
	err = SetItem(client, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadFootageCheckData 함수는 필요한 파일들을 모두 업로드했는지 체크하고, /addfootage-success 페이지로 redirect한다.
func handleUploadFootageCheckData(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting Adminsetting
		Item         Item
	}
	rcp := recipe{}
	rcp.Token = token
	// objectID로 item을 가져온다.
	objectID, err := GetObjectIDfromRequestHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	//rcp에 adminsetting 추가
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting
	//rcp에 item 추가
	item, err := GetItem(client, objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Item = item
	if !item.DataUploaded {
		err = TEMPLATES.ExecuteTemplate(w, "checkfootage-file", rcp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Redirect(w, r, "/addfootage-success", http.StatusSeeOther)
}

func handleAddFootageSuccess(w http.ResponseWriter, r *http.Request) {
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
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addfootage-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEditFootage(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		ID            primitive.ObjectID `json:"id" bson:"id"`
		ItemType      string             `json:"itemtype" bson:"itemtype"`
		Author        string             `json:"author" bson:"author"`
		Title         string             `json:"title" bson:"title"`
		Description   string             `json:"description" bson:"description"`
		Tags          []string           `json:"tags" bson:"tags"`
		Attributes    map[string]string  `json:"attributes" bson:"attributes"`
		InColorspace  string             `json:"incolorspace"`
		OutColorspace string             `json:"outcolorspace"`
		Fps           string             `json:"fps"`
		Colorspaces   []Colorspace       `json:"colorspaces"`
		Token
		Adminsetting Adminsetting
		User         User
	}
	q := r.URL.Query()
	id := q.Get("id")
	if id == "" {
		http.Error(w, "URL에 id를 입력해주세요", http.StatusBadRequest)
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

	item, err := SearchItem(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ocioConfig, err := loadOCIOConfig(adminsetting.OCIOConfig)
	if err != nil {
		http.Redirect(w, r, "/error-ocio", http.StatusSeeOther)
		return
	}
	rcp := recipe{
		ID:            item.ID,
		ItemType:      item.ItemType,
		Author:        item.Author,
		Title:         item.Title,
		Description:   item.Description,
		Tags:          item.Tags,
		Attributes:    item.Attributes,
		InColorspace:  item.InColorspace,
		OutColorspace: item.OutColorspace,
		Colorspaces:   ocioConfig.Colorspaces,
		Token:         token,
		Adminsetting:  adminsetting,
		Fps:           item.Fps,
	}
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	err = TEMPLATES.ExecuteTemplate(w, "editfootage", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//handleEditFootageSubmit 함수는 footage아이템을 수정하는 페이지에서 UPDATE버튼을 누르면 작동하는 함수다.
func handleEditFootageSubmit(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	id := r.FormValue("id")
	attrNum, err := strconv.Atoi(r.FormValue("attributesNum"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	attr := make(map[string]string)
	for i := 0; i < attrNum; i++ {
		key := r.FormValue(fmt.Sprintf("key%d", i))
		value := r.FormValue(fmt.Sprintf("value%d", i))
		if key == "" || value == "" {
			continue
		}
		attr[key] = value
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
	item, err := SearchItem(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.Author = r.FormValue("author")
	item.Title = r.FormValue("title")
	item.Description = r.FormValue("description")
	item.Tags = Str2List(r.FormValue("tags"))
	item.InColorspace = r.FormValue("incolorspace")
	item.OutColorspace = r.FormValue("outcolorspace")
	item.Attributes = attr
	err = item.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = SetItem(client, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/editfootage-success", http.StatusSeeOther)
}

func handleEditFootageSuccess(w http.ResponseWriter, r *http.Request) {
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
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "editfootage-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAPISearchFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Only POST", http.StatusMethodNotAllowed)
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
	user, err := GetUserFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !(user.AccessLevel == "admin" || user.AccessLevel == "default") {
		http.Error(w, "Need permission", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "path 를 설정해주세요", http.StatusBadRequest)
		return
	}
	// 경로에서 파일을 불러온다.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Rename struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	}
	rcp := []Rename{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileInfo := Rename{}
		fileInfo.Src = file.Name()
		fileInfo.Dst = file.Name()
		rcp = append(rcp, fileInfo)
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPISearchFootages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Only POST", http.StatusMethodNotAllowed)
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
	user, err := GetUserFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !(user.AccessLevel == "admin" || user.AccessLevel == "default") {
		http.Error(w, "Need permission", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "path 를 설정해주세요", http.StatusBadRequest)
		return
	}
	items, err := searchSeq(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIAddFootage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Only POST", http.StatusMethodNotAllowed)
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
	// 권한체크
	user, err := GetUserFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !(user.AccessLevel == "admin" || user.AccessLevel == "default") {
		http.Error(w, "Need permission", http.StatusUnauthorized)
		return
	}
	// 옵션 파싱
	r.ParseForm()
	item := Item{}
	item.ID = primitive.NewObjectID()
	item.InputData.Base = r.FormValue("base")
	item.InputData.Dir = r.FormValue("dir")
	frameIn, err := strconv.Atoi(r.FormValue("framein"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.InputData.FrameIn = frameIn
	frameOut, err := strconv.Atoi(r.FormValue("frameout"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.InputData.FrameOut = frameOut
	item.RequireMkdirInProcess = true
	item.RequireCopyInProcess = true
	item.DataUploaded = true
	item.Title = r.FormValue("title")
	item.Author = r.FormValue("author")
	item.Description = r.FormValue("description")
	item.InColorspace = r.FormValue("incolorspace")
	item.OutColorspace = r.FormValue("outcolorspace")
	item.Fps = r.FormValue("fps")
	item.Tags = Str2List(r.FormValue("tags"))
	item.Categories = Str2List(r.FormValue("categories"))
	item.ItemType = "footage"
	if r.FormValue("premult") == "" {
		item.Premultiply = false
	} else {
		item.Premultiply = true
	}
	// set Attribute
	item.Attributes = make(map[string]string) // init
	jsonAttributes := JsonAttributes{}
	item.Attributes, err = jsonAttributes.ToAttributes(r.FormValue("attributes"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.Status = "fileuploaded" // <- 이 부분은 모순이 있지만 input 소스를 이용해서 경로가 입력된 상황이다.
	item.Logs = append(item.Logs, "아이템이 생성되었습니다.")
	item.ThumbImgUploaded = false
	item.ThumbClipUploaded = false

	objIDpath, err := idToPath(item.ID.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rootpath, err := GetRootPath(client)
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
	item.OutputProxyImgPath = rootpath + objIDpath + "/proxy/"

	// 에러체크
	err = item.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 연산 목록에 등록
	err = AddItem(client, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json 반환
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
