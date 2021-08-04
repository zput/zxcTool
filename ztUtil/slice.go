package ztUtil

func Contains(str string, strs ...string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}
