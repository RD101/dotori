package main

import (
	"gopkg.in/mgo.v2"
)

// AddItem 은 데이터베이스에 Item을 넣는 함수이다.
func AddItem(session *mgo.Session, i Item) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("dotori").C(i.Type)
	err := c.Insert(i)
	if err != nil {
		return err
	}
	return nil
}
