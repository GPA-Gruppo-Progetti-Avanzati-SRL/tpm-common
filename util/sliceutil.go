package util

func RemoveStringDuplicates(d []string) []string {

	check := make(map[string]int)

	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for s, _ := range check {
		res = append(res, s)
	}

	return res
}
