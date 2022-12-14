package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// AddItem 은 데이터베이스에 Item을 넣는 함수이다.
func AddItem(client *mongo.Client, i Item) error {
	i.CreateTime = time.Now().Format(time.RFC3339)
	i.Updatetime = i.CreateTime
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}

// GetItem 은 데이터베이스에 Item을 가지고 오는 함수이다.
func GetItem(client *mongo.Client, id string) (Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"id": "setting.admin"}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments { // document가 존재하지 않는 경우, 즉 adminsetting이 없는 경우
			return Adminsetting{}, nil
		}
		return Adminsetting{}, err
	}
	return result, nil
}

// RmItem 는 id를 받아서, 해당 컬렉션에서 id가 일치하는 Item을 삭제한다.
func RmItem(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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

// RmFavoriteItem 는 id를 받아서 해당 id를 즐겨찾기 하고 있는 User가 있다면 favoriteassetids 필드에서 삭제한다.
func RmFavoriteItem(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"favoriteassetids": id}
	update := bson.M{"$pull": bson.M{"favoriteassetids": id}}
	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetAllItems 는 DB에서 전체 아이템 정보를 가져오는 함수입니다.
func GetAllItems(client *mongo.Client) ([]Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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

// GetRecentlyCreatedItems 는 DB에서 최근생성 된 num건의 아이템 정보를 가져오는 함수입니다.
func GetRecentlyCreatedItems(client *mongo.Client, limitnum int64, page int64) ([]Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Item
	opts := options.Find()
	opts.SetSort(bson.M{"createtime": -1})
	opts.SetSkip(int64((page - 1) * limitnum))
	opts.SetLimit(int64(limitnum))
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetTopUsingItems 는 DB에서 Using숫자가 높은 순서, num건의 아이템 정보를 가져오는 함수입니다.
func GetTopUsingItems(client *mongo.Client, limitnum int64, page int64) ([]Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Item
	opts := options.Find()
	opts.SetSort(bson.M{"usingrate": -1})
	opts.SetSkip(int64((page - 1) * limitnum))
	opts.SetLimit(int64(limitnum))
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetAllItemsNum 는 DB에서 전체 아이템의 개수 정보를 가져오는 함수입니다.
func GetAllItemsNum(client *mongo.Client) (int64, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	n, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return n, err
	}
	return n, nil
}

// SetItem 은 컬렉션 이름과 Item을 받아서, Item을 업데이트한다.
func SetItem(client *mongo.Client, i Item) error {
	i.Updatetime = time.Now().Format(time.RFC3339)
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": i.ID},
		bson.D{{Key: "$set", Value: i}},
	)
	if err != nil {
		return err
	}
	return nil
}

// Search 는 words를 입력받아 해당 아이템을 검색한다.
// http_restapi.go에서 사용중
func Search(client *mongo.Client, itemType, words string) ([]Item, error) {
	var results []Item
	//검색어가 존재하지 않으면 빈 결과를 반환한다.
	if words == "" {
		return results, nil
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wordsQueries := []bson.M{}
	if itemType != "" {
		querys := []bson.M{}
		querys = append(querys, bson.M{"itemtype": itemType})
		wordsQueries = append(wordsQueries, bson.M{"$and": querys})
	}
	for _, word := range strings.Split(words, " ") {
		if word == "" {
			continue
		}
		// 쉼표가 존재한다면 쉼표를 제거합니다.
		word = strings.Trim(word, ",")
		querys := []bson.M{}
		//"tag:"가 앞에 붙어있으면 태그에서 검색한다.
		if strings.HasPrefix(word, "tag:") {
			querys = append(querys, bson.M{"tags": strings.TrimPrefix(word, "tag:")})
		} else if strings.HasPrefix(word, "categories:") {
			querys = append(querys, bson.M{"categories": strings.TrimPrefix(word, "categories:")})
		} else if strings.HasPrefix(word, "category:") {
			querys = append(querys, bson.M{"categories": strings.TrimPrefix(word, "category:")})
		} else if strings.HasPrefix(word, "author:") {
			querys = append(querys, bson.M{"author": strings.TrimPrefix(word, "author:")})
		} else if strings.HasPrefix(word, "title:") {
			querys = append(querys, bson.M{"title": strings.TrimPrefix(word, "title:")})
		} else if strings.Contains(word, ":") {
			key := strings.Split(word, ":")[0]
			value := strings.Split(word, ":")[1]
			querys = append(querys, bson.M{"attributes." + key: primitive.Regex{Pattern: value, Options: "i"}})
		} else {
			querys = append(querys, bson.M{"title": primitive.Regex{Pattern: word, Options: "i"}})
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
	// word에 "," 표시가 존재한다면 or 검색을 진행합니다.
	if strings.Contains(words, ",") {
		q = bson.M{"$or": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.
	}
	opts := options.Find()
	opts.SetSort(bson.M{"usingrate": -1})
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

// SearchPage 는 words, 해당 page를 입력받아 해당 아이템을 검색한다. 검색된 아이템과 그 개수를 반환한다.
// http.go에서 사용중
func SearchPage(client *mongo.Client, itemType, words string, page, limitnum int64) (int64, int64, []Item, error) {
	var results []Item
	//검색어가 존재하지 않으면 빈 결과를 반환한다.
	if words == "" {
		return 0, 0, results, nil
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wordsQueries := []bson.M{}

	if itemType != "" {
		querys := []bson.M{}
		querys = append(querys, bson.M{"itemtype": itemType})
		wordsQueries = append(wordsQueries, bson.M{"$and": querys})
	}

	for _, word := range strings.Split(words, " ") {
		if word == "" {
			continue
		}
		// 쉼표가 존재한다면 쉼표를 제거합니다.
		word = strings.Trim(word, ",")
		querys := []bson.M{}
		//"tag:"가 앞에 붙어있으면 태그에서 검색한다.
		if strings.HasPrefix(word, "tag:") {
			querys = append(querys, bson.M{"tags": strings.TrimPrefix(word, "tag:")})
		} else if strings.HasPrefix(word, "categories:") {
			querys = append(querys, bson.M{"categories": strings.TrimPrefix(word, "categories:")})
		} else if strings.HasPrefix(word, "category:") {
			querys = append(querys, bson.M{"categories": strings.TrimPrefix(word, "category:")})
		} else if strings.HasPrefix(word, "author:") {
			querys = append(querys, bson.M{"author": strings.TrimPrefix(word, "author:")})
		} else if strings.HasPrefix(word, "title:") {
			querys = append(querys, bson.M{"title": strings.TrimPrefix(word, "title:")})
		} else if strings.Contains(word, ":") {
			key := strings.Split(word, ":")[0]
			value := strings.Split(word, ":")[1]
			querys = append(querys, bson.M{"attributes." + key: primitive.Regex{Pattern: value, Options: "i"}})
		} else {
			switch strings.ToLower(word) {
			case "all":
				querys = append(querys, bson.M{})
			default:
				querys = append(querys, bson.M{"title": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"author": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"tags": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"description": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"type": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"inputpath": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"outputpath": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"createtime": primitive.Regex{Pattern: word, Options: "i"}})
				querys = append(querys, bson.M{"updatetime": primitive.Regex{Pattern: word, Options: "i"}})
			}
		}
		wordsQueries = append(wordsQueries, bson.M{"$or": querys})
	}
	// 사용률이 많은 소스가 위로 출력되도록 한다.

	q := bson.M{"$and": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.

	// word에 "," 표시가 존재한다면 or 검색을 진행합니다.
	if strings.Contains(words, ",") {
		q = bson.M{"$or": wordsQueries} // 최종 쿼리는 BSON type 오브젝트가 되어야 한다.
	}
	opts := options.Find()
	opts.SetSort(bson.M{"usingrate": -1})
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

// SearchItem 은 id를 받아서, 해당 컬렉션에서 id가 일치하는 item을 검색, 반환한다.
func SearchItem(client *mongo.Client, id string) (Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	var result Item
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SearchTags 는 tag를 입력받아 tag의 값이 일치하면 반환하는 함수입니다.
func SearchTags(client *mongo.Client, tag string) ([]Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// SetUser 함수는 사용자 정보를 업데이트하는 함수이다.
func SetUser(client *mongo.Client, u User) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"id": u.ID},
		bson.D{{Key: "$set", Value: u}},
	)
	if err != nil {
		return err
	}
	return nil
}

// GetUser 함수는 id를 입력받아서 사용자 정보를 반환한다.
func GetUser(client *mongo.Client, id string) (User, error) {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u := User{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetAllUsers 함수는 DB에서 전체 사용자 정보를 가져오는 함수이다.
func GetAllUsers(client *mongo.Client) ([]User, error) {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []User
	opts := options.Find()
	opts.SetSort(bson.M{"id": 1})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, err
}

// SetAdminSetting 은 입력받은 어드민셋팅으로 업데이트한다.
func SetAdminSetting(client *mongo.Client, a Adminsetting) error {
	collection := client.Database(*flagDBName).Collection("setting.admin")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"id": "setting.admin"},
		bson.D{{Key: "$set", Value: a}},
	)
	if err != nil {
		return err
	}
	return nil
}

// GetFileUploadedItem 은 DB에서 FileUploaded상태인 Item을 하나 가져온다.
func GetFileUploadedItem() (Item, error) {
	var result Item
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		return result, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return result, err
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return result, err
	}
	collection := client.Database(*flagDBName).Collection("items")
	filter := bson.M{"status": "fileuploaded"}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetFileUploadedItemsNum 함수는 fileuploaded 상태 갯수가 몇개인지 체크한다,
func GetFileUploadedItemsNum(client *mongo.Client) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(*flagDBName).Collection("items")
	filter := bson.M{"status": "fileuploaded"}
	n, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// GetUndoneItem 는 status가 done이 아닌 모든 아이템을 가져온다.
func GetUndoneItem(client *mongo.Client) ([]Item, error) {
	var results []Item
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 조건: 프로세스가 종료되지 않은 모든 아이템
	filter := bson.M{"status": bson.M{"$ne": "done"}}
	cursor, err := client.Database(*flagDBName).Collection("items").Find(ctx, filter)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

//SetStatus 함수는 인수로 받은 item의 Status를 status로 바꾼다
func SetStatus(client *mongo.Client, item Item, status string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// item의 Status를 업데이트 한다.
	filter := bson.M{"_id": item.ID}
	update := bson.M{
		"$set": bson.M{"status": status},
	}
	err := collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return err
	}
	return nil
}

//SetStatusAndGetItem 함수는 인수로 받은 item의 Status를 update 하고 update된 item을 return 한다
func SetStatusAndGetItem(item Item, status string) (Item, error) {
	var result Item

	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		return result, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return result, err
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return result, err
	}

	// item의 Status를 업데이트 한다.
	filter := bson.M{"_id": item.ID}
	update := bson.M{
		"$set": bson.M{"status": status},
	}
	collection := client.Database(*flagDBName).Collection("items")
	option := *options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, filter, update, &option).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//SetErrStatus 함수는 인수로 받은 item의 Status를 error status로 바꾼다
func SetErrStatus(client *mongo.Client, id, errmsg string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// item의 Status를 업데이트 한다.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"status": "error",
		},
		"$push": bson.M{"logs": errmsg},
	}
	err = collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetProcessingItemNum 함수는 현재 연산 중인 아이템의 개수를 구한다.
func GetProcessingItemNum(client *mongo.Client) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"$and": []interface{}{
		bson.M{"status": bson.M{"$ne": "ready"}},
		bson.M{"status": bson.M{"$ne": "done"}},
	}}
	collection := client.Database(*flagDBName).Collection("items")
	n, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// SetLog 함수는 id와 로그메세지를 받아서 해당 item에 로그메세지를 더한다.
func SetLog(client *mongo.Client, id, msg string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// 로그를 업데이트 한다.
	filter := bson.M{"_id": objID}
	update := bson.M{"$push": bson.M{"logs": msg}}
	err = collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return err
	}
	return nil
}

//GetIncompleteItems 함수는 프로세스 이후 필요한 정보(썸네일 이미지, 썸네일 클립, 데이터 등)가 없는 아이템을 가지고 온다.
func GetIncompleteItems(client *mongo.Client) ([]Item, error) {
	var results []Item
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"$or": []interface{}{
		// maya, max, fusion360, nuke, houdini, blender, footage, alembic, usd, modo, katana, openvdb 아이템타입의 조건
		// 필요 조건: 썸네일 이미지, 썸네일 클립, 데이터
		bson.M{"$and": []interface{}{
			// 아이템 타입
			bson.M{"$or": []interface{}{
				bson.M{"itemtype": "maya"},
				bson.M{"itemtype": "max"},
				bson.M{"itemtype": "fusion360"},
				bson.M{"itemtype": "nuke"},
				bson.M{"itemtype": "houdini"},
				bson.M{"itemtype": "blender"},
				bson.M{"itemtype": "footage"},
				bson.M{"itemtype": "alembic"},
				bson.M{"itemtype": "usd"},
				bson.M{"itemtype": "modo"},
				bson.M{"itemtype": "katana"},
				bson.M{"itemtype": "openvdb"},
			}},
			// 조건
			bson.M{"$or": []interface{}{
				bson.M{"thumbimguploaded": false},
				bson.M{"thumbclipuploaded": false},
				bson.M{"datauploaded": false},
			}},
		}},
		// sound, pdf, hwp, texture, clip, ies, ppt, unreal 타입 아이템의 조건
		// 필요 조건: 데이터
		bson.M{"$and": []interface{}{
			// 아이템 타입
			bson.M{"$or": []interface{}{
				bson.M{"itemtype": "sound"},
				bson.M{"itemtype": "pdf"},
				bson.M{"itemtype": "hwp"},
				bson.M{"itemtype": "texture"},
				bson.M{"itemtype": "lut"},
				bson.M{"itemtype": "clip"},
				bson.M{"itemtype": "ies"},
				bson.M{"itemtype": "ppt"},
				bson.M{"itemtype": "unreal"},
			}},
			// 조건
			bson.M{"datauploaded": false},
		}},
		// lut, hdri 타입 아이템의 조건
		// 필요 조건: 썸네일 이미지, 데이터
		bson.M{"$and": []interface{}{
			// 아이템 타입
			bson.M{"$or": []interface{}{
				bson.M{"itemtype": "lut"},
				bson.M{"itemtype": "hdri"},
			}},
			// 조건
			bson.M{"$or": []interface{}{
				bson.M{"thumbimguploaded": false},
				bson.M{"datauploaded": false},
			}},
		}},
	}}

	cursor, err := client.Database(*flagDBName).Collection("items").Find(ctx, filter)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetUsingRate 함수는 id를 받아서 해당 아이템의 UsingRate을 가져온다.
func GetUsingRate(client *mongo.Client, id string) (int64, error) {
	var item Item
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"_id": objID}
	err = collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		return 0, err
	}
	return item.UsingRate, nil
}

// UpdateUsingRate 함수는 id를 받아서 해당 아이템의 usingrate을 1만큼 올린다.
func UpdateUsingRate(client *mongo.Client, id string) (int64, error) {
	var item Item
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"usingrate": 1}}
	option := *options.FindOneAndUpdate().SetReturnDocument(options.After) // option은 4.0부터 추가됨.
	err = collection.FindOneAndUpdate(ctx, filter, update, &option).Decode(&item)
	if err != nil {
		return 0, err
	}
	return item.UsingRate, nil
}

//SetThumbImgUploaded 함수는 인수로 받은 item의 ThumbImgUploaded 값을 바꾼다.
func SetThumbImgUploaded(client *mongo.Client, item Item, status bool) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// item의 Status를 업데이트 한다.
	filter := bson.M{"_id": item.ID}
	update := bson.M{
		"$set": bson.M{"thumbimguploaded": status},
	}
	err := collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return err
	}
	return nil
}

//SetThumbClipUploaded 함수는 인수로 받은 item의 ThumbClipUploaded 값을 바꾼다.
func SetThumbClipUploaded(client *mongo.Client, item Item, status bool) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// item의 Status를 업데이트 한다.
	filter := bson.M{"_id": item.ID}
	update := bson.M{
		"$set": bson.M{"thumbclipuploaded": status},
	}
	err := collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetTags 는 DB에서 전체 태그 정보를 가져오는 함수입니다.
func GetTags(client *mongo.Client) ([]string, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []string
	values, err := collection.Distinct(ctx, "tags", bson.D{})
	if err != nil {
		return results, err
	}
	for _, value := range values {
		results = append(results, fmt.Sprintf("%v", value))
	}
	sort.Strings(results)
	return results, nil
}

func addCategory(client *mongo.Client, c Category) error {
	err := c.CheckError()
	if err != nil {
		return err
	}
	collection := client.Database(*flagDBName).Collection("category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"name": c.Name, "parentid": c.ParentID})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름을 가진 데이터가 있습니다")
	}
	_, err = collection.InsertOne(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

func GetCategory(client *mongo.Client, id string) (Category, error) {
	collection := client.Database(*flagDBName).Collection("category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result Category
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

func GetRootCategories(client *mongo.Client) ([]Category, error) {
	collection := client.Database(*flagDBName).Collection("category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Category
	opts := options.Find()
	opts.SetSort(bson.M{"name": 1})
	cursor, err := collection.Find(ctx, bson.M{"parentid": ""}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func GetSubCategories(client *mongo.Client, parentid string) ([]Category, error) {
	collection := client.Database(*flagDBName).Collection("category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Category
	opts := options.Find()
	opts.SetSort(bson.M{"name": 1})
	cursor, err := collection.Find(ctx, bson.M{"parentid": parentid}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func SetCategory(client *mongo.Client, c Category) error {
	collection := client.Database(*flagDBName).Collection("category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": c.ID},
		bson.D{{Key: "$set", Value: c}},
	)
	if err != nil {
		return err
	}
	return nil
}

func RmCategory(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// 부모를 삭제한다.
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	// 자식들을 지운다.
	_, err = collection.DeleteMany(ctx, bson.M{"parentid": id})
	if err != nil {
		return err
	}
	return nil
}
