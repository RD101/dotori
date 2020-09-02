package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleAPIItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
		tags := Str2Tags(r.FormValue("tags"))
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
		return

	} else if r.Method == http.MethodDelete {
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
		// 삭제 함수 호출
		err = RmItem(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	} else if r.Method == http.MethodGet {
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
		return
	} else {
		http.Error(w, "Not Supported Method", http.StatusMethodNotAllowed)
		return
	}
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
	return
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
	return
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
	return
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
	return
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
	return
}
