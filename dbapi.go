package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddItem 은 데이터베이스에 Item을 넣는 함수이다.
func AddItem(session *mgo.Session, i Item) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(i.Type)
	err := c.Insert(i)
	if err != nil {
		return err
	}
	return nil
}

// allItems는 DB에서 전체 아이템 정보를 가져오는 함수입니다.
func allItems(session *mgo.Session, itemType string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(itemType)
	var result []Item
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
