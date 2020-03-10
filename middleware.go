package main

import (
	"errors"
	"net/http"
	"net/url"
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
