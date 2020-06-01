package slice

func StringSliceContains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func StringSliceIntersect(a []string, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return []string{}
	}

	aMap := StringSliceToSet(a)
	var ans []string

	for _, v := range b {
		if aMap[v] {
			ans = append(ans, v)
		}
	}

	return ans
}

func StringSliceToSet(a []string) map[string]bool {
	ans := map[string]bool{}

	for _, v := range a {
		ans[v] = true
	}

	return ans
}
