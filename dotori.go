package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagAdd      = flag.Bool("add", false, "add")
	flagID       = flag.String("id", "", "id")
	flagTag      = flag.String("tag", "", "tag")
	flagThumbimg = flag.String("thumbimg", "", "path of thumbnail image")
)

func main() {
	flag.Parse()
	if *flagAdd {
		i := Item{}
		i.ID = *flagID
		i.Tags = append(i.Tags, *flagTag)
		i.Thumbimg = *flagThumbimg

		fmt.Println(i)
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

}
