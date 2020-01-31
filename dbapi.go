package main

import (
	"strings"

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

// SearchTags 는 itemType, tag를 입력받아 tag의 값이 일치하면 반환하는 함수입니다.
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

// Search 는 itemType, words를 입력받아 해당 아이템을 검색한다.
func Search(session *mgo.Session, itemType string, words string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(itemType)
	var results []Item
	wordsQueries := []bson.M{}
	for _, word := range strings.Split(words, " ") {
		if word == "" {
			continue
		}
		querys := []bson.M{}
		querys = append(querys, bson.M{"_id": bson.ObjectIdHex(word)})
		querys = append(querys, bson.M{"author": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"tags": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"description": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"type": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"inputpath": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"outputpath": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"createtime": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"updatetime": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"attributes.*": word})
		wordsQueries = append(wordsQueries, bson.M{"$or": querys})
	}
	// 사용률이 많은 소스가 위로 출력되도록 한다.
	q := bson.M{"$and": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.
	err := c.Find(q).Sort("-usingrate").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// TotalPage 함수는 아이템의 갯수를 입력받아 필요한 총 페이지 수를 구한다.
func TotalPage(itemNum int) int {
	page := itemNum / *flagPagenum
	if itemNum%*flagPagenum != 0 {
		page++
	}
	return page
}

// SearchPage 는 itemType, words, 해당 page를 입력받아 해당 아이템을 검색한다.
func SearchPage(session *mgo.Session, itemType string, words string, page int) ([]Item, error) {
	var results []Item
	if words == "" {
		return results, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(itemType)

	wordsQueries := []bson.M{}
	for _, word := range strings.Split(words, " ") {
		if word == "" {
			continue
		}
		querys := []bson.M{}
		querys = append(querys, bson.M{"_id": bson.ObjectIdHex(word)})
		querys = append(querys, bson.M{"author": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"tags": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"description": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"type": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"inputpath": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"outputpath": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"createtime": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"updatetime": &bson.RegEx{Pattern: word, Options: "i"}})
		querys = append(querys, bson.M{"attributes.*": word})
		wordsQueries = append(wordsQueries, bson.M{"$or": querys})
	}
	// 사용률이 많은 소스가 위로 출력되도록 한다.
	q := bson.M{"$and": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.
	err := c.Find(q).Sort("-usingrate").Skip(page - 1).Limit(*flagPagenum).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchItem 은 컬렉션 이름(itemType)과 id를 받아서, 해당 컬렉션에서 id가 일치하는 item을 검색, 반환한다.
func SearchItem(session *mgo.Session, itemType, id string) (Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(itemType)
	var result Item
	err := c.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
