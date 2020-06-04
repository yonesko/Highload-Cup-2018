package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/yonesko/Highload-Cup-2018/api"
	"github.com/yonesko/Highload-Cup-2018/db"
)

func main() {
	db.LoadAccounts()
	gin.SetMode(gin.ReleaseMode)
	r := buildGin()
	log.Fatal(r.Run(":80"))
}

func buildGin() *gin.Engine {
	r := gin.Default()
	r.GET("/accounts/filter/", api.AccountsFilter)
	return r
}
