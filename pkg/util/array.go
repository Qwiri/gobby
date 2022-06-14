package util

func ConvertArray[T any](arr []T) (r []interface{}) {
	r = make([]interface{}, len(arr))
	for i, v := range arr {
		r[i] = v
	}
	return r
}
