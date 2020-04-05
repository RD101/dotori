package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddItem 은 데이터베이스에 Item을 넣는 함수이다.
func AddItem(client *mongo.Client, i Item) error {
	collection := client.Database(*flagDBName).Collection(i.ItemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}

// GetItem 은 데이터베이스에 Item을 가지고 오는 함수이다.
func GetItem(client *mongo.Client, itemType, id string) (Item, error) {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var result Item
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetAdminSetting 은 관리자 셋팅값을 가지고 온다.
func GetAdminSetting(client *mongo.Client) (Adminsetting, error) {
	//monotonic이 필수적인가? 필수적이라면 대응하는 기능은 무엇인가
	//session.SetMode(mgo.Monotonic, true)
	collection := client.Database(*flagDBName).Collection("setting.admin")
	var result Adminsetting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{"id": "setting.admin"}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNilDocument { // document가 존재하지 않는 경우, 즉 adminsetting이 없는 경우
			return Adminsetting{}, nil
		}
		return Adminsetting{}, err
	}
	return result, nil
}

// RmItem 는 컬렉션 이름과 id를 받아서, 해당 컬렉션에서 id가 일치하는 Item을 삭제한다.
func RmItem(client *mongo.Client, itemType, id string) error {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

// AllItems 는 DB에서 전체 아이템 정보를 가져오는 함수입니다.
func AllItems(client *mongo.Client, itemType string) ([]Item, error) {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var results []Item
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// UpdateItem 은 컬렉션 이름과 Item을 받아서, Item을 업데이트한다.
func UpdateItem(client *mongo.Client, itemType string, item Item) error {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// Search 는 itemType, words를 입력받아 해당 아이템을 검색한다.
// http_restapi.go에서 사용중
func Search(client *mongo.Client, itemType string, words string) ([]Item, error) {
	var results []Item
	//검색어가 존재하지 않으면 빈 결과를 반환한다.
	if words == "" {
		return results, nil
	}
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

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
			querys = append(querys, bson.M{"attributes." + key: primitive.Regex{Pattern: value, Options: "i"}})
		} else {
			querys = append(querys, bson.M{"author": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"tags": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"description": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"type": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"inputpath": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"outputpath": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"createtime": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"updatetime": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"attributes.*": word})
		}
		wordsQueries = append(wordsQueries, bson.M{"$or": querys})
	}
	// 사용률이 많은 소스가 위로 출력되도록 한다.
	q := bson.M{"$and": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.
	opts := options.Find()
	opts.SetSort(bson.D{{"usingrate", -1}})
	cursor, err := collection.Find(ctx, q, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchPage 는 itemType, words, 해당 page를 입력받아 해당 아이템을 검색한다. 검색된 아이템과 그 개수를 반환한다.
// http.go에서 사용중
func SearchPage(client *mongo.Client, itemType string, words string, page, limitnum int64) (int64, int64, []Item, error) {
	var results []Item
	//검색어가 존재하지 않으면 빈 결과를 반환한다.
	if words == "" {
		return 0, 0, results, nil
	}
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

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
			querys = append(querys, bson.M{"attributes." + key: primitive.Regex{Pattern: value, Options: "i"}})
		} else {
			querys = append(querys, bson.M{"author": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"tags": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"description": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"type": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"inputpath": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"outputpath": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"createtime": primitive.Regex{Pattern: word, Options: "i"}})
			querys = append(querys, bson.M{"updatetime": primitive.Regex{Pattern: word, Options: "i"}})
		}
		wordsQueries = append(wordsQueries, bson.M{"$or": querys})
	}
	// 사용률이 많은 소스가 위로 출력되도록 한다.
	q := bson.M{"$and": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.
	opts := options.Find()
	opts.SetSort(bson.D{{"usingrate", -1}})
	opts.SetSkip(int64((page - 1) * limitnum))
	opts.SetLimit(int64(limitnum))
	cursor, err := collection.Find(ctx, q, opts)
	if err != nil {
		return 0, 0, nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return 0, 0, nil, err
	}
	totalNum, err := collection.CountDocuments(ctx, q)
	if err != nil {
		return 0, 0, nil, err
	}
	return TotalPage(totalNum, limitnum), totalNum, results, nil
}

// SearchItem 은 컬렉션 이름(itemType)과 id를 받아서, 해당 컬렉션에서 id가 일치하는 item을 검색, 반환한다.
func SearchItem(client *mongo.Client, itemType, id string) (Item, error) {
	collection := client.Database(*flagDBName).Collection(itemType)
	var result Item
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SearchTags 는 itemType, tag를 입력받아 tag의 값이 일치하면 반환하는 함수입니다.
func SearchTags(client *mongo.Client, itemType string, tag string) ([]Item, error) {
	collection := client.Database(*flagDBName).Collection(itemType)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var results []Item
	cursor, err := collection.Find(ctx, bson.M{"tags": tag})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// AddUser 는 데이터베이스에 User를 넣는 함수이다.
func AddUser(client *mongo.Client, u User) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	n, err := collection.CountDocuments(ctx, bson.M{"id": u.ID})
	if err != nil {
		return err
	}
	if n != 0 {
		return errors.New("already exists user ID")
	}
	_, err = collection.InsertOne(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

// RmUser 는 데이터베이스에 User를 삭제하는 함수이다.
func RmUser(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// SetUser 함수는 사용자 정보를 업데이트하는 함수이다.
func SetUser(client *mongo.Client, u User) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	n, err := collection.CountDocuments(ctx, bson.M{"id": u.ID})
	if err != nil {
		return err
	}
	if n != 0 {
		return errors.New("already exists user ID")
	}
	_, err = collection.UpdateOne(ctx, bson.M{"id": u.ID}, u)
	if err != nil {
		return err
	}
	return nil
}

// GetUser 함수는 id를 입력받아서 사용자 정보를 반환한다.
func GetUser(client *mongo.Client, id string) (User, error) {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	u := User{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// SetAdminSetting 은 입력받은 어드민셋팅으로 업데이트한다.
func SetAdminSetting(client *mongo.Client, a Adminsetting) error {
	collection := client.Database(*flagDBName).Collection("setting.admin")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	n, err := collection.CountDocuments(ctx, bson.M{"id": "setting.admin"})
	if err != nil {
		return err
	}
	a.ID = "setting.admin"
	if n == 0 {
		_, err = collection.InsertOne(ctx, a)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = collection.UpdateOne(ctx, bson.M{"id": "setting.admin"}, a)
	if err != nil {
		return err
	}
	return nil
}

// GetOngoingProcess 는 처리 중인 아이템을 가져온다.
func GetOngoingProcess(client *mongo.Client) ([]Item, error) {
	var results []Item
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//콜렉션 리스트를 가져온다.
	collections, err := client.Database(*flagDBName).ListCollectionNames(ctx, bson.M{})
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
		cursor, err := client.Database(*flagDBName).Collection(c).Find(ctx, bson.M{"status": bson.M{"$ne": Done}})
		if err != nil {
			return results, err
		}
		err = cursor.All(ctx, &results)
		if err != nil {
			return results, err
		}
		results = append(results, items...)
	}
	return results, nil
}
