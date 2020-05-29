package main

import (
	"github.com/gin-gonic/gin"
)

func accountsFilter(c *gin.Context) {
	//for _,p := range c.Params {
	//
	//}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
