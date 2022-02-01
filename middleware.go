package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetObjectIDfromRequestHeader 미들웨어는 리퀘스트헤더에서 ObjectID를 가지고 온다.
func GetObjectIDfromRequestHeader(r *http.Request) (string, error) {
	// 리퀘스트헤더에서 ObjectID를 가지고 온다.
	uri := r.Header.Get("Referer")
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	urlValues, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	var objectID string
	// urlValues에 objectid가 존재하는지 채크한다.
	if value, has := urlValues["objectid"]; has {
		// urlValues["objectid"] 갯수가 1개인지 체크한다.
		if len(value) != 1 {
			return "", errors.New("objectid 값이 1개가 아닙니다")
		}
		objectID = value[0]
	}
	if objectID == "" {
		return "", errors.New("objectid 값이 빈 문자열입니다")
	}
	return objectID, nil
}

// GetTokenFromHeader 함수는 쿠키에서 Token 값을 반환한다.
func GetTokenFromHeader(w http.ResponseWriter, r *http.Request) (Token, error) {
	// Token을 열기위해서 헤더 쿠키에서 필요한 정보를 불러온다.
	sessionToken := ""
	sessionSignkey := ""
	for _, cookie := range r.Cookies() {
		if cookie.Name == "SessionToken" {
			sessionToken = cookie.Value
			continue
		}
		if cookie.Name == "SessionSignKey" {
			sessionSignkey = cookie.Value
			continue
		}
	}
	tk := Token{}
	// Singkey로 Token 정보를 연다.
	token, err := jwt.ParseWithClaims(sessionToken, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(sessionSignkey), nil
	})
	if err != nil {
		return tk, err
	}
	if !token.Valid {
		return tk, errors.New("Token key is not valid")
	}
	return tk, nil
}

// GetAccessLevelFromHeader 함수는 restapi 사용 시 토큰을 체크하고 accesslevel을 반환하는 함수이다.
func GetAccessLevelFromHeader(r *http.Request, client *mongo.Client) (string, error) {
	//header에서 token을 가져온다.
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", errors.New("authorization failed")
	}
	token := auth[1]

	//DB 검색
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := User{}
	err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.AccessLevel, nil
}

// GetUserFromHeader 함수는 restapi 사용 시 토큰을 체크하고 User값을 반환하는 함수이다.
func GetUserFromHeader(r *http.Request, client *mongo.Client) (User, error) {
	user := User{}
	//header에서 token을 가져온다.
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return user, errors.New("authorization failed")
	}
	token := auth[1]
	//DB 검색
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
