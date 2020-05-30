package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yonesko/Highload-Cup-2018/account"
	"github.com/yonesko/Highload-Cup-2018/db"
	"github.com/yonesko/Highload-Cup-2018/slice"
)

type filterField struct {
	name        string
	Ops         []string
	selectivity int
}

var filterFieldsMap = map[string]filterField{}

func init() {
	var filterFields = []filterField{
		{
			name:        "sex",
			Ops:         []string{"eq"},
			selectivity: 100,
		},
		{
			name:        "email",
			Ops:         []string{"domain", "lt", "gt"},
			selectivity: 0,
		},
		{
			name:        "status",
			Ops:         []string{"eq", "neq"},
			selectivity: 90,
		},
		{
			name:        "fname",
			Ops:         []string{"eq", "any", "null"},
			selectivity: 0,
		},
		{
			name:        "sname",
			Ops:         []string{"eq", "starts", "null"},
			selectivity: 0,
		},
		{
			name:        "phone",
			Ops:         []string{"code", "null"},
			selectivity: 0,
		},
		{
			name:        "country",
			Ops:         []string{"eq", "null"},
			selectivity: 80,
		},
		{
			name:        "city",
			Ops:         []string{"eq", "any", "null"},
			selectivity: 70,
		},
		{
			name:        "birth",
			Ops:         []string{"lt", "gt", "year"},
			selectivity: 0,
		},
		{
			name:        "interests",
			Ops:         []string{"contains", "any"},
			selectivity: 0,
		},
		{
			name:        "likes",
			Ops:         []string{"contains"},
			selectivity: 0,
		},
		{
			name:        "premium",
			Ops:         []string{"now", "null"},
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
	case "birth":
		switch p.op {
		case "lt":

		case "gt":

		case "year":

		}
	case "city":
		switch p.op {
		case "eq":
			var ans []int64
			for _, a := range db.Accounts {
				if a.City == p.val {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "any":

		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val == "1" && a.City == "") || (p.val == "0" && a.City != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans

		}
	case "country":
		switch p.op {
		case "eq":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Country == p.val {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val == "1" && a.Country == "") || (p.val == "0" && a.Country != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans

		}
	case "email":
		switch p.op {
		case "domain":

		case "lt":

		case "gt":

		}
	case "fname":
		switch p.op {
		case "eq":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Fname == p.val {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "any":

		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val == "1" && a.Fname == "") || (p.val == "0" && a.Fname != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans

		}
	case "interests":
		switch p.op {
		case "contains":

		case "any":

		}
	case "likes":
		switch p.op {
		case "contains":

		}
	case "phone":
		switch p.op {
		case "code":

		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val == "1" && a.Phone == "") || (p.val == "0" && a.Phone != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans

		}
	case "premium":
		switch p.op {
		case "now":

		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val == "1" && a.Premium == account.Premium{}) || (p.val == "0" && a.Premium != account.Premium{}) {
					ans = append(ans, a.ID)
				}
			}
			return ans

		}
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
	case "sname":
		switch p.op {
		case "eq":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Sname == p.val {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "starts":

		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val == "1" && a.Sname == "") || (p.val == "0" && a.Sname != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans

		}
	case "status":
		switch p.op {
		case "eq":

			var ans []int64
			for _, a := range db.Accounts {
				if a.Status == p.val {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "neq":

		}
	}
	return nil
}

var debug = true

func AccountsFilter(c *gin.Context) {
	limit, ok := parseLimit(c)
	if !ok {
		c.Status(400)
		return
	}
	preds, ok := parsePredicates(c)
	if debug {
		fmt.Printf("parsePredicates %v\n", preds)
	}
	if !ok {
		c.Status(400)
		return
	}
	accountIds := preds[0].filter()
	for i := 1; i < len(preds); i++ {
		if debug {
			fmt.Printf("accountIds %d\n", len(accountIds))
		}
		if len(accountIds) == 0 {
			c.Status(200)
			c.JSON(200, gin.H{"accounts": []gin.H{}})
			return
		}
		accountIds = db.IntersectSorted(accountIds, preds[i].filter())
	}
	if debug {
		fmt.Printf("accountIds %d\n", len(accountIds))
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
		acc := db.Accounts[accountIds[i]]
		ans[i] = gin.H{"id": acc.ID, "email": acc.Email}
		for _, p := range preds {
			switch p.field {
			case "phone":
				if acc.Phone != "" {
					ans[i]["phone"] = acc.Phone
				}
			case "country":
				if acc.Country != "" {
					ans[i]["country"] = acc.Country
				}
			case "birth":
				if acc.Birth != 0 {
					ans[i]["birth"] = acc.Birth
				}
			case "premium":
				if acc.Premium != (account.Premium{}) {
					ans[i]["premium"] = acc.Premium
				}
			case "sex":
				if acc.Sex != "" {
					ans[i]["sex"] = acc.Sex
				}
			case "status":
				if acc.Status != "" {
					ans[i]["status"] = acc.Status
				}
			case "sname":
				if acc.Sname != "" {
					ans[i]["sname"] = acc.Sname
				}
			case "city":
				if acc.City != "" {
					ans[i]["city"] = acc.City
				}
			case "fname":
				if acc.Fname != "" {
					ans[i]["fname"] = acc.Fname
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
	if ff, ok := filterFieldsMap[p.field]; ok && slice.StringSliceContains(ff.Ops, p.op) {
		return true
	}
	return false
}
