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

	flagID          = flag.String("id", "", "id")
	flagAuthor		  = flag.String("author", "", "author")
	flagTag         = flag.String("tag", "", "tag")
	flagDescription = flag.String("description"), "", "description"
	flagThumbimg    = flag.String("thumbimg", "", "path of thumbnail image")
	flagThumbmov    = flag.String("thumbmov", "", "path of thumbnail mov")
	flagInputpath   = flag.String("inputpath", "", "input path")
	flagOutputpath  = flag.String("outputpath", "", "output path")
	flagType        = flag.String("type", "", "type of asset")
	flagStatus      = flag.String("status", "", "status of asset")
	flagLog				  = flag.String("log", "", "log")
	flagCreatetime  = flag.String("createtime", "", "created time")
	flagUpdatetime  = flag.String("updatetime", "", "updated time")
	flagUsingRate		= flag.Int64("usingrate", 0 , "using rate")
	flagStorage			= flag.String("storage", "", "info of storage")
	flagAttributes	= flag.String("attributes", "", "detail info of file")

	flagDBIP = flag.String("dbip", "", "DB IP")
)

func main() {
	flag.Parse()
	if *flagAdd {
		i := Item{}

		i.ID = *flagID
		i.Author = *flagAuthor
		i.Tags = append(i.Tags, *flagTag)
		i.Description = *flagDescription
		i.Thumbimg = *flagThumbimg
		i.Thumbmov = *flagThumbmov
		i.Inputpath = *flagInputpath
		i.Outputpath = *flagOutputpath
		i.Type = *flagType
		i.Status = *flagStatus
		i.Log = *flagLog
		i.CreateTime = *flagCreatetime
		i.Updatetime = *flagUpdatetime
		i.UsingRate = *flagUsingRate
		i.Storage = *flagStorage
		i.Attributes = *flagAttributes

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
