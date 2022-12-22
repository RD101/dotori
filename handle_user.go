package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

func handleProfile(w http.ResponseWriter, r *http.Request) {
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	type recipe struct {
		Token
		User         User
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
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting
	// Profile 페이지를 띄운다.
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "profile", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleFavoriteAssets 는 사용자가 즐겨찾는 아이템 모음 페이지로 연결해주는 함수다.
func handleFavoriteAssets(w http.ResponseWriter, r *http.Request) {
	// token 체크
	token, err := GetTokenFromHeader(w, r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	// mongoDB client 연결
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

	// recipe에 FavoriteItems 담기
	type recipe struct {
		Token
		User           User
		FavoriteItems  []Item
		Adminsetting   Adminsetting
		RootCategories []Category
	}
	rcp := recipe{}
	rcp.Token = token
	adminsetting, err := GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Adminsetting = adminsetting
	rcp.RootCategories, err = GetRootCategories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := GetUser(client, token.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = user
	for _, itemid := range user.FavoriteAssetIDs {
		item, err := GetItem(client, itemid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}
		rcp.FavoriteItems = append(rcp.FavoriteItems, item)
	}

	// 즐겨찾기 모음 페이지 response
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "items-favoriteassets", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSignupSuccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "signup-success", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSignupSubmit(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")
	if id == "" {
		http.Error(w, "ID 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	pw := r.FormValue("Password")
	if pw == "" {
		http.Error(w, "Password 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	if pw != r.FormValue("ConfirmPassword") {
		http.Error(w, "입력받은 2개의 패스워드가 서로 다릅니다", http.StatusBadRequest)
		return
	}
	encryptedPW, err := Encrypt(pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u := User{}
	u.AccessLevel = "default"
	u.ID = id
	u.Password = encryptedPW
	u.NewsNum = 4
	u.TopNum = 4
	err = u.CreateToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	err = AddUser(client, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/signup-success", http.StatusSeeOther)
}

func handleSignin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "signin", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSignOut(w http.ResponseWriter, r *http.Request) {
	tokenKey := http.Cookie{
		Name:   "SessionToken",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &tokenKey)
	signKey := http.Cookie{
		Name:   "SessionSignKey",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &signKey)
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func handleSigninSubmit(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")
	if id == "" {
		http.Error(w, "ID 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	pw := r.FormValue("Password")
	if pw == "" {
		http.Error(w, "Password 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	// DB에서 id로 사용자를 가지고 와서 Password를 비교한다.
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
	u, err := GetUser(client, id)
	if err != nil {
		err := TEMPLATES.ExecuteTemplate(w, "signin-fail", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if err != nil {
		err := TEMPLATES.ExecuteTemplate(w, "signin-fail", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// Token을 쿠키에 저장한다.
	cookieToken := http.Cookie{
		Name:    "SessionToken",
		Value:   u.Token,
		Expires: time.Now().Add(time.Duration(*flagCookieAge) * time.Hour),
	}
	http.SetCookie(w, &cookieToken)
	signKey := http.Cookie{
		Name:    "SessionSignKey",
		Value:   u.SignKey,
		Expires: time.Now().Add(time.Duration(*flagCookieAge) * time.Hour),
	}
	http.SetCookie(w, &signKey)
	// "/" 로 리다이렉션 한다.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleInvalidAccess 는 접근 권한이 없는 uri로의 접속을 시도했을 때 invalid access 페이지를 열어주는 함수다.
func handleInvalidAccess(w http.ResponseWriter, r *http.Request) {
	// token 체크
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

	// 연결할 페이지에서 필요한 데이터 recipe에 정리
	type recipe struct {
		Adminsetting Adminsetting
		Token
	}
	rcp := recipe{}
	rcp.Adminsetting, err = GetAdminSetting(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Token = token

	// 페이지 출력
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "invalidaccess", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAPIUserAutoplay(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	autoplay := q.Get("value")
	if autoplay == "" {
		http.Error(w, "autoplay value must be true or false", http.StatusBadRequest)
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
	user, err := GetUserFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !(user.AccessLevel == "admin" || user.AccessLevel == "default") {
		http.Error(w, "Need permission", http.StatusUnauthorized)
		return
	}
	user.Autoplay = str2bool(autoplay)
	err = SetUser(client, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

func handleAPIUserNewsNum(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	value := q.Get("value")
	if value == "" {
		http.Error(w, "value must be only number", http.StatusBadRequest)
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
	user, err := GetUserFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !(user.AccessLevel == "admin" || user.AccessLevel == "default") {
		http.Error(w, "Need permission", http.StatusUnauthorized)
		return
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if n == 0 { // 만약 사용자 값이 0이면 4로 설정한다.
		n = 4
	}
	user.NewsNum = n
	err = SetUser(client, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

func handleAPIUserTopNum(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	value := q.Get("value")
	if value == "" {
		http.Error(w, "value must be only number", http.StatusBadRequest)
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
	user, err := GetUserFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !(user.AccessLevel == "admin" || user.AccessLevel == "default") {
		http.Error(w, "Need permission", http.StatusUnauthorized)
		return
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if n == 0 { // 만약 사용자 값이 0이면 4로 설정한다.
		n = 4
	}
	user.TopNum = n
	err = SetUser(client, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

func handleAPIUserAccessLevel(w http.ResponseWriter, r *http.Request) {
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
	accessUser, err := GetUserFromHeader(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// RestAPI에 접근한 User의 AccessLevel을 체크한다.
	if !(accessUser.AccessLevel == "admin") {
		http.Error(w, "Need permission", http.StatusUnauthorized)
		return
	}
	// client로 부터 받은 데이터를 파싱한다.
	data := User{}
	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// 파싱된 ID를 이용해서 권한을 변경할 유저의 토큰과 엑세스 레벨을 수정하고 업데이트 한다.
	user, err := GetUser(client, data.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user.AccessLevel = data.AccessLevel
	err = user.CheckAccessLevel()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.CreateToken()
	err = SetUser(client, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	resp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
