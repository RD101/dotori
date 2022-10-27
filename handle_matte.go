package main

import (
	"context"
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

// handleAddMatte 함수는 URL에 objectID를 붙여서 /addmatte-item 페이지로 redirect한다.
func handleAddMatte(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	objectID := primitive.NewObjectID().Hex()
	http.Redirect(w, r, fmt.Sprintf("/addmatte-item?objectid=%s", objectID), http.StatusSeeOther)
}

func handleAddMatteItem(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		Adminsetting   Adminsetting
		Colorspaces    []Colorspace
		RootCategories []Category
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
	rcp.Adminsetting, err = GetAdminSetting(client)
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
	err = TEMPLATES.ExecuteTemplate(w, "addmatte-item", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUploadMatteItem(w http.ResponseWriter, r *http.Request) {
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
	item.Tags = Str2List(r.FormValue("tags"))
	item.Categories = []string{r.FormValue("rootcategory-string"), r.FormValue("subcategory-string")}
	item.ItemType = "matte"
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
	http.Redirect(w, r, fmt.Sprintf("/addmatte-file?objectid=%s", objectID), http.StatusSeeOther)
}

// handleAddMatteFile 함수는 Matte 소스 파일을 추가하는 페이지 이다.
func handleAddMatteFile(w http.ResponseWriter, r *http.Request) {
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
	rcp.Adminsetting, err = GetAdminSetting(client)
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
	err = TEMPLATES.ExecuteTemplate(w, "addmatte-file", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadMatteFile 함수는 Matte 소스 파일을 DB에 업로드하는 페이지를 연다. dropzone에 파일을 올릴 경우 실행된다.
func handleUploadMatteFile(w http.ResponseWriter, r *http.Request) {
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
	uploadMatteFile(w, r, objectID)
}

// uploadMatteFile 함수는 Matte 소스 파일 정보를 DB에 업로드하고 파일을 storage에 복사한다.
func uploadMatteFile(w http.ResponseWriter, r *http.Request, objectID string) {
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
			case "application/octet-stream", "image/jpeg", "image/png", "image/tiff", "image/x-tiff", "image/x-tga":
				ext := strings.ToLower(filepath.Ext(f.Filename))
				if !(ext == ".exr" || ext == ".png" || ext == ".jpg" || ext == ".tga" || ext == ".tif" || ext == ".tiff" || ext == ".hdr" || ext == ".hdri") {
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
			tags := FilenameToTags(f.Filename)
			for _, tag := range tags {
				has := false // 중복되는 tag가 있다면 append하지 않는다.
				for _, t := range item.Tags {
					if tag == t {
						has = true
					}
				}
				if !has {
					item.Tags = append(item.Tags, tag)
				}
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

func handleUploadMatteCheckData(w http.ResponseWriter, r *http.Request) {
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
	rcp.Adminsetting, err = GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		err = TEMPLATES.ExecuteTemplate(w, "checkmatte-file", rcp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Redirect(w, r, "/addmatte-success", http.StatusSeeOther)
}

func handleAddMatteSuccess(w http.ResponseWriter, r *http.Request) {
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
	rcp.Adminsetting, err = GetAdminSetting(client)
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
	err = TEMPLATES.ExecuteTemplate(w, "addmatte-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEditMatte(w http.ResponseWriter, r *http.Request) {
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
	}
	rcp.User, err = GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = TEMPLATES.ExecuteTemplate(w, "editmatte", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//handleEditMatteSubmit 함수는 Matte 아이템을 수정하는 페이지에서 UPDATE 버튼을 누르면 작동하는 함수다.
func handleEditMatteSubmit(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/editmatte-success", http.StatusSeeOther)
}

func handleEditMatteSuccess(w http.ResponseWriter, r *http.Request) {
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
	rcp.Adminsetting, err = GetAdminSetting(client)
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
	err = TEMPLATES.ExecuteTemplate(w, "editmatte-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
