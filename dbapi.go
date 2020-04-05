package main

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddItem 은 데이터베이스에 Item을 넣는 함수이다.
// func AddItem(session *mgo.Session, i Item) error {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C(i.ItemType)
// 	err := c.Insert(i)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// GetItem 은 데이터베이스에 Item을 가지고 오는 함수이다.
// func GetItem(session *mgo.Session, itemType, id string) (Item, error) {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C(itemType)
// 	var result Item
// 	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
// 	if err != nil {
// 		return result, err
// 	}
// 	return result, nil
// }

// RmItem 는 컬렉션 이름과 id를 받아서, 해당 컬렉션에서 id가 일치하는 Item을 삭제한다.
// func RmItem(session *mgo.Session, itemType, id string) error {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C(itemType)
// 	err := c.RemoveId(bson.ObjectIdHex(id))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// UpdateItem 은 컬렉션 이름과 Item을 받아서, Item을 업데이트한다.
// func UpdateItem(session *mgo.Session, itemType string, item Item) error {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C(itemType)
// 	err := c.Update(bson.M{"_id": item.ID}, item)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// allItems는 DB에서 전체 아이템 정보를 가져오는 함수입니다.
// func allItems(session *mgo.Session, itemType string) ([]Item, error) {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C(itemType)
// 	var result []Item
// 	err := c.Find(bson.M{}).All(&result)
// 	if err != nil {
// 		return result, err
// 	}
// 	return result, nil
// }

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
		//"tag:"가 앞에 붙어있으면 태그에서 검색한다.
		if strings.HasPrefix(word, "tag:") {
			querys = append(querys, bson.M{"tags": strings.TrimPrefix(word, "tag:")})
		} else if strings.Contains(word, ":") {
			key := strings.Split(word, ":")[0]
			value := strings.Split(word, ":")[1]
			querys = append(querys, bson.M{"attributes." + key: &bson.RegEx{Pattern: value, Options: "i"}})
		} else {
			querys = append(querys, bson.M{"author": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"tags": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"description": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"type": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"inputpath": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"outputpath": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"createtime": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"updatetime": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"attributes.*": word})
		}
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

// SearchPage 는 itemType, words, 해당 page를 입력받아 해당 아이템을 검색한다. 검색된 아이템과 그 개수를 반환한다.
func SearchPage(session *mgo.Session, itemType string, words string, page, limitnum int) (int, int, []Item, error) {
	var results []Item
	//검색어가 존재하지 않으면 빈 결과를 반환한다.
	if words == "" {
		return 0, 0, results, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C(itemType)

	wordsQueries := []bson.M{}
	for _, word := range strings.Split(words, " ") {
		if word == "" {
			continue
		}
		querys := []bson.M{}
		//"tag:"가 앞에 붙어있으면 태그에서 검색한다.
		if strings.HasPrefix(word, "tag:") {
			querys = append(querys, bson.M{"tags": strings.TrimPrefix(word, "tag:")})
		} else if strings.Contains(word, ":") {
			key := strings.Split(word, ":")[0]
			value := strings.Split(word, ":")[1]
			querys = append(querys, bson.M{"attributes." + key: &bson.RegEx{Pattern: value, Options: "i"}})
		} else {
			querys = append(querys, bson.M{"author": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"tags": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"description": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"type": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"inputpath": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"outputpath": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"createtime": &bson.RegEx{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"updatetime": &bson.RegEx{Pattern: word, Options: "i"}})
		}
		wordsQueries = append(wordsQueries, bson.M{"$or": querys})
	}
	// 사용률이 많은 소스가 위로 출력되도록 한다.
	q := bson.M{"$and": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.
	err := c.Find(q).Sort("-usingrate").Skip((page - 1) * limitnum).Limit(limitnum).All(&results)
	if err != nil {
		return 0, 0, nil, err
	}
	totalNum, err := c.Find(q).Count()
	if err != nil {
		return 0, 0, nil, err
	}
	return TotalPage(totalNum, limitnum), totalNum, results, nil
}

// SearchItem 은 컬렉션 이름(itemType)과 id를 받아서, 해당 컬렉션에서 id가 일치하는 item을 검색, 반환한다.
// func SearchItem(session *mgo.Session, itemType, id string) (Item, error) {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C(itemType)
// 	var result Item
// 	err := c.FindId(bson.ObjectIdHex(id)).One(&result)
// 	if err != nil {
// 		return result, err
// 	}
// 	return result, nil
// }

// SetAdminSetting 은 입력받은 어드민셋팅으로 업데이트한다.
func SetAdminSetting(session *mgo.Session, a Adminsetting) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("setting.admin")
	num, err := c.Find(bson.M{"id": "setting.admin"}).Count()
	if err != nil {
		return err
	}
	a.ID = "setting.admin"
	if num == 0 {
		err = c.Insert(a)
		if err != nil {
			return err
		}
		return nil
	}
	err = c.Update(bson.M{"id": "setting.admin"}, a)
	if err != nil {
		return err
	}
	return nil
}

// GetAdminSetting 은 관리자 셋팅값을 가지고 온다.
// func GetAdminSetting(session *mgo.Session) (Adminsetting, error) {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C("setting.admin")
// 	var result Adminsetting
// 	err := c.Find(bson.M{"id": "setting.admin"}).One(&result)
// 	if err != nil {
// 		if err == mgo.ErrNotFound {
// 			return Adminsetting{}, nil
// 		}
// 		return Adminsetting{}, err
// 	}
// 	return result, nil
// }

// GetOngoingProcess 는 처리 중인 아이템을 가져온다.
func GetOngoingProcess(session *mgo.Session) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	var results []Item
	//콜렉션 리스트를 가져온다.
	collections, err := session.DB(*flagDBName).CollectionNames()
	if err != nil {
		return results, err
	}
	// 콜렉션마다 돌면서 Status가 Done이 아닌 아이템을 가져온다.
	for _, c := range collections {
		var items []Item
		if c == "system.indexs" { //mongodb의 기본 컬렉션. 제외한다.
			continue
		}
		if c == "setting.admin" { //admin setting값을 저장하는 컬렉션. 제외한다.
			continue
		}
		err = session.DB(*flagDBName).C(c).Find(bson.M{"status": bson.M{"$ne": Done}}).All(&items)
		if err != nil {
			return results, err
		}
		results = append(results, items...)
	}
	return results, nil
}

// GetReadyItem 은 DB에서 ready상태인 Item을 하나 가져온다.
func GetReadyItem(session *mgo.Session) (Item, error) {
	session.SetMode(mgo.Monotonic, true)
	var result Item
	collections, err := session.DB(*flagDBName).CollectionNames()
	if err != nil {
		return result, err
	}
	// 컬렉션을 for문 돌면서 Ready 상태인 Item을 찾는다.
	for _, c := range collections {
		if c == "setting.admin" { // setting.admin 컬렉션은 제외한다.
			continue
		}
		if c == "system.indexs" { //mongodb의 기본 컬렉션. 제외한다.
			continue
		}
		cur := session.DB(*flagDBName).C(c)
		// 해당 컬렉션에 ready상태인 Item이 없으면 다음 컬렉션을 체크한다
		num, err := cur.Find(bson.M{"status": Ready}).Count()
		if err != nil {
			return result, err
		}
		if num == 0 {
			continue
		}
		// ready상태인 Item이 있다면 가져와서 Status를 업데이트 한다.
		err = cur.Find(bson.M{"status": Ready}).One(&result)
		if err != nil {
			return result, err
		}
		err = cur.Update(bson.M{"_id": result.ID}, bson.M{"$set": bson.M{"status": StartProcessing}})
		if err != nil {
			return result, err
		}
		// 해당 Item을 반환한다.
		return result, nil
	}
	return result, errors.New("ready상태인 Item이 없습니다")
}

// AddUser 는 데이터베이스에 User를 넣는 함수이다.
// func AddUser(session *mgo.Session, u User) error {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C("users")
// 	n, err := c.Find(bson.M{"id": u.ID}).Count()
// 	if err != nil {
// 		return err
// 	}
// 	if n != 0 {
// 		return errors.New("already exists user ID")
// 	}
// 	err = c.Insert(u)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// RmUser 는 데이터베이스에 User를 삭제하는 함수이다.
// func RmUser(session *mgo.Session, id string) error {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C("users")
// 	err := c.Remove(bson.M{"id": id})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// SetUser 함수는 사용자 정보를 업데이트하는 함수이다.
// func SetUser(session *mgo.Session, u User) error {
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB(*flagDBName).C("users")
// 	num, err := c.Find(bson.M{"id": u.ID}).Count()
// 	if err != nil {
// 		return err
// 	}
// 	if num != 1 {
// 		return errors.New("해당 유저가 존재하지 않습니다")
// 	}
// 	err = c.Update(bson.M{"id": u.ID}, u)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// GetUser 함수는 id를 입력받아서 사용자 정보를 반환한다.
func GetUser(session *mgo.Session, id string) (User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("users")
	u := User{}
	err := c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}
