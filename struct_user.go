package main

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

// Token 자료구조. JWT 방식을 사용한다. restAPI 사용시 보안체크를 위해 http 헤더에 들어간다.
type Token struct {
	ID          string `json:"id" bson:"id"`                   // 사용자 ID
	AccessLevel string `json:"accesslevel" bson:"accesslevel"` // admin, manager, default
	jwt.StandardClaims
}

// User 는 사용자 자료구조이다.
type User struct {
	ID             string   `json:"id" bson:"id"`                         // 사용자 ID
	Password       string   `json:"password" bson:"password"`             // 암호화된 암호
	Token          string   `json:"token" bson:"token"`                   // JWT 토큰
	SignKey        string   `json:"signkey" bson:"signkey"`               // JWT 토큰을 만들 때 사용하는 SignKey
	AccessLevel    string   `json:"accesslevel" bson:"accesslevel"`       // admin, manager, default
	FavoriteAssets []string `json:"favoriteassets" bson:"favoriteassets"` // 즐겨찾는 어셋 id 리스트
}

// CreateToken 메소드는 토큰을 생성합니다.
func (u *User) CreateToken() error {
	if u.ID == "" {
		return errors.New("ID is an empty string")
	}
	if u.Password == "" {
		return errors.New("Password is an empty string")
	}
	if u.AccessLevel == "" {
		return errors.New("AccessLevel is an empty string")
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{
		ID:          u.ID,
		AccessLevel: u.AccessLevel,
	})
	signKey, err := Encrypt(u.Password)
	if err != nil {
		return err
	}
	u.SignKey = signKey
	tokenString, err := token.SignedString([]byte(signKey))
	if err != nil {
		return err
	}
	u.Token = tokenString
	return nil
}
