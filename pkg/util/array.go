package util

func StringToAnyArray(str []string) (r []interface{}) {
	r = make([]interface{}, len(str))
	for i, v := range str {
		r[i] = v
	}
	return
}
