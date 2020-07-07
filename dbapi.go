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

// AllItems 는 DB에서 전체 아이템 정보를 가져오는 함수입니다.
func AllItems(client *mongo.Client) ([]Item, error) {
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

// UpdateItem 은 컬렉션 이름과 Item을 받아서, Item을 업데이트한다.
func UpdateItem(client *mongo.Client, i Item) error {
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
		querys := []bson.M{}
		//"tag:"가 앞에 붙어있으면 태그에서 검색한다.
		if strings.HasPrefix(word, "tag:") {
			querys = append(querys, bson.M{"tags": strings.TrimPrefix(word, "tag:")})
		} else if strings.HasPrefix(word, "author:") {
			querys = append(querys, bson.M{"author": strings.TrimPrefix(word, "author:")})
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
		querys := []bson.M{}
		//"tag:"가 앞에 붙어있으면 태그에서 검색한다.
		if strings.HasPrefix(word, "tag:") {
			querys = append(querys, bson.M{"tags": strings.TrimPrefix(word, "tag:")})
		} else if strings.HasPrefix(word, "author:") {
			querys = append(querys, bson.M{"author": strings.TrimPrefix(word, "author:")})
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
func GetFileUploadedItem(client *mongo.Client) (Item, error) {
	var result Item
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(*flagDBName).Collection("items")
	filter := bson.M{"status": "fileuploaded"}
	err := collection.FindOne(ctx, filter).Decode(&result)
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

// GetOngoingProcess 는 처리 중인 아이템을 가져온다.
func GetOngoingProcess(client *mongo.Client) ([]Item, error) {
	var results []Item
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"$and": []interface{}{
		// 프로세스가 필요한 아이템 타입
		bson.M{"$or": []interface{}{
			bson.M{"itemtype": "maya"},
			bson.M{"itemtype": "nuke"},
			bson.M{"itemtype": "houdini"},
			bson.M{"itemtype": "blender"},
			bson.M{"itemtype": "footage"},
			bson.M{"itemtype": "alembic"},
			bson.M{"itemtype": "usd"},
		}},
		// 조건: 프로세스가 종료되지 않은 모든 아이템
		bson.M{"status": bson.M{"$ne": "done"}},
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

//GetIncompleteItems 함수는 썸네일 이미지, 썸네일 클립, 데이터 중 하나라도 없는 아이템을 모두 가져온다.
func GetIncompleteItems(client *mongo.Client) ([]Item, error) {
	var results []Item
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"$or": []interface{}{
		// sound, pdf, hwp를 제외한 itemtype의 조건
		bson.M{"$and": []interface{}{
			// 아이템 타입
			bson.M{"$or": []interface{}{
				bson.M{"itemtype": "maya"},
				bson.M{"itemtype": "nuke"},
				bson.M{"itemtype": "houdini"},
				bson.M{"itemtype": "blender"},
				bson.M{"itemtype": "footage"},
				bson.M{"itemtype": "alembic"},
				bson.M{"itemtype": "usd"},
			}},
			// 조건
			bson.M{"$or": []interface{}{
				bson.M{"thumbimguploaded": false},
				bson.M{"thumbclipuploaded": false},
				bson.M{"datauploaded": false},
			}},
		}},
		// sound, pdf, hwp 타입 아이템의 조건
		bson.M{"$and": []interface{}{
			// 아이템 타입
			bson.M{"$or": []interface{}{
				bson.M{"itemtype": "sound"},
				bson.M{"itemtype": "pdf"},
				bson.M{"itemtype": "hwp"},
			}},
			// 조건
			bson.M{"datauploaded": false},
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
