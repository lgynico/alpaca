package helper

import (
	"strings"
	"unicode"
)

func UnderlineToCamelCase(value string, capitalizeFirst bool) string {
	ss := strings.Split(value, "_")
	for i, v := range ss {
		if i == 0 && !capitalizeFirst {
			ss[i] = strings.ToLower(string(v[0])) + v[1:]
		} else {
			ss[i] = strings.ToUpper(string(v[0])) + v[1:]
		}
	}
	return strings.Join(ss, "")
}

func CapitalizeLeading(value string) string {
	if len(value) == 0 {
		return value
	}

	runes := []rune(value)
	if unicode.IsLetter(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes)
}
