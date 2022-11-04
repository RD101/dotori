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

// handleAddClip 함수는 URL에 objectID를 붙여서 /addclip-item 페이지로 redirect한다.
func handleAddClip(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	objectID := primitive.NewObjectID().Hex()
	http.Redirect(w, r, fmt.Sprintf("/addclip-item?objectid=%s", objectID), http.StatusSeeOther)
}

func handleAddClipItem(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting   Adminsetting
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
	rcp.User, err = GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.RootCategories, err = GetRootCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addclip-item", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddClipItems(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting   Adminsetting
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
	rcp.RootCategories, err = GetRootCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addclip-items", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadClip 핸들러는 clip 아이템을 생성한다.
func handleUploadClipItem(w http.ResponseWriter, r *http.Request) {
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
	item.Fps = r.FormValue("fps")
	tags := Str2List(r.FormValue("tags"))
	item.Tags = tags
	item.ItemType = "clip"
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
	http.Redirect(w, r, fmt.Sprintf("/addclip-file?objectid=%s", objectID), http.StatusSeeOther)
}

// handleAddClipFile 함수는 Clip 파일을 추가하는 페이지 이다.
func handleAddClipFile(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "addclip-file", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadClipFile 함수는 Clip 파일을 DB에 업로드하는 페이지를 연다. dropzone에 파일을 올릴 경우 실행된다.
func handleUploadClipFile(w http.ResponseWriter, r *http.Request) {
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
	uploadClipFile(w, r, objectID)
}

// uploadClipFile 함수는 Clip 파일 정보를 DB에 업로드하고 파일을 storage에 복사한다.
func uploadClipFile(w http.ResponseWriter, r *http.Request, objectID string) {
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
				http.Error(w, "파일사이즈가 0 바이트입니다", http.StatusInternalServerError)
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
			case "video/quicktime", "video/mp4":
				ext := strings.ToLower(filepath.Ext(f.Filename))
				if !(ext == ".mov" || ext == ".mp4") { // .mov, .mp4 외에는 허용하지 않는다.
					http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
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
				http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
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

// handleUploadClipCheckData 함수는 필요한 파일들을 모두 업로드했는지 체크하고, /addclip-success 페이지로 redirect한다.
func handleUploadClipCheckData(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting Adminsetting
		Item         Item
		User         User
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
	rcp.User, err = GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//rcp에 item 추가
	item, err := GetItem(client, objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Item = item
	if !item.DataUploaded {
		err = TEMPLATES.ExecuteTemplate(w, "checkclip-file", rcp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Redirect(w, r, "/addclip-success", http.StatusSeeOther)
}

func handleAddClipSuccess(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "addclip-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEditClip(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		ID          primitive.ObjectID `json:"id" bson:"id"`
		ItemType    string             `json:"itemtype" bson:"itemtype"`
		Author      string             `json:"author" bson:"author"`
		Title       string             `json:"title" bson:"title"`
		Description string             `json:"description" bson:"description"`
		Tags        []string           `json:"tags" bson:"tags"`
		Attributes  map[string]string  `json:"attributes" bson:"attributes"`
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

	rcp := recipe{
		ID:           item.ID,
		ItemType:     item.ItemType,
		Author:       item.Author,
		Title:        item.Title,
		Description:  item.Description,
		Tags:         item.Tags,
		Attributes:   item.Attributes,
		Token:        token,
		Adminsetting: adminsetting,
	}

	err = TEMPLATES.ExecuteTemplate(w, "editclip", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//handleEditClipSubmit 함수는 clip 아이템을 수정하는 페이지에서 UPDATE버튼을 누르면 작동하는 함수다.
func handleEditClipSubmit(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/editclip-success", http.StatusSeeOther)
}

func handleEditClipSuccess(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "editclip-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAPISearchClips(w http.ResponseWriter, r *http.Request) {
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
	clips, err := searchClip(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(clips)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIAddClip(w http.ResponseWriter, r *http.Request) {
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
	item.RequireMkdirInProcess = true
	item.RequireCopyInProcess = true
	item.DataUploaded = true
	item.Title = r.FormValue("title")
	item.Author = r.FormValue("author")
	item.Description = r.FormValue("description")
	item.Tags = Str2List(r.FormValue("tags"))
	item.Categories = Str2List(r.FormValue("categories"))
	item.ItemType = "clip"
	// set Attribute
	item.Attributes = make(map[string]string) // init
	jsonAttributes := JsonAttributes{}
	item.Attributes, err = jsonAttributes.ToAttributes(r.FormValue("attributes"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.Status = "fileuploaded" // <- 이미 input 소스를 이용해서 경로가 입력된 상황이다.
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
