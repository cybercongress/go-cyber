package util

import "strings"

func HasPrefixIgnoreCase(s, prefix string) bool {
	lowerCaseS := strings.ToLower(s)
	lowerCasePrefix := strings.ToLower(prefix)

	return strings.HasPrefix(lowerCaseS, lowerCasePrefix)
}
