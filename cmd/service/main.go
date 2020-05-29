package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/api"
	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/db"
)

func main() {
	db.LoadAccounts()

	r := gin.Default()
	r.GET("/accounts/filter/", api.AccountsFilter)
	log.Fatal(r.Run(":80"))
}
