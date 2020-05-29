package db

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/account"
)

var Accounts []account.Account

func LoadData() {
	reader, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		log.Fatalf("can't load data: %s", err)
	}
	defer reader.Close()
	for _, f := range reader.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Fatal(err)
		}
		_ = rc.Close()
		type Accs struct {
			Accounts []account.Account `json:"accounts"`
		}
		accs := Accs{}
		err = json.Unmarshal(bytes, &accs)
		if err != nil {
			log.Fatal(err)
		}
		Accounts = append(Accounts, accs.Accounts...)
	}
	fmt.Printf("Load completed, len = %d\n", len(Accounts))
	//fmt.Printf("%#v", Accounts[0])
}
