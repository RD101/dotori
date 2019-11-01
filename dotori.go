package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var (
	flagAdd = flag.Bool("add", false, "add")
	flagRm  = flag.Bool("remove", false, "remove")

	flagID         = flag.String("id", "", "id")
	flagTag        = flag.String("tag", "", "tag")
	flagThumbimg   = flag.String("thumbimg", "", "path of thumbnail image")
	flagThumbmov   = flag.String("thumbmov", "", "path of thumbnail mov")
	flagInputpath  = flag.String("inputpath", "", "input path")
	flagOutputpath = flag.String("outputpath", "", "output path")
	flagType       = flag.String("type", "", "type of asset")
	flagStatus     = flag.String("status", "", "status of asset")
	flagUpdatetime = flag.String("updatetime", "", "updated time")
	flagCreatetime = flag.String("createtime", "", "created time")

	flagDBIP = flag.String("dbip", "", "DB IP")
)

func main() {
	flag.Parse()
	if *flagAdd {
		i := Item{}

		i.ID = *flagID
		i.Tags = append(i.Tags, *flagTag)
		i.Thumbimg = *flagThumbimg
		i.Thumbmov = *flagThumbmov
		i.Inputpath = *flagInputpath
		i.Outputpath = *flagOutputpath
		i.Type = *flagType
		i.Status = *flagStatus
		i.Updatetime = *flagUpdatetime
		i.CreateTime = *flagCreatetime

		err := i.CheckError()
		if err != nil {
			log.Fatal(err)
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = AddItem(session, i)
		if err != nil {
			log.Print(err)
		}
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

}
