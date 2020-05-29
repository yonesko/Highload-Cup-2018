package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type row struct {
	method             string
	query              string
	expectedStatus     int
	expectedAnswerBody string
}

var rows []row

func main() {
	loadAmmo()
	for _, r := range rows {
		if strings.Contains(r.query, "/accounts/filter/") {
			fmt.Println(r.query)
		}
	}
}

func loadAmmo() {
	file, err := ioutil.ReadFile("/Users/gbdanichev/Downloads/test_accounts_240119/answers/phase_1_get.answ")
	if err != nil {
		log.Fatal(err)
	}
	fileStr := string(file)

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
}
