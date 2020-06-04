package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"

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
	val   interface{}
}

func (p predicate) filter() []int64 {
	switch p.field {
	case "birth":
		switch p.op {
		case "lt":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Birth < p.val.(int64) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "gt":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Birth > p.val.(int64) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "year":
			var ans []int64
			for _, a := range db.Accounts {
				if a.UTCBirthYear() == p.val.(int) {
					ans = append(ans, a.ID)
				}
			}
			return ans
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
			var ans []int64
			for _, a := range db.Accounts {
				if slice.StringSliceContains(p.val.([]string), a.City) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val.(bool) && a.City == "") || (!p.val.(bool) && a.City != "") {
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
				if (p.val.(bool) && a.Country == "") || (!p.val.(bool) && a.Country != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	case "email":
		switch p.op {
		case "domain":
			var ans []int64
			for _, a := range db.Accounts {
				if a.EmailDomain() == p.val.(string) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "lt":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Email < p.val.(string) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "gt":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Email > p.val.(string) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	case "fname":
		switch p.op {
		case "eq":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Fname == p.val.(string) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "any":
			var ans []int64
			for _, a := range db.Accounts {
				if slice.StringSliceContains(p.val.([]string), a.Fname) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val.(bool) && a.Fname == "") || (!p.val.(bool) && a.Fname != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	case "interests":
		switch p.op {
		case "contains":
			var ans []int64
			for _, a := range db.Accounts {
				if len(slice.StringSliceIntersect(p.val.([]string), a.Interests)) == len(p.val.([]string)) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "any":
			var ans []int64
			for _, a := range db.Accounts {
				if len(slice.StringSliceIntersect(p.val.([]string), a.Interests)) > 0 {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	case "likes":
		switch p.op {
		case "contains":
			var ans []int64
			for _, a := range db.Accounts {
				if len(funk.Join(p.val.([]int64), a.LikesIds(), funk.InnerJoin).([]int64)) == len(p.val.([]int64)) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	case "phone":
		switch p.op {
		case "code":
			var ans []int64
			for _, a := range db.Accounts {
				if a.PhoneCode() == p.val {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val.(bool) && a.Phone == "") || (!p.val.(bool) && a.Phone != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	case "premium":
		switch p.op {
		case "now":
			if p.val.(bool) {
				var ans []int64
				for _, a := range db.Accounts {
					if a.Premium.Start <= time.Now().Unix() && a.Premium.Finish >= time.Now().Unix() {
						ans = append(ans, a.ID)
					}
				}
				return ans
			}
		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val.(bool) && a.Premium == account.Premium{}) || (!p.val.(bool) && a.Premium != account.Premium{}) {
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
				if a.Sex == p.val.(string) {
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
				if a.Sname == p.val.(string) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "starts":
			var ans []int64
			for _, a := range db.Accounts {
				if strings.HasPrefix(a.Sname, p.val.(string)) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "null":
			var ans []int64
			for _, a := range db.Accounts {
				if (p.val.(bool) && a.Sname == "") || (!p.val.(bool) && a.Sname != "") {
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
				if a.Status == p.val.(string) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		case "neq":
			var ans []int64
			for _, a := range db.Accounts {
				if a.Status != p.val.(string) {
					ans = append(ans, a.ID)
				}
			}
			return ans
		}
	}
	return nil
}

var debug = false

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
	if len(preds) == 0 {
		c.JSON(200, gin.H{"accounts": respBody(funk.Keys(db.Accounts).([]int64), limit, preds)})
		return
	}
	accountIds := preds[0].filter()
	for i := 1; i < len(preds); i++ {
		if debug {
			fmt.Printf("filtered by %v got %d\n", preds[i-1], len(accountIds))
		}
		if len(accountIds) == 0 {
			c.Status(200)
			c.JSON(200, gin.H{"accounts": []gin.H{}})
			return
		}
		accountIds = funk.JoinInt64(accountIds, preds[i].filter(), funk.InnerJoinInt64)
	}
	if debug {
		fmt.Printf("filtered by %v got %d\n", preds[len(preds)-1], len(accountIds))
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
	respAccountsSize := min(limit, len(accountIds))
	ans := make([]gin.H, respAccountsSize)

	sort.Slice(accountIds, func(i, j int) bool {
		return accountIds[i] > accountIds[j]
	})
	for i := 0; i < respAccountsSize; i++ {
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
	predVal, ok := parsePredVal(p)
	if !ok {
		return predicate{}, false
	}
	p.val = predVal
	return p, true
}
func validatePred(p predicate) bool {
	if ff, ok := filterFieldsMap[p.field]; ok && slice.StringSliceContains(ff.Ops, p.op) {
		return true
	}
	return false
}

func parsePredVal(p predicate) (interface{}, bool) {
	if p.op == "null" {
		if p.val.(string) == "1" {
			return true, true
		} else if p.val.(string) == "0" {
			return false, true
		} else {
			return nil, false
		}
	}
	switch p.field {
	case "birth":
		switch p.op {
		case "gt", "lt":
			dt, err := strconv.ParseInt(p.val.(string), 10, 64)
			if err != nil {
				return nil, false
			}
			return dt, true
		case "year":
			y, err := strconv.Atoi(p.val.(string))
			if err != nil {
				return nil, false
			}
			return y, true
		}
	case "city":
		switch p.op {
		case "any":
			return strings.Split(p.val.(string), ","), true
		}
	case "country":
		switch p.op {
		case "eq":
			return p.val, true
		}
	case "fname":
		switch p.op {
		case "any":
			return strings.Split(p.val.(string), ","), true
		}
	case "interests":
		switch p.op {
		case "any", "contains":
			return strings.Split(p.val.(string), ","), true
		default:
			return nil, false
		}
	case "likes":
		switch p.op {
		case "contains":
			var ans []int64
			for _, idStr := range strings.Split(p.val.(string), ",") {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					return nil, false
				}
				ans = append(ans, id)
			}
			return ans, true
		}
	case "phone":
		switch p.op {
		case "code":
			return p.val, true
		}
	case "premium":
		switch p.op {
		case "now":
			if p.val.(string) == "1" {
				return true, true
			} else {
				return nil, false
			}
		}
	case "sex":
		switch p.op {
		case "eq":
			if p.val.(string) == "m" {
				return "m", true
			} else if p.val.(string) == "f" {
				return "f", true
			} else {
				return nil, false
			}
		}
	case "sname":
		switch p.op {
		case "starts", "eq":
			return p.val.(string), true
		}
	case "status":
		if slice.StringSliceContains([]string{"свободны", "заняты", "всё сложно"}, p.val.(string)) {
			return p.val.(string), true
		} else {
			return nil, false
		}
	}

	return p.val, true
}
