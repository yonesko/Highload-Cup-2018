package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
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
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			actualRespBody := string(bytes)
			if resp.StatusCode != r.expectedStatus || actualRespBody != r.expectedAnswerBody {
				fmt.Println(request.URL)
				require.Equal(t, r.expectedStatus, resp.StatusCode)
				require.Equal(t, parseBody(r.expectedAnswerBody), parseBody(actualRespBody))
			}
		}
	}
}

func parseBody(s string) body {
	var b body
	err := json.Unmarshal([]byte(s), &b)
	if err != nil {
		panic(err)
	}
	return b
}

type body struct {
	Accounts map[string]string `json:"accounts"`
}

type row struct {
	method             string
	query              string
	expectedStatus     int
	expectedAnswerBody string
}

var rows []row

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
