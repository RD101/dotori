package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// TEMPLATES 는 kalena에서 사용하는 템플릿 글로벌 변수이다.
	TEMPLATES = template.New("")

	flagAdd               = flag.Bool("add", false, "add mode")
	flagRm                = flag.Bool("remove", false, "remove mode")
	flagSeek              = flag.Bool("seek", false, "seek mode") // 해당 폴더를 탐색할 때 사용합니다.
	flagSearch            = flag.Bool("search", false, "search mode")
	flagSearchWord        = flag.String("searchword", "", "search word") // DB를 검색할 때 사용합니다.
	flagSearchID          = flag.Bool("searchid", false, "search a item by its id")
	flagGetOngoingProcess = flag.Bool("getongoingprocess", false, "get ongoing process") // 완료되지 않은 프로세스를 가져옵니다.
	flagProcess           = flag.Bool("process", false, "start processing item")         // 프로세스를 실행시킨다

	flagAuthor      = flag.String("author", "", "author")
	flagTag         = flag.String("tag", "", "tag")
	flagDescription = flag.String("description", "", "description")
	flagInputpath   = flag.String("inputpath", "", "input path")
	flagItemType    = flag.String("itemtype", "", "type of asset")
	flagAttributes  = flag.String("attributes", "", "detail info of file") // "key:value,key:value"

	// 서비스에 필요한 인수
	flagDBIP      = flag.String("dbip", "", "DB IP")
	flagDBName    = flag.String("dbname", "dotori", "DB name")
	flagHTTPPort  = flag.String("http", "", "Web Service Port Number")
	flagPagenum   = flag.Int64("pagenum", 9, "maximum number of items in a page")
	flagCookieAge = flag.Int("cookieage", 4, "cookie age (hour)") // MPAA 기준 4시간이다.

	flagItemID = flag.String("itemid", "", "bson ObjectID assigned by mongodb")
)

func main() {
	flag.Parse()
	if *flagSeek {
		items, err := searchSeq(*flagInputpath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(items)
		os.Exit(0)
	} else if *flagSearch {
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		if *flagItemType == "" {
			log.Fatal("itemtype이 빈 문자열입니다")
		}
		items, err := Search(client, *flagItemType, *flagSearchWord)
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range items {
			fmt.Println(item)
		}
	} else if *flagAdd {
		if *flagItemType == "" {
			log.Fatal("itemtype이 빈 문자열입니다")
		}

		i := Item{}
		i.ID = primitive.NewObjectID()
		i.Author = *flagAuthor
		i.Tags = SplitBySpace(*flagTag)
		i.Description = *flagDescription
		i.ItemType = *flagItemType
		i.Attributes = StringToMap(*flagAttributes)

		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		// admin settin에서 rootpath를 가져와서 경로를 생성한다.
		rootpath, err := GetRootPath(client)
		if err != nil {
			log.Fatal(err)
		}
		objIDpath, err := idToPath(i.ID.Hex())
		if err != nil {
			log.Fatal(err)
		}
		i.InputThumbnailImgPath = rootpath + objIDpath + "/originalthumbimg/"
		i.InputThumbnailClipPath = rootpath + objIDpath + "/originalthumbmov/"
		i.OutputThumbnailPngPath = rootpath + objIDpath + "/thumbnail/thumbnail.png"
		i.OutputThumbnailMp4Path = rootpath + objIDpath + "/thumbnail/thumbnail.mp4"
		i.OutputThumbnailOggPath = rootpath + objIDpath + "/thumbnail/thumbnail.ogg"
		i.OutputThumbnailMovPath = rootpath + objIDpath + "/thumbnail/thumbnail.mov"
		i.OutputDataPath = rootpath + objIDpath + "/data/"

		err = i.CheckError()
		if err != nil {
			log.Fatal(err)
		}
		if *flagDBIP != "" {
			if !regexIPv4.MatchString(*flagDBIP) { // 입력받은 DB IP의 형식이 맞는지 확인
				log.Fatal(err)
			}
		}
		if *flagDBName != "" {
			if !regexLower.MatchString(*flagDBName) { // 입력받은 dbname이 소문자인지 확인
				log.Fatal(err)
			}
		}
		err = AddItem(client, i)
		if err != nil {
			log.Print(err)
		}
	} else if *flagRm {
		if *flagItemType == "" {
			log.Fatal("flagType이 빈 문자열 입니다")
		}
		if *flagItemID == "" {
			log.Fatal("id가 빈 문자열 입니다")
		}
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		err = RmItem(client, *flagItemType, *flagItemID)
		if err != nil {
			log.Print(err)
		}
	} else if *flagHTTPPort != "" {
		ip, err := serviceIP()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Service start: http://%s\n", ip)
		webserver()
	} else if *flagSearchID {
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		item, err := SearchItem(client, *flagItemType, *flagItemID)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(item)
	} else if *flagGetOngoingProcess {
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		defer client.Disconnect(ctx)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		items, err := GetOngoingProcess(client)
		fmt.Println(items)
	} else if *flagProcess {
		err := processingItem()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("done")
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

}
