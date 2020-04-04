package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2"
)

func handleAPIItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		i := Item{}
		i.ID = primitive.NewObjectID()
		//ParseForm parses the raw query from the URL and updates r.Form.
		r.ParseForm()
		for key, values := range r.PostForm {
			switch key {
			case "itemtype":
				if len(values) != 1 {
					http.Error(w, "URL에 itemtype을 입력해주세요", http.StatusBadRequest)
					return
				}
				i.ItemType = values[0]
			case "author":
				if len(values) != 1 {
					http.Error(w, "URL에 author를 입력해주세요", http.StatusBadRequest)
					return
				}
				i.Author = values[0]
			}
		}

		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		// admin settin에서 rootpath를 가져와서 경로를 생성한다.
		// rootpath, err := GetRootPath(session)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// objIDpath, err := idToPath(i.ID.Hex())
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// i.InputThumbnailImgPath = rootpath + objIDpath + "/originalthumbimg/"
		// i.InputThumbnailClipPath = rootpath + objIDpath + "/originalthumbmov/"
		// i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
		// i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
		// i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
		// i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
		// i.OutputDataPath = rootpath + objIDpath + "/data/"

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
		data, _ := json.Marshal(i)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	} else if r.Method == http.MethodDelete {
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
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = RmItem(client, itemtype, id)
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
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		i, err := GetItem(client, itemtype, id)
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	item, err := Search(session, itemtype, searchword)
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
