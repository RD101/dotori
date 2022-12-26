package main

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleUsers 함수는 Users 페이지를 띄워주는 함수이다.
func handleUsers(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
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
	// Access level 체크

	type recipe struct {
		User         User
		Token        Token
		Adminsetting Adminsetting
		Users        []User
	}
	rcp := recipe{}
	rcp.User, err = GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Token = token

	if !(rcp.User.AccessLevel == "admin" || rcp.User.AccessLevel == "manager") {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}

	rcp.Adminsetting, err = GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Users, err = GetAllUsers(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "users", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
