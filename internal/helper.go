package internal

import "strings"

func ContainsAll(container []string, containing []string) bool {
	for _, v := range containing {
		if !Contains(container, v) {
			return false
		}
	}
	return true
}

func Contains(arr []string, elt string) bool {
	for _, v := range arr {
		if elt == v {
			return true
		}
	}
	return false
}

func MatchAny(container string, containing []string) bool {
	container = strings.ToLower(container)
	for _, v := range containing {
		if strings.Contains(container, v) {
			return true
		}
	}
	return false
}
