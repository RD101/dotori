package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

// handleAddMayaProcess 함수는 Maya 파일을 처리하는 페이지 이다.
func handleAddMayaProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "addmaya-process", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadMaya 함수는 Maya파일을 DB에 업로드하는 페이지를 연다.
func handleUploadMaya(w http.ResponseWriter, r *http.Request) {
	//dropzone setting
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
		tempDir, err := ioutil.TempDir("", "")
		fmt.Println(tempDir)
		path := filepath.Dir(tempDir) + "/dotori/thumbnail"
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return
		}
		fmt.Println(path)
		err = ioutil.WriteFile(path, data, 0666)
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
		tempDir, err := ioutil.TempDir("", "")
		fmt.Println(tempDir)
		path := filepath.Dir(tempDir) + "/dotori/preview"
		err = os.MkdirAll(path, 0777) //0766
		if err != nil {
			return
		}
		fmt.Println(path)
		err = ioutil.WriteFile(path, data, 0777) //0666
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	case "application/octet-stream":
		//ext := filepath.Ext()
	default:
		//컨텐츠가 따로 있는게 편할지 같이 있는게 편할지
	}

	log.Println(mimeType)
}
