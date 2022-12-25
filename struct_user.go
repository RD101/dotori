package main

import (
	"errors"

	jwt "github.com/golang-jwt/jwt"
)

// Token 자료구조. JWT 방식을 사용한다. restAPI 사용시 보안체크를 위해 http 헤더에 들어간다.
type Token struct {
	ID          string `json:"id" bson:"id"`                   // 사용자 ID
	AccessLevel string `json:"accesslevel" bson:"accesslevel"` // admin, manager, default
	jwt.StandardClaims
}

// User 는 사용자 자료구조이다.
type User struct {
	ID               string   `json:"id"`               // 사용자 ID
	Password         string   `json:"-"`                // 암호화된 암호. json으로 반환되지 않도록 설정한다.
	Token            string   `json:"-"`                // JWT 토큰. json으로 반환되지 않도록 설정한다.
	SignKey          string   `json:"-"`                // JWT 토큰을 만들 때 사용하는 SignKey. json으로 반환되지 않도록 설정한다.
	AccessLevel      string   `json:"accesslevel"`      // admin, manager, default
	FavoriteAssetIDs []string `json:"favoriteassetids"` // 즐겨찾는 어셋 id 리스트
	Autoplay         bool     `json:"autoplay"`         // 영상 자동재생 옵션
	NewsNum          int      `json:"newsnum"`          // 새로 추가된 에셋 표시갯수
	TopNum           int      `json:"topnum"`           // 자주 사용하는 에셋 표시갯수
}

func (u *User) CheckAccessLevel() error {
	switch u.AccessLevel {
	case "admin":
		return nil
	case "manager":
		return nil
	case "default":
		return nil
	default:
		return errors.New("not support " + u.AccessLevel)
	}
}

// CreateToken 메소드는 토큰을 생성합니다.
func (u *User) CreateToken() error {
	if u.ID == "" {
		return errors.New("ID is an empty string")
	}
	if u.Password == "" {
		return errors.New("password is an empty string")
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
