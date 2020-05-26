package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleDownloadItem(w http.ResponseWriter, r *http.Request) {
	_, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	q := r.URL.Query()
	itemtype := q.Get("itemtype")
	id := q.Get("id")
	if itemtype == "" {
		http.Error(w, "URL에 itemtype을 입력해주세요", http.StatusBadRequest)
		return
	}
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
	item, err := GetItem(client, itemtype, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 임시 디렉토리를 생성한다.
	tempDir, err := ioutil.TempDir("", "zip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)
	zipFileName := item.ID.Hex() + ".zip"
	zipFilePath, err := genZipfile(tempDir, zipFileName, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Using Rate(사용률)을 업데이트 한다.
	_, err = UpdateUsingRate(client, itemtype, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/zip")
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s", zipFileName))
	http.ServeFile(w, r, zipFilePath)
}

func genZipfile(tempDir, zipFileName string, item Item) (string, error) {
	// zip 파일을 생성한다.
	zipFilePath := tempDir + "/" + zipFileName
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return zipFilePath, err
	}
	defer zipFile.Close()
	// zip 파일에 쓰기할 준비를 한다.
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// item의 data경로에 존재하는 파일 리스트를 불러온다.
	dataPath := item.OutputDataPath
	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return zipFilePath, err
	}

	// 데이터 파일을 돌면서 zip 파일에 데이터 파일 추가한다.
	for _, f := range files {
		fileName, err := os.Open(dataPath + f.Name())
		if err != nil {
			return zipFilePath, err
		}
		defer fileName.Close()

		// 파일정보를 가지고 온다.
		info, err := fileName.Stat()
		if err != nil {
			return zipFilePath, err
		}
		// 압축할 때 zip 파일에 파일정보를 헤더로 설정한다.
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return zipFilePath, err
		}
		header.Method = zip.Deflate
		// 헤더정보를 zip 파일에 쓴다.
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return zipFilePath, err
		}
		// 파일의 실제 내용을 zip 파일에 복사한다.
		_, err = io.Copy(writer, fileName)
		if err != nil {
			return zipFilePath, err
		}
	}
	return zipFilePath, nil
}
