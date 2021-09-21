package slices

func InSlice(v interface{}, slice []interface{}) bool {
	for _, s := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func InSliceStr(v string, slices []string) bool {
	for _, s := range slices {
		if v == s {
			return true
		}
	}
	return false
}
