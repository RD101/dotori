package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/sys/unix"
)

// handleAddMayaFile 함수는 Maya 파일을 추가하는 페이지 이다.
func handleAddMayaFile(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "addmaya-file", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddMayaItem(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "addmaya-item", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddMayaSubmit 함수는 URL에 objectID를 붙여서 /addmaya-item 페이지로 redirect한다.
func handleAddMaya(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	objectID := primitive.NewObjectID().Hex()
	http.Redirect(w, r, fmt.Sprintf("/addmaya-item?objectid=%s", objectID), http.StatusSeeOther)
}

func handleAddHoudini(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addhoudini", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddBlender(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addblender", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddAlembic(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addalembic", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddUSD(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addusd", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddHoudiniProcess(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addhoudini-process", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddAlembicProcess(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addalembic-process", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUploadMayaItem(w http.ResponseWriter, r *http.Request) {
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
	tags := SplitBySpace(r.FormValue("tag"))
	item.Tags = tags
	item.ItemType = "maya"
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
	http.Redirect(w, r, fmt.Sprintf("/addmaya-file?objectid=%s", objectID), http.StatusSeeOther)
}

// handleUploadMaya 함수는 Maya파일을 DB에 업로드하는 페이지를 연다. dropzone에 파일을 올릴 경우 실행된다.
func handleUploadMayaFile(w http.ResponseWriter, r *http.Request) {
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
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item, err := GetItem(client, "maya", objectID)
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
			case "image/jpeg", "image/png":
				data, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				path, filename := path.Split(item.InputThumbnailImgPath)
				if filename != "" { // 파일이 이미 존재하는 경우, 지우고 경로를 새로 지정한다.
					err = os.Remove(item.InputThumbnailImgPath)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					item.InputThumbnailImgPath = path
				}
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
				err = ioutil.WriteFile(path+f.Filename, data, os.FileMode(filePerm))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				item.InputThumbnailImgPath = path + f.Filename
				item.ThumbImgUploaded = true
			case "video/quicktime", "video/mp4", "video/ogg", "application/ogg":
				data, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				path, filename := path.Split(item.InputThumbnailClipPath)
				if filename != "" { // 파일이 이미 존재하는 경우, 지우고 경로를 새로 지정한다.
					err = os.Remove(item.InputThumbnailClipPath)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					item.InputThumbnailClipPath = path
				}
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
				err = ioutil.WriteFile(path+f.Filename, data, os.FileMode(filePerm))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				item.InputThumbnailClipPath = path + f.Filename
				item.ThumbClipUploaded = true
			case "application/octet-stream":
				ext := filepath.Ext(f.Filename)
				if ext != ".mb" && ext != ".ma" { // .ma .mb 외에는 허용하지 않는다.
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
	if item.ThumbImgUploaded && item.ThumbClipUploaded && item.DataUploaded {
		item.Status = "fileuploaded"
	}
	UpdateItem(client, "maya", item)
}

func handleUploadMayaCheckData(w http.ResponseWriter, r *http.Request) {
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
	item, err := GetItem(client, "maya", objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 썸네일이미지가 있는지 체크
	if !item.ThumbImgUploaded {
		http.Error(w, "썸네일이미지를 업로드해주세요", http.StatusBadRequest)
		return
	}
	// 썸네일클립이 있는지 체크
	if !item.ThumbClipUploaded {
		http.Error(w, "썸네일클립을 업로드해주세요", http.StatusBadRequest)
		return
	}
	// 마야파일이 있는지 체크
	if !item.DataUploaded {
		http.Error(w, "마야 파일을 업로드해주세요", http.StatusBadRequest)
		return
	}
	// addmaya-success 페이지로 연결
	http.Redirect(w, r, "/addmaya-success", http.StatusSeeOther)
}

func handleAddMayaSuccess(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "addmaya-success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddHoudiniProcess 함수는 Houdini 파일을 처리하는 페이지 이다.
func handleUploadHoudini(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	unix.Umask(0)
	mimeType := header.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg", "image/png":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		path := os.TempDir() + "/dotori/thumbnail"
		err = os.MkdirAll(path, 0770)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(path+"/"+header.Filename, data, 0440)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	case "video/quicktime", "video/mp4", "video/ogg", "application/ogg":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		path := os.TempDir() + "/dotori/preview"
		err = os.MkdirAll(path, 0770)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(path+"/"+header.Filename, data, 0440)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	case "application/octet-stream":
		ext := filepath.Ext(header.Filename)
		if ext == ".hip" || ext == ".hda" {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
				return
			}
			path := os.TempDir() + "/dotori"
			err = os.MkdirAll(path, 0770)
			if err != nil {
				return
			}
			err = ioutil.WriteFile(path+"/"+header.Filename, data, 0440)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
				return
			}
		}
	default:
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		path := os.TempDir() + "/dotori"
		err = os.MkdirAll(path, 0770)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(path+"/"+header.Filename, data, 0440)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	}
}

// handleUploadAlembic 함수는 Alembic 파일을 처리하는 페이지 이다.
func handleUploadAlembic(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	mimeType := header.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg", "image/png":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		path := os.TempDir() + "/dotori/thumbnail"
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(path+"/"+header.Filename, data, 0666)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	case "video/quicktime", "video/mp4", "video/ogg", "application/ogg":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		path := os.TempDir() + "/dotori/preview"
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(path+"/"+header.Filename, data, 0666)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	case "application/octet-stream":
		ext := filepath.Ext(header.Filename)
		if ext == ".abc" {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
				return
			}
			path := os.TempDir() + "/dotori"
			err = os.MkdirAll(path, 0766)
			if err != nil {
				return
			}
			err = ioutil.WriteFile(path+"/"+header.Filename, data, 0666) // 악성 코드가 들어올 수 있으므로 실행권한은 주지 않는다.
			if err != nil {
				fmt.Fprintf(w, "%v", err)
				return
			}
		}
	default:
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		path := os.TempDir() + "/dotori"
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(path+"/"+header.Filename, data, 0666)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	}
}
