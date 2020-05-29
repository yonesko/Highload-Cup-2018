package db

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/account"
)

var Accounts = map[int64]account.Account{}

func LoadAccounts() {
	reader, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		log.Fatalf("can't load accounts: %s", err)
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
		accs := struct {
			Accounts []account.Account `json:"accounts"`
		}{}
		err = json.Unmarshal(bytes, &accs)
		if err != nil {
			log.Fatal(err)
		}
		for _, a := range accs.Accounts {
			Accounts[a.ID] = a
		}
	}
	fmt.Printf("Load completed, len = %d\n", len(Accounts))
	//fmt.Printf("%#v", Accounts[0])
}
