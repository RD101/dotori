package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
)

func handleAPIItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		i := Item{}
		//ParseForm parses the raw query from the URL and updates r.Form.
		r.ParseForm()
		for key, values := range r.PostForm {
			switch key {
			case "type":
				if len(values) != 1 {
					http.Error(w, "type을 설정해 주세요", http.StatusBadRequest)
					return
				}
				i.Type = values[0]
			case "author":
				if len(values) != 1 {
					http.Error(w, "author를 설정해 주세요", http.StatusBadRequest)
					return
				}
				i.Author = values[0]
			case "inputpath":
				if len(values) != 1 {
					http.Error(w, "inputpath를 설정해 주세요", http.StatusBadRequest)
					return
				}
				i.Inputpath = values[0]
			case "outputpath":
				if len(values) != 1 {
					http.Error(w, "outputpath를 설정해 주세요", http.StatusBadRequest)
					return
				}
				i.Outputpath = values[0]
			case "thumbimg":
				if len(values) != 1 {
					http.Error(w, "thumbnail image의 경로를 설정해 주세요", http.StatusBadRequest)
					return
				}
				i.Thumbimg = values[0]
			case "thumbmov":
				if len(values) != 1 {
					http.Error(w, "thumbnail mov의 경로를 설정해 주세요", http.StatusBadRequest)
					return
				}
				i.Thumbmov = values[0]
			}
		}
		err := i.CheckError()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer session.Close()
		err = AddItem(session, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, _ := json.Marshal(i)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else if r.Method == http.MethodDelete {
		var itemtype, id string
		r.ParseForm()
		for key, values := range r.Form {
			switch key {
			case "type":
				if len(values) != 1 {
					http.Error(w, "type을 설정해 주세요", http.StatusBadRequest)
					return
				}
				itemtype = values[0]
			case "id":
				if len(values) != 1 {
					http.Error(w, "id를 설정해 주세요", http.StatusBadRequest)
					return
				}
				id = values[0]
			}
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer session.Close()
		fmt.Printf("itemtype  %T  %s\n", itemtype, itemtype)
		fmt.Printf("id  %T  %s\n", id, id)
		RmItem(session, itemtype, id)
		data, _ := json.Marshal(id)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else {
		http.Error(w, "Not Supported Method", http.StatusMethodNotAllowed)
		return
	}

}
