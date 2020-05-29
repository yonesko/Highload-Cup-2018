package api

import (
	"github.com/gin-gonic/gin"

	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/db"
)

type predicate struct {
	field string
	op    string
}

type predicateI interface {
	//sorted acc ids
	filter() []int64
}

func AccountsFilter(c *gin.Context) {
	preds, err := parsePredicate(c)
	if err != nil {
		c.Status(400)
		return
	}
	var accountIds []int64
	for _, p := range preds {
		result := p.filter()
		if len(result) == 0 {
			c.Status(200)
			c.JSON(200, gin.H{"accounts": []gin.H{}})
			return
		}
		accountIds = db.IntersectSorted(accountIds, result)
	}
	c.JSON(200, gin.H{"accounts": respBody(accountIds)})
}

func respBody(accountIds []int64) []gin.H {
	ans := make([]gin.H, len(accountIds))

	for _, id := range accountIds {
		account := db.Accounts[id]
		ans = append(ans, gin.H{
			"email":   account.Email,
			"country": account.Country,
			"id":      account.ID,
			"status":  account.Status,
			"birth":   account.Birth,
		})
	}

	return ans
}

//return sorted by selectivity predicates
func parsePredicate(c *gin.Context) ([]predicateI, error) {
	return nil, nil
}
