package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func addMayaItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".ma" && ext != ".mb" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addMaxItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".max" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addFusion360ItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".f3d" && ext != ".step" && ext != ".stp" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addOpenVDBItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".vdb" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addModoItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".lxo" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addKatanaItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".katana" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addHoudiniItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".hda" && ext != ".hip" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addBlenderItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".blend" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addLutItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		// "lut", "blut", "cms", "csp", "cub", "vfz"
		if ext != ".cube" && ext != ".3dl" && ext != ".vf" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 4. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addClipItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagFPS == "" {
		log.Fatal("fps가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Fps = *flagFPS
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".mov" && ext != ".mp4" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addPdfItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".pdf" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addIesItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".ies" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addPptItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".key" && ext != ".ppt" && ext != ".pptx" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addSoundItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".wav" && ext != ".mp3" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addTextureItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".exr" && ext != ".png" && ext != ".jpg" && ext != ".tga" && ext != ".tif" && ext != ".tiff" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addHwpItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".hwp" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addHdriItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	if *flagInColorspace == "" {
		log.Fatal("incolorspace가 빈 문자열입니다")
	}
	if *flagOutColorspace == "" {
		log.Fatal("outcolorspace가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InColorspace = *flagInColorspace
	i.OutColorspace = *flagOutColorspace
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".hdr" && ext != ".hdri" && ext != ".exr" { // .hdr .hdri .exr 외에는 허용하지 않는다.
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addUnrealItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".cpp" && ext != ".uasset" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addFootageItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagFPS == "" {
		log.Fatal("fps가 빈 문자열입니다")
	}
	if *flagInColorspace == "" {
		log.Fatal("incolorspace가 빈 문자열입니다")
	}
	if *flagOutColorspace == "" {
		log.Fatal("outcolorspace가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Fps = *flagFPS
	i.InColorspace = *flagInColorspace
	i.OutColorspace = *flagOutColorspace
	i.Tags = Str2List(*flagTag)
	i.Premultiply = *flagPremultiply
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"
	i.OutputProxyImgPath = rootpath + objIDpath + "/proxy/"

	// 1. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext := filepath.Ext(path)
		if ext != ".dpx" && ext != ".exr" && ext != ".zip" && ext != ".jpg" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addUSDItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".usd" && ext != ".usdc" && ext != ".usda" && ext != ".usdz" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 4. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addAlembicItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".abc" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func addNukeItemCmd() {
	if *flagTitle == "" {
		log.Fatal("title이 빈 문자열입니다")
	}
	if *flagAuthor == "" {
		log.Fatal("author이 빈 문자열입니다")
	}
	if *flagDescription == "" {
		log.Fatal("description이 빈 문자열입니다")
	}
	if *flagTag == "" {
		log.Fatal("tag가 빈 문자열입니다")
	}
	if *flagInputThumbImgPath == "" {
		log.Fatal("inputthumbimgpath가 빈 문자열입니다")
	}
	if *flagInputThumbClipPath == "" {
		log.Fatal("inputthumbclippath가 빈 문자열입니다")
	}
	if *flagInputDataPath == "" {
		log.Fatal("inputdatapath가 빈 문자열입니다")
	}
	i := Item{}
	i.ID = primitive.NewObjectID()
	i.ItemType = *flagItemType
	i.Title = *flagTitle
	i.Author = *flagAuthor
	i.Description = *flagDescription
	i.Tags = Str2List(*flagTag)
	attr, err := StringToMap(*flagAttributes)
	if err != nil {
		log.Fatal(err)
	}
	i.Attributes = attr
	i.InputThumbnailImgPath = *flagInputThumbImgPath
	i.InputThumbnailClipPath = *flagInputThumbClipPath
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	i.ThumbImgUploaded = false
	i.ThumbClipUploaded = false
	i.DataUploaded = false

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// admin settin에서 rootpath를 가져와서 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		log.Fatal(err)
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"
	i.OutputProxyImgPath = rootpath + objIDpath + "/proxy/"

	// 1. 썸네일 이미지
	// 썸네일 이미지 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err := os.Stat(*flagInputThumbImgPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext := filepath.Ext(*flagInputThumbImgPath)
	if ext != ".jpg" && ext != ".png" {
		log.Fatal("지원하지 않는 썸네일 이미지 포맷입니다")
	}
	// 존재하고 유효하면 ThumbImgUploaded true로 바꾸기
	i.ThumbImgUploaded = true

	// 2. 썸네일 클립
	// 썸네일 클립 경로에 실재 파일이 존재하는지 체크.
	err = FileExists(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	// 레귤러 파일이 아니면 에러처리 한다.
	stat, err = os.Stat(*flagInputThumbClipPath)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Fatal("폴더, 심볼릭 링크 등은 복사할 수 없습니다")
	}
	// 유효한 파일인지 체크.
	ext = filepath.Ext(*flagInputThumbClipPath)
	if ext != ".mov" && ext != ".mp4" && ext != ".ogg" {
		log.Fatal("지원하지 않는 썸네일 클립 포맷입니다.")
	}
	// 존재하고 유효하면 ThumbClipUploaded true로 바꾸기
	i.ThumbClipUploaded = true

	// 3. DB에 Asset 추가
	err = i.CheckError()
	if err != nil {
		log.Fatal(err)
	}
	err = AddItem(client, i)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 데이터 복사
	datapaths := QuotesPaths2Paths(*flagInputDataPath)
	var filteredPaths []string
	for _, path := range datapaths {
		if HasWildcard(path) {
			// 파일명에 와일드카드(?,*)가 존재할 때
			matches, err := filepath.Glob(*flagInputDataPath)
			if err != nil {
				log.Fatal(err)
			}
			filteredPaths = append(filteredPaths, matches...)
		} else {
			filteredPaths = append(filteredPaths, path)
		}
	}
	for _, path := range filteredPaths {
		// 데이터 경로에 실재 파일이 존재하는지 체크.
		err = FileExists(path)
		if err != nil {
			log.Fatal(err)
		}
		// 레귤러 파일이 아니면 에러처리 한다.
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !stat.Mode().IsRegular() {
			// 레귤러 파일이 아니면 복사할 수 없다.(ex. 폴더, symlinks, 디바이스 등등) cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			log.Fatal("폴더, 심볼릭 링크 등등은 복사할 수 없습니다")
		}
		// 유효한 파일인지 체크.
		ext = filepath.Ext(path)
		if ext != ".nk" && ext != ".gizmo" && ext != ".zip" {
			log.Fatal("지원하지 않는 데이터 포맷입니다.")
		}
		// 있으면 OutputData 경로로 복사하기
		err = copyFile(path, i.OutputDataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Asset status 업데이트
	updateItem, err := GetItem(client, i.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}
	// file upload 완료를 의미하는 status로 변경
	updateItem.DataUploaded = true
	updateItem.Status = "fileuploaded"
	err = SetItem(client, updateItem)
	if err != nil {
		log.Print(err)
	}
}

func rmItemCmd() {
	if *flagItemID == "" {
		log.Fatal("id가 빈 문자열 입니다")
	}
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	// 실제 데이터를 폴더 트리에서 삭제
	err = RmData(client, *flagItemID)
	if err != nil {
		log.Print(err)
	}
	// DB에서 데이터 삭제
	err = RmItem(client, *flagItemID)
	if err != nil {
		log.Print(err)
	}
	err = RmFavoriteItem(client, *flagItemID)
	if err != nil {
		log.Print(err)
	}

}
