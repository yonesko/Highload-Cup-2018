package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/yonesko/Highload-Cup-2018/api"
	"github.com/yonesko/Highload-Cup-2018/db"
)

func main() {
	db.LoadAccounts()

	r := gin.Default()
	r.GET("/accounts/filter/", api.AccountsFilter)
	log.Fatal(r.Run(":80"))
}
