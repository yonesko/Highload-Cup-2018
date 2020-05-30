package api

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	tmpl := `case "%s":
				if account.%s != "" {
					ans[i]["%s"] = account.%s
				}`
	ans := ""

	for k := range filterFieldsMap {
		k2 := UpFirst(k)
		ans += fmt.Sprintf(tmpl+"\n", k, k2, k, k2)
	}

	fmt.Println(ans)
}

func UpFirst(k string) string {
	return strings.ToUpper(k[:1]) + k[1:]
}

func TestGenerateFilter(t *testing.T) {
	tmpl := `{{range $key, $ff := .}}
case "{{$key}}":
switch p.op {
	{{range $ff.Ops -}}
		case "{{.}}":
			{{if eq . "null" -}}
				var ans []int64
			for _, a := range db.Accounts {
				if (p.val == "1" && a.{{UpFirst $key}} == "") || (p.val == "0" && a.{{UpFirst $key}} != "") {
					ans = append(ans, a.ID)
				}
			}
			return ans
			{{- end}}
	{{- end}}
}
{{- end}}
`
	tm := template.Must(template.New("").Funcs(template.FuncMap{
		"UpFirst": UpFirst,
	}).Parse(tmpl))
	require.Nil(t, tm.Execute(os.Stdout, filterFieldsMap))
}
