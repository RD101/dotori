package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleAPIDeleteItem(w http.ResponseWriter, r *http.Request) {

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
	//accesslevel 체크
	accesslevel, err := GetAccessLevelFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accesslevel != "admin" {
		http.Error(w, "삭제 권한이 없는 계정입니다", http.StatusUnauthorized)
		return
	}
	// 실제 데이터 삭제
	err = RmData(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// DB에서 데이터 삭제
	err = RmItem(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = RmFavoriteItem(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIPostItem(w http.ResponseWriter, r *http.Request) {
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
	//accesslevel 체크
	accesslevel, err := GetAccessLevelFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accesslevel != "default" && accesslevel != "manager" && accesslevel != "admin" {
		http.Error(w, "등록 권한이 없는 계정입니다", http.StatusUnauthorized)
		return
	}

	// 아이템 생성
	i := Item{}
	i.ID = primitive.NewObjectID()
	// 아이템 정보 Parsing
	itemtype := r.FormValue("itemtype")
	if itemtype == "" {
		http.Error(w, "itemtype을 설정해주세요", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "title을 설정해주세요", http.StatusBadRequest)
		return
	}
	author := r.FormValue("author")
	if author == "" {
		http.Error(w, "author를 설정해주세요", http.StatusBadRequest)
		return
	}
	description := r.FormValue("description")
	if description == "" {
		http.Error(w, "description을 설정해주세요", http.StatusBadRequest)
		return
	}
	tags := Str2List(r.FormValue("tags"))
	if len(tags) == 0 {
		http.Error(w, "tags를 설정해주세요", http.StatusBadRequest)
		return
	}
	attributes, err := StringToMap(r.FormValue("attributes"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	i.ItemType = itemtype
	i.Title = title
	i.Author = author
	i.Description = description
	i.Tags = tags
	i.Attributes = attributes
	i.Status = "ready"
	i.Logs = append(i.Logs, "아이템이 생성되었습니다.")
	// admin setting에서 rootpath를 가져와 경로를 생성한다.
	rootpath, err := GetRootPath(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	objIDpath, err := idToPath(i.ID.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	i.InputThumbnailImgPath = rootpath + objIDpath + "/originalthumbimg/"
	i.InputThumbnailClipPath = rootpath + objIDpath + "/originalthumbmov/"
	i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
	i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
	i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
	i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
	i.OutputDataPath = rootpath + objIDpath + "/data/"

	// 아이템 추가
	err = i.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = AddItem(client, i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 아이템에 파일 업데이트
	if itemtype == "alembic" {
		uploadAlembicFile(w, r, i.ID.Hex())
	}
	if itemtype == "blender" {
		uploadBlenderFile(w, r, i.ID.Hex())
	}
	if itemtype == "footage" {
		uploadFootageFile(w, r, i.ID.Hex())
	}
	if itemtype == "fusion360" {
		uploadFusion360File(w, r, i.ID.Hex())
	}
	if itemtype == "hdri" {
		uploadHDRIFile(w, r, i.ID.Hex())
	}
	if itemtype == "houdini" {
		uploadHoudiniFile(w, r, i.ID.Hex())
	}
	if itemtype == "hwp" {
		uploadHwpFile(w, r, i.ID.Hex())
	}
	if itemtype == "katana" {
		uploadKatanaFile(w, r, i.ID.Hex())
	}
	if itemtype == "lut" {
		uploadLutFile(w, r, i.ID.Hex())
	}
	if itemtype == "max" {
		uploadMaxFile(w, r, i.ID.Hex())
	}
	if itemtype == "maya" {
		uploadMayaFile(w, r, i.ID.Hex())
	}
	if itemtype == "modo" {
		uploadModoFile(w, r, i.ID.Hex())
	}
	if itemtype == "nuke" {
		uploadNukeFile(w, r, i.ID.Hex())
	}
	if itemtype == "openvdb" {
		uploadOpenVDBFile(w, r, i.ID.Hex())
	}
	if itemtype == "pdf" {
		uploadPdfFile(w, r, i.ID.Hex())
	}
	if itemtype == "ppt" {
		uploadPptFile(w, r, i.ID.Hex())
	}
	if itemtype == "sound" {
		uploadSoundFile(w, r, i.ID.Hex())
	}
	if itemtype == "texture" {
		uploadClipFile(w, r, i.ID.Hex())
	}
	if itemtype == "unreal" {
		uploadUnrealFile(w, r, i.ID.Hex())
	}
	if itemtype == "usd" {
		uploadUSDFile(w, r, i.ID.Hex())
	}

	// Response
	item, err := GetItem(client, i.ID.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIGetItem(w http.ResponseWriter, r *http.Request) {
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
	i, err := GetItem(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIPutItem(w http.ResponseWriter, r *http.Request) {
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
	//accesslevel 체크
	accesslevel, err := GetAccessLevelFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accesslevel != "manager" && accesslevel != "admin" {
		http.Error(w, "need permission", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	item := Item{}
	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&item)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	err = SetItem(client, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISearch 는 아이템을 검색하는 함수입니다.
func handleAPISearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	itemtype := r.FormValue("itemtype")
	if itemtype == "" {
		http.Error(w, "itemtype을 설정해주세요", http.StatusBadRequest)
		return
	}
	searchword := r.FormValue("searchword")
	if searchword == "" {
		http.Error(w, "searchword를 설정해주세요", http.StatusBadRequest)
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
	item, err := Search(client, itemtype, searchword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIAdminSetting(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
		admin, err := GetAdminSetting(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(admin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	http.Error(w, "Not Supported Method", http.StatusMethodNotAllowed)
}

func handleAPIUsingRate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		itemtype := r.FormValue("itemtype")
		if itemtype == "" {
			http.Error(w, "itemtype을 입력해주세요", http.StatusBadRequest)
			return
		}
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "id를 입력해주세요", http.StatusBadRequest)
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
		usingrate, err := UpdateUsingRate(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(usingrate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	http.Error(w, "Not Supported Method", http.StatusMethodNotAllowed)
}

// handleAPIRecentItem 는 최근생성된 아이템들을 반환하는 함수임니다.
func handleAPIRecentItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		r.ParseForm()
		recentlypage, err := strconv.ParseInt(r.FormValue("recentlypage"), 10, 64)
		if err != nil {
			http.Error(w, "recentlypage를 입력해주세요", http.StatusBadRequest)
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
		usingrate, err := GetRecentlyCreatedItems(client, 4, recentlypage) // 해당페이지(recentlypage)의 4개 아이템을 가져온다.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(usingrate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	http.Error(w, "Not Supported Method", http.StatusMethodNotAllowed)
}

// handleAPITopUsingItem 는 많이 사용되는 아이템들을 반환하는 함수임니다.
func handleAPITopUsingItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		r.ParseForm()
		topusingpage, err := strconv.ParseInt(r.FormValue("usingpage"), 10, 64)
		if err != nil {
			http.Error(w, "usingpage를 입력해주세요", http.StatusBadRequest)
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
		usingrate, err := GetTopUsingItems(client, 4, topusingpage) // 해당페이지(topusingpage)의 4개 아이템을 가져온다.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(usingrate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	http.Error(w, "Not Supported Method", http.StatusMethodNotAllowed)
}

// handleAPIFavoriteAsset는 FavoriteAssetIds에 아이템 id를 추가하거나 제거하는 함수다.
func handleAPIFavoriteAsset(w http.ResponseWriter, r *http.Request) {

	// mongoDB Client 생성
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
	// accesslevel 체크
	accesslevel, err := GetAccessLevelFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accesslevel != "default" && accesslevel != "manager" && accesslevel != "admin" {
		http.Error(w, "즐겨찾기 수정 권한이 없습니다", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		// Get : Get FavoriteAssetIDs

		// 전송받은 데이터 parsing
		q := r.URL.Query()
		userid := q.Get("userid")
		if userid == "" {
			http.Error(w, "URL에 userid를 입력해주세요", http.StatusBadRequest)
			return
		}

		// Delete itemid from FavoriteAssetIds of User
		user := User{}
		user, err = GetUser(client, userid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		favoriteAssetIds := user.FavoriteAssetIDs
		reponseIds := make(map[string][]string)
		reponseIds["favoriteAssetIds"] = favoriteAssetIds

		data, err := json.Marshal(reponseIds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return

	} else if r.Method == http.MethodPost {
		// POST : FavoriteAssetIDs 자료구조에 itemid를 추가

		// 전송받은 데이터 parsing
		itemid := r.FormValue("itemid")
		if itemid == "" {
			http.Error(w, "itemid를 설정해주세요", http.StatusBadRequest)
			return
		}
		userid := r.FormValue("userid")
		if userid == "" {
			http.Error(w, "userid를 설정해주세요", http.StatusBadRequest)
		}

		// Add itemid to FavoriteAssetIds of User
		user := User{}
		user, err = GetUser(client, userid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for i := 0; i < len(user.FavoriteAssetIDs); i++ {
			if itemid == user.FavoriteAssetIDs[i] {
				http.Error(w, "즐겨찾기 목록에 이미 존재하는 itemid입니다", http.StatusBadRequest)
				return
			}
		}
		user.FavoriteAssetIDs = append(user.FavoriteAssetIDs, itemid)
		err = SetUser(client, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		user, err = GetUser(client, userid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return

	} else if r.Method == http.MethodDelete {
		// DELETE : FavoriteAssetsId 자료구조에 itemid를 추가

		// 전송받은 데이터 parsing
		q := r.URL.Query()
		itemid := q.Get("itemid")
		if itemid == "" {
			http.Error(w, "URL에 itemid를 입력해주세요", http.StatusBadRequest)
			return
		}
		userid := q.Get("userid")
		if userid == "" {
			http.Error(w, "URL에 userid를 입력해주세요", http.StatusBadRequest)
			return
		}

		// Delete itemid from FavoriteAssetIds of User
		user := User{}
		user, err = GetUser(client, userid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		deleteBool := false
		for i := 0; i < len(user.FavoriteAssetIDs); i++ {
			if itemid == user.FavoriteAssetIDs[i] {
				user.FavoriteAssetIDs = append(user.FavoriteAssetIDs[:i], user.FavoriteAssetIDs[i+1:]...)
				deleteBool = true
			}
		}
		if !deleteBool {
			http.Error(w, "즐겨찾기에 존재하지 않는 itemid입니다", http.StatusBadRequest)
			return
		}
		err = SetUser(client, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		user, err = GetUser(client, userid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

}

// handleAPIInitPassword 함수는 rest API를 이용하여 사용자의 비밀번호를 초기화하는 함수이다.
func handleAPIInitPassword(w http.ResponseWriter, r *http.Request) {
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
	//accesslevel 체크
	accesslevel, err := GetAccessLevelFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accesslevel != "admin" {
		http.Error(w, "사용자의 패스워드 초기화 권한이 없는 계정입니다", http.StatusUnauthorized)
		return
	}

	adminSetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encryptedPW, err := Encrypt(adminSetting.InitPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u := User{}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&u)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if u.ID == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}

	user, err := GetUser(client, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = encryptedPW
	user.CreateToken()
	err = SetUser(client, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
