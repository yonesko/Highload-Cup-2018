package main

import (
	"log"
	"net/http"

	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/db"
)

func main() {
	db.LoadAccounts()

	http.HandleFunc("/accounts/filter/", accountsFilter)
	log.Fatal(http.ListenAndServe(":80", nil))
}
