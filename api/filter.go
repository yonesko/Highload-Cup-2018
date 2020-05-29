package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/db"
	"go.avito.ru/github.com/yonesko/Highload-Cup-2018/slice"
)

type filterField struct {
	name        string
	ops         []string
	selectivity int
}

var filterFields = []filterField{
	{
		name:        "sex",
		ops:         []string{"eq"},
		selectivity: 100,
	},
	{
		name:        "email",
		ops:         []string{"domain", "lt", "gt"},
		selectivity: 0,
	},
	{
		name:        "status",
		ops:         []string{"eq", "neq"},
		selectivity: 0,
	},
	{
		name:        "fname",
		ops:         []string{"eq", "any", "null"},
		selectivity: 0,
	},
	{
		name:        "sname",
		ops:         []string{"eq", "starts", "null"},
		selectivity: 0,
	},
	{
		name:        "phone",
		ops:         []string{"code", "null"},
		selectivity: 0,
	},
	{
		name:        "country",
		ops:         []string{"eq", "null"},
		selectivity: 0,
	},
	{
		name:        "city",
		ops:         []string{"eq", "any", "null"},
		selectivity: 0,
	},
	{
		name:        "birth",
		ops:         []string{"lt", "gt", "year"},
		selectivity: 0,
	},
	{
		name:        "interests",
		ops:         []string{"contains", "any"},
		selectivity: 0,
	},
	{
		name:        "likes",
		ops:         []string{"eq", "contains"},
		selectivity: 0,
	},
	{
		name:        "premium",
		ops:         []string{"now", "null"},
		selectivity: 0,
	},
}

type predicate struct {
	field string
	op    string
	val   string
}

func (p predicate) filter() []int64 {
	panic("implement me")
}

func AccountsFilter(c *gin.Context) {
	preds, ok := parsePredicates(c)
	if !ok {
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
func parsePredicates(c *gin.Context) ([]predicate, bool) {
	var ans []predicate
	for k, vals := range c.Request.URL.Query() {
		if k == "query_id" {
			continue
		}
		if p, ok := parsePred(k, vals[0]); !ok {
			return nil, false
		} else {
			ans = append(ans, p)
		}
	}
	//SORT by SELECT
	return ans, true
}

func parsePred(key string, val string) (predicate, bool) {
	split := strings.Split(key, "_")
	if len(split) != 2 {
		return predicate{}, false
	}
	p := predicate{field: split[0], op: split[1], val: val}
	fmt.Println(p)
	if !validatePred(p) {
		return predicate{}, false
	}
	return p, true
}
func validatePred(p predicate) bool {
	for _, ff := range filterFields {
		if ff.name == p.field && slice.StringSliceContains(ff.ops, p.op) {
			return true
		}
	}
	return false
}
