package api

import (
	"fmt"
	"strings"
	"testing"
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
