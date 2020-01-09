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

// RmItem 는 컬렉션 이름과 id를 받아서, 해당 컬렉션에서 id가 일치하는 Item을 삭제한다.
func RmItem(session *mgo.Session, itemType, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(itemType)
	err := c.RemoveId(bson.ObjectIdHex(id))
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

// SearchTag는 itemType, tag를 입력받아 tag의 값이 일치하면 반환하는 함수입니다.
func SearchTags(session *mgo.Session, itemType string, tag string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(itemType)
	var results []Item
	err := c.Find(bson.M{"tags": tag}).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}