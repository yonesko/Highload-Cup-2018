package db

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
)

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
		rc.Close()
		fmt.Println(string(bytes))
	}
}
