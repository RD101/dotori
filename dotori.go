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

	flagAdd    = flag.Bool("add", false, "add mode")
	flagRm     = flag.Bool("remove", false, "remove mode")
	flagSeek   = flag.Bool("seek", false, "seek mode")    // 해당 폴더를 탐색할 때 사용합니다.
	flagSearch = flag.String("search", "", "search mode") // DB를 검색할 때 사용합니다.

	flagAuthor      = flag.String("author", "", "author")
	flagTag         = flag.String("tag", "", "tag")
	flagDescription = flag.String("description", "", "description")
	flagThumbimg    = flag.String("thumbimg", "", "path of thumbnail image")
	flagThumbmov    = flag.String("thumbmov", "", "path of thumbnail mov")
	flagInputpath   = flag.String("inputpath", "", "input path")
	flagOutputpath  = flag.String("outputpath", "", "output path")
	flagType        = flag.String("type", "", "type of asset")
	flagAttributes  = flag.String("attributes", "", "detail info of file") // "key:value,key:value"

	// 서비스에 필요한 인수
	flagDBIP     = flag.String("dbip", "", "DB IP")
	flagDBName   = flag.String("dbname", "dotori", "DB name")
	flagHTTPPort = flag.String("http", "", "Web Service Port Number")
	flagPagenum  = flag.Int("pagenum", 9, "maximum number of items in a page")

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
	} else if *flagSearch != "" {
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		items, err := Search(session, *flagType, *flagSearch)
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range items {
			fmt.Println(item)
		}
	} else if *flagAdd {
		i := Item{}

		i.Author = *flagAuthor
		i.Tags = append(i.Tags, *flagTag)
		i.Description = *flagDescription
		i.Thumbimg = *flagThumbimg
		i.Thumbmov = *flagThumbmov
		i.Inputpath = *flagInputpath
		i.Outputpath = *flagOutputpath
		i.Type = *flagType
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
		if *flagType == "" {
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
		err = RmItem(session, *flagType, *flagItemID)
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
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

}
