package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	client := &http.Client{}
	for _, r := range rows {
		if strings.Contains(r.query, "/accounts/filter/") {
			request, err := http.NewRequest(r.method, fmt.Sprintf("http://localhost:80%s", r.query), nil)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode != r.expectedStatus {
				fmt.Printf("Unexpected status, want = %d, got = %d\n", r.expectedStatus, resp.StatusCode)
				fmt.Println(r)
				bytes, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(bytes))
			}
		}
		os.Exit(0)
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
