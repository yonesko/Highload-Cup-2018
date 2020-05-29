package api

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yonesko/Highload-Cup-2018/db"
	"github.com/yonesko/Highload-Cup-2018/slice"
)

type filterField struct {
	name        string
	ops         []string
	selectivity int
}

var filterFieldsMap = map[string]filterField{}

func init() {
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
			selectivity: 90,
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
			selectivity: 80,
		},
		{
			name:        "city",
			ops:         []string{"eq", "any", "null"},
			selectivity: 70,
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

	for _, ff := range filterFields {
		filterFieldsMap[ff.name] = ff
	}
}

type predicate struct {
	field string
	op    string
	val   string
}

func (p predicate) filter() []int64 {
	switch p.field {
	case "sex":
		switch p.op {
		case "eq":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Sex == p.val {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	case "phone":
		switch p.op {
		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if p.val == "1" && a.Phone == "" {
					ans = append(ans, a.ID)
				}
				if p.val == "0" && a.Phone != "" {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	}
	return nil
}

func AccountsFilter(c *gin.Context) {
	limit, ok := parseLimit(c)
	if !ok {
		c.Status(400)
		return
	}
	preds, ok := parsePredicates(c)
	if !ok {
		c.Status(400)
		return
	}
	accountIds := preds[0].filter()
	for i := 1; i < len(preds); i++ {
		if len(accountIds) == 0 {
			c.Status(200)
			c.JSON(200, gin.H{"accounts": []gin.H{}})
			return
		}
		accountIds = db.IntersectSorted(accountIds, preds[i].filter())
	}

	c.JSON(200, gin.H{"accounts": respBody(accountIds, limit, preds)})
}

func parseLimit(c *gin.Context) (int, bool) {
	l, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return 0, false
	}
	return l, true
}
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func respBody(accountIds []int64, limit int, preds []predicate) []gin.H {
	ans := make([]gin.H, limit)

	sort.Slice(accountIds, func(i, j int) bool {
		return accountIds[i] > accountIds[j]
	})
	for i := 0; i < min(limit, len(accountIds)); i++ {
		account := db.Accounts[accountIds[i]]
		ans[i] = gin.H{"id": account.ID, "email": account.Email}
		for _, p := range preds {
			switch p.field {
			case "sex":
				ans[i]["sex"] = account.Sex
			case "phone":
				if account.Phone != "" {
					ans[i]["phone"] = account.Phone
				}
			}
		}
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
		if k == "limit" {
			continue
		}
		if p, ok := parsePred(k, vals[0]); !ok {
			return nil, false
		} else {
			ans = append(ans, p)
		}
	}
	sort.Slice(ans, func(i, j int) bool {
		return filterFieldsMap[ans[i].field].selectivity > filterFieldsMap[ans[j].field].selectivity
	})
	return ans, true
}

func parsePred(key string, val string) (predicate, bool) {
	split := strings.Split(key, "_")
	if len(split) != 2 {
		return predicate{}, false
	}
	p := predicate{field: split[0], op: split[1], val: val}
	if !validatePred(p) {
		return predicate{}, false
	}
	return p, true
}
func validatePred(p predicate) bool {
	if ff, ok := filterFieldsMap[p.field]; ok && slice.StringSliceContains(ff.ops, p.op) {
		return true
	}
	return false
}
