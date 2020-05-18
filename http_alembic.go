package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
