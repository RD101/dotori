package main

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/mgo.v2"
)

func processingItem() error {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()

	item, err := GetReadyItem(session)
	if err != nil {
		return err
	}
	// Status : 복사중
	item.Status = Copying
	src := item.Inputpath
	dst := item.Outputpath
	err = copyDir(src, dst)
	if err != nil {
		return err
	}
	return nil
}

func copyDir(src, dst string) (err error) {
	// 목적지 폴더를 만든다.
	err = os.MkdirAll(dst, 0777)
	if err != nil {
		return err
	}
	// 소스 폴더 하위의 파일 목록을 가져온다.
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	// 파일 리스트를 for문 돌면서 하나씩 복사한다.
	for _, f := range files {
		if f.IsDir() { // 디렉토리라면 복사하지 않는다.
			continue
		}
		s := src + "/" + f.Name()
		d := dst + "/" + f.Name()
		err = copyFile(s, d)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}
	return nil
}
