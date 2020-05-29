package main

import (
	"log"
	"net/http"

	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/db"
)

func DumbHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello, I'm Web Server!"))
}

func main() {
	db.LoadAccounts()

	http.HandleFunc("/", DumbHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
