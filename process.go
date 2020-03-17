package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

func processingItem() error {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()
	// Status가 Ready인 item을 가져온다.
	item, err := GetReadyItem(session)
	if err != nil {
		return err
	}
	fmt.Println(item)
	return nil
}
