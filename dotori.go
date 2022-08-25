package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/user"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// TEMPLATES 는 dotori에서 사용하는 템플릿 글로벌 변수이다.
	TEMPLATES = template.New("")

	flagAdd               = flag.Bool("add", false, "add mode")
	flagRm                = flag.Bool("remove", false, "remove mode")
	flagSeek              = flag.Bool("seek", false, "seek mode") // 해당 폴더를 탐색할 때 사용합니다.
	flagExt               = flag.String("ext", "", "extenstion")  // 확장자
	flagSearch            = flag.Bool("search", false, "search mode")
	flagSearchWord        = flag.String("searchword", "", "search word") // DB를 검색할 때 사용합니다.
	flagSearchID          = flag.Bool("searchid", false, "search a item by its id")
	flagGetOngoingProcess = flag.Bool("getongoingprocess", false, "get ongoing process")  // 완료되지 않은 프로세스를 가져옵니다.
	flagProcess           = flag.Bool("process", false, "start processing item")          // 프로세스를 실행시킨다
	flagDebug             = flag.Bool("debug", false, "debug mode")                       // debug모드
	flagAccesslevel       = flag.String("accesslevel", "default", "access level of user") // 사용자의 accesslevel을 지정합니다. admin, manager, default

	flagAuthor             = flag.String("author", "", "author")
	flagTitle              = flag.String("title", "", "title")
	flagTag                = flag.String("tag", "", "tag")
	flagDescription        = flag.String("description", "", "description")
	flagInputThumbImgPath  = flag.String("inputthumbimgpath", "", "input path of thumbnail image")
	flagInputThumbClipPath = flag.String("inputthumbclippath", "", "input path of thumbnail clip")
	flagInputDataPath      = flag.String("inputdatapath", "", "input path of data")
	flagItemType           = flag.String("itemtype", "", "type of asset")
	flagAttributes         = flag.String("attributes", "", "detail info of file") // "key:value,key:value"
	flagUserID             = flag.String("userid", "", "ID of user")
	flagFPS                = flag.String("fps", "", "frame per second")
	flagInColorspace       = flag.String("incolorspace", "", "in color space")
	flagOutColorspace      = flag.String("outcolorspace", "", "out color space")
	flagPremultiply        = flag.Bool("premultiply", false, "premultiply")

	// 서비스에 필요한 인수
	flagMongoDBURI    = flag.String("mongodburi", "mongodb://localhost:27017", "mongoDB URI ex)mongodb://localhost:27017")
	flagDBName        = flag.String("dbname", "dotori", "DB name")
	flagHTTPPort      = flag.String("http", "", "Web Service Port Number")
	flagPagenum       = flag.Int64("pagenum", 9, "maximum number of items in a page")
	flagCookieAge     = flag.Int("cookieage", 4, "cookie age (hour)") // MPAA 기준 4시간이다.
	flagMaxProcessNum = flag.Int("maxprocessnum", 4, "maximum number of process")
	flagCertFullchain = flag.String("certfullchain", "", "certification fullchain path")
	flagCertPrivkey   = flag.String("certprivkey", "", "certification privkey path")

	flagItemID = flag.String("itemid", "", "bson ObjectID assigned by mongodb")

	// SHA1VER 값은 git 커밋 로그 이다.
	SHA1VER = ""
	// BUILDTIME 값은 빌드시간 변수이다.
	BUILDTIME = ""
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("dotori: ")
	flag.Parse()

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if *flagSeek && *flagExt == "" {
		items, err := searchSeqAndClip(*flagInputDataPath)
		if err != nil {
			log.Fatal(err)
		}
		data, err := json.Marshal(items)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
		os.Exit(0)
	} else if *flagSeek && *flagExt != "" {
		paths, err := searchExt(*flagInputDataPath, *flagExt)
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range paths {
			fmt.Println(path)
		}
		os.Exit(0)
	} else if *flagSearch {
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
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
		if *flagDBName != "" {
			if !regexLower.MatchString(*flagDBName) { // 입력받은 dbname이 소문자인지 확인
				log.Fatal(err)
			}
		}
		switch *flagItemType {
		case "maya":
			addMayaItemCmd()
		case "houdini":
			addHoudiniItemCmd()
		case "blender":
			addBlenderItemCmd()
		case "clip":
			addClipItemCmd()
		case "footage":
			addFootageItemCmd()
		case "nuke":
			addNukeItemCmd()
		case "alembic":
			addAlembicItemCmd()
		case "usd":
			addUSDItemCmd()
		case "unreal":
			addUnrealItemCmd()
		case "hwp":
			addHwpItemCmd()
		case "pdf":
			addPdfItemCmd()
		case "texture":
			addTextureItemCmd()
		case "sound":
			addSoundItemCmd()
		case "openvdb":
			addOpenVDBItemCmd()
		case "modo":
			addModoItemCmd()
		case "katana":
			addKatanaItemCmd()
		case "ppt":
			addPptItemCmd()
		case "ies":
			addIesItemCmd()
		case "lut":
			addLutItemCmd()
		case "hdri":
			addHdriItemCmd()
		case "fusion360":
			addFusion360ItemCmd()
		case "max":
			addMaxItemCmd()
		default:
			log.Fatal("command를 지원하지 않는 아이템타입입니다.")
		}
	} else if *flagRm {
		if user.Username != "root" {
			log.Fatal(errors.New("item을 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		rmItemCmd()
	} else if *flagHTTPPort != "" {
		ip, err := serviceIP()
		if err != nil {
			log.Fatal(err)
		}
		// 프로세스 연산을 실행한다.
		// webserver와 같이 실행해야하기 때문에 go를 붙혀서 실행한다.
		// go 명령어가 없다면, webserver() 함수가 실행되지 않는다.
		go ProcessMain()
		// 웹서버 실행
		fmt.Printf("Service start: http://%s\n", ip)
		webserver()

	} else if *flagSearchID {
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		item, err := SearchItem(client, *flagItemID)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(item)
	} else if *flagGetOngoingProcess {
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}

		items, err := GetUndoneItem(client)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(items)
	} else if *flagProcess {
		ProcessMain()
		fmt.Println("done")
	} else if *flagAccesslevel != "" && *flagUserID != "" {
		if user.Username != "root" {
			log.Fatal(errors.New("사용자의 레벨을 수정하기 위해서는 root 권한이 필요합니다"))
		}
		//mongoDB client 연결
		client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}
		// userID를 이용해서 사용자정보를 가져온다.
		u, err := GetUser(client, *flagUserID)
		if err != nil {
			log.Fatal(err)
		}
		u.AccessLevel = *flagAccesslevel
		//수정된 사용자 정보를 DB에 업데이트한다.
		u.CreateToken()
		err = SetUser(client, u)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
