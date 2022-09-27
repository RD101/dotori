package main

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`  // ID
	Name     string             `json:"name" bson:"name"`         // 카테고리 Name
	ParentID string             `json:"parentid" bson:"parentid"` // 빈 문자열이면 Root이다. 빈문자열이 아니라면 서브카테고리이며 이 이름이 부모 카테고리이다.
}

func (i Category) CheckError() error {
	if i.Name == "" {
		return errors.New("need name")
	}
	return nil
}
