package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

// assets 폴더의 하위 모든 파일을 ../assets_vfsdata.go 파일로 만드는 코드이다.
func main() {
	var fs http.FileSystem = http.Dir("assets")
	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}
}
