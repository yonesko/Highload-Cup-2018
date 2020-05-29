package api

import "github.com/gin-gonic/gin"

type predicate struct {
	field string
	op    string
}

func AccountsFilter(c *gin.Context) {
	//fmt.Println(c.Params)

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
