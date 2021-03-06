package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonesko/Highload-Cup-2018/db"
)

func Test(t *testing.T) {
	db.LoadAccounts()
	ammo := loadAmmo()
	gin := buildGin()
	for _, r := range ammo {
		if strings.Contains(r.query, "/accounts/filter/") && !ignore(r.query) {
			rrun := r
			t.Run(rrun.query, func(t *testing.T) {
				t.Parallel()
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(rrun.method, rrun.query, nil)
				gin.ServeHTTP(w, req)
				actualRespBody := w.Body.String()
				assert.Equal(t, rrun.expectedStatus, w.Code)
				if w.Code == 200 {
					assert.Equal(t, parseBody(rrun.expectedAnswerBody), parseBody(actualRespBody))
				}
			})
		}
	}
}

func ignore(query string) bool {
	return strings.Contains(query, "premium_now")
}

func parseBody(s string) body {
	var b body
	err := json.Unmarshal([]byte(s), &b)
	if err != nil {
		fmt.Println(s)
		panic(err)
	}
	return b
}

type body struct {
	Accounts []map[string]interface{} `json:"accounts"`
}

type row struct {
	method             string
	query              string
	expectedStatus     int
	expectedAnswerBody string
}

func loadAmmo() []row {
	file, err := ioutil.ReadFile("/Users/gbdanichev/Downloads/test_accounts_240119/answers/phase_1_get.answ")
	if err != nil {
		log.Fatal(err)
	}
	fileStr := string(file)
	var rows []row
	for _, line := range strings.Split(fileStr, "\n") {
		vals := strings.Split(line, "\t")
		if len(vals) < 3 {
			continue
		}
		status, err := strconv.Atoi(vals[2])
		if err != nil {
			log.Fatal(err)
		}

		r := row{method: vals[0], query: vals[1], expectedStatus: status}
		if len(vals) == 4 {
			r.expectedAnswerBody = vals[3]
		}
		rows = append(rows, r)
	}
	return rows
}
