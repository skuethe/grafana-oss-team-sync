package helpers

func RemoveFromSlice(s []string, r string) []string {
	for k, v := range s {
		if v == r {
			return append(s[:k], s[k+1:]...)
		}
	}
	return s
}
