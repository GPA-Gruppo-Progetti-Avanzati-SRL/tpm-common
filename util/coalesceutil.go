package util

func StringCoalesce(args ...string) string {
	for _, str := range args {
		if str != "" {
			return str
		}
	}
	return ""
}

func IntCoalesce(args ...int) int {
	for _, str := range args {
		if str > 0 {
			return str
		}
	}
	return 0
}
