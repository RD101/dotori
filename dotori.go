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

	flagAdd = flag.Bool("add", false, "add")
	flagRm  = flag.Bool("remove", false, "remove")

	flagAuthor      = flag.String("author", "", "author")
	flagTag         = flag.String("tag", "", "tag")
	flagDescription = flag.String("description", "", "description")
	flagThumbimg    = flag.String("thumbimg", "", "path of thumbnail image")
	flagThumbmov    = flag.String("thumbmov", "", "path of thumbnail mov")
	flagInputpath   = flag.String("inputpath", "", "input path")
	flagOutputpath  = flag.String("outputpath", "", "output path")
	flagType        = flag.String("type", "", "type of asset")
	flagAttributes  = flag.String("attributes", "", "detail info of file") // "key:value,key:value"

	flagDBIP   = flag.String("dbip", "", "DB IP")
	flagDBName = flag.String("dbname", "dotori", "DB name")

	flagHTTPPort = flag.String("http", "", "Web Service Port Number")
)

func main() {
	flag.Parse()
	if *flagAdd {
		i := Item{}

		i.Author = *flagAuthor
		i.Tags = append(i.Tags, *flagTag)
		i.Description = *flagDescription
		i.Thumbimg = *flagThumbimg
		i.Thumbmov = *flagThumbmov
		i.Inputpath = *flagInputpath
		i.Outputpath = *flagOutputpath
		i.Type = *flagType

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

		if *flagDBName != "" {
			if !regexLower.MatchString(*flagDBName) { // 입력받은 dbname이 소문자인지 확인
				log.Fatal(err)
			}
		}

		defer session.Close()
		err = AddItem(session, i)
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
