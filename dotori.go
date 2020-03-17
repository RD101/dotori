package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"gopkg.in/mgo.v2"
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
	flagThumbimg    = flag.String("thumbimg", "", "path of thumbnail image")
	flagInputpath   = flag.String("inputpath", "", "input path")
	flagOutputpath  = flag.String("outputpath", "", "output path")
	flagItemType    = flag.String("itemtype", "", "type of asset")
	flagAttributes  = flag.String("attributes", "", "detail info of file") // "key:value,key:value"

	// 서비스에 필요한 인수
	flagDBIP      = flag.String("dbip", "", "DB IP")
	flagDBName    = flag.String("dbname", "dotori", "DB name")
	flagHTTPPort  = flag.String("http", "", "Web Service Port Number")
	flagPagenum   = flag.Int("pagenum", 9, "maximum number of items in a page")
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
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		if *flagItemType == "" {
			log.Fatal("itemtype이 빈 문자열입니다")
		}
		items, err := Search(session, *flagItemType, *flagSearchWord)
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range items {
			fmt.Println(item)
		}
	} else if *flagAdd {
		i := Item{}

		i.Author = *flagAuthor
		i.Tags = SplitBySpace(*flagTag)
		i.Description = *flagDescription
		i.Thumbimg = *flagThumbimg
		i.Outputpath = *flagOutputpath
		i.ItemType = *flagItemType
		i.Attributes = StringToMap(*flagAttributes)

		err := i.CheckError()
		if err != nil {
			log.Fatal(err)
		}
		if *flagDBIP != "" {
			if !regexIPv4.MatchString(*flagDBIP) { // 입력받은 DB IP의 형식이 맞는지 확인
				log.Fatal(err)
			}
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		if *flagDBName != "" {
			if !regexLower.MatchString(*flagDBName) { // 입력받은 dbname이 소문자인지 확인
				log.Fatal(err)
			}
		}
		err = AddItem(session, i)
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
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = RmItem(session, *flagItemType, *flagItemID)
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
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		item, err := SearchItem(session, *flagItemType, *flagItemID)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(item)
	} else if *flagGetOngoingProcess {
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		items, err := GetOngoingProcess(session)
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
