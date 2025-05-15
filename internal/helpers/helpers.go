package helpers

import "strings"

func RemoveFromSlice(searchHere []string, removeThis string, caseSensitive bool) []string {
	for k, v := range searchHere {
		var validated bool
		if caseSensitive {
			validated = v == removeThis
		} else {
			validated = strings.EqualFold(v, removeThis)
		}
		if validated {
			return append(searchHere[:k], searchHere[k+1:]...)
		}
	}
	return searchHere
}
