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
		k2 := strings.ToUpper(k[:1]) + k[1:]
		ans += fmt.Sprintf(tmpl+"\n", k, k2, k, k2)
	}

	fmt.Println(ans)
}

func TestGenerateFilter(t *testing.T) {
	//tmpl := `case "%s":
	//	switch p.op {
	//	case "%s":
	//		var ans []int64
	//		for _, a := range db.Accounts {
	//			if a.%s == p.val {
	//				ans = append(ans, a.ID)
	//			}
	//		}
	//		return ans
	//	}` + "\n"
	tmpl := `{{range $key, $ff := .}}
case "{{$key}}":
switch p.op {
	{{range $ff.Ops }}
		case "{{.}}":
	{{end}}
}
{{end}}
`
	tm := template.Must(template.New("").Parse(tmpl))
	require.Nil(t, tm.Execute(os.Stdout, filterFieldsMap))
}
