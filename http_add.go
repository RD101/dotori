package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// handleAddMaya 함수는 Maya 파일을 추가하는 페이지 이다.
func handleAddMaya(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addmaya", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddMayaItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addmaya-item", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddMayaSubmit 함수는 URL에 objectID를 붙여서 /addmaya 페이지로 redirect한다.
func handleAddMayaSubmit(w http.ResponseWriter, r *http.Request) {
	objectID := bson.NewObjectId().Hex()
	http.Redirect(w, r, fmt.Sprintf("/addmaya-item?objectid=%s", objectID), http.StatusSeeOther)
}

func handleAddNuke(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addnuke", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddHoudini(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addhoudini", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddBlender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addblender", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddAlembic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addalembic", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddUSD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addusd", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddNukeProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addnuke-process", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddHoudiniProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addhoudini-process", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddAlembicProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addalembic-process", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadMaya 함수는 Maya파일을 DB에 업로드하는 페이지를 연다. dropzone에 파일을 올릴 경우 실행된다.
func handleUploadMaya(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	objectID, err := GetObjectIDfromRequestHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rootpath, err := GetRootPath(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	adminsetting, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Close()
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
	// mongoDB objectID를 이용해서 경로 생성
	objectIDpath, err := idToPath(objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, files := range r.MultipartForm.File {
		for _, f := range files {
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
				path := rootpath + objectIDpath + "/originalthumbimg"
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
			case "video/quicktime", "video/mp4", "video/ogg", "application/ogg":
				data, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				path := rootpath + objectIDpath + "/originalthumbmov"
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
			case "application/octet-stream":
				ext := filepath.Ext(f.Filename)
				if ext != ".mb" && ext != ".ma" { // .ma .mb 외에는 허용하지 않는다.
					http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
					return
				}
				data, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Fprintf(w, "%v", err)
					return
				}
				path := rootpath + objectIDpath
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
			default:
				//허용하지 않는 파일 포맷입니다.
				http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
				return
			}
		}
	}
}

func handleUploadMayaOnDB(w http.ResponseWriter, r *http.Request) {
	item := Item{}
	objectID, err := GetObjectIDfromRequestHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.ID = bson.ObjectIdHex(objectID)
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
	item.Status = Ready
	time := time.Now()
	item.CreateTime = time.Format("2006-01-02 15:04:05")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(session)
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
	err = AddItem(session, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/addmaya", http.StatusSeeOther)
}

func handleAddMayaSuccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addmaya-success", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadMaya 함수는 Nuke파일을 DB에 업로드하는 페이지를 연다.
func handleUploadNuke(w http.ResponseWriter, r *http.Request) {
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
		if ext == ".nk" {
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

// handleAddHoudiniProcess 함수는 Houdini 파일을 처리하는 페이지 이다.
func handleUploadHoudini(w http.ResponseWriter, r *http.Request) {
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
