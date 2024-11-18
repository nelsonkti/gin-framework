package helper

func If(bool2 bool, value interface{}, value2 interface{}) interface{} {
	if bool2 {
		return value
	}
	return value2
}

// Max php max
func Max(args ...int) int {
	max := args[0]
	for _, v := range args {
		if v > max {
			max = v
		}
	}
	return max
}

// MaxInt64 php max
func MaxInt64(args ...int64) int64 {
	max := args[0]
	for _, v := range args {
		if v > max {
			max = v
		}
	}
	return max
}
