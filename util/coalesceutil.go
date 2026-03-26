package util

func StringCoalesce(args ...string) string {
	for _, str := range args {
		if str != "" {
			return str
		}
	}
	return ""
}

func PositiveIntCoalesce(args ...int) int {
	for _, str := range args {
		if str > 0 {
			return str
		}
	}
	return 0
}

func IntCoalesce(args ...int) int {
	for _, str := range args {
		if str != 0 {
			return str
		}
	}
	return 0
}

func Int64Coalesce(args ...int64) int64 {
	for _, str := range args {
		if str != 0 {
			return str
		}
	}
	return 0
}

func Int32Coalesce(args ...int32) int32 {
	for _, str := range args {
		if str != 0 {
			return str
		}
	}
	return 0
}

func Float64Coalesce(args ...float64) float64 {
	for _, str := range args {
		if str != 0 {
			return str
		}
	}
	return 0.0
}

func CoalesceError(err ...error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}
	return nil
}
