package foundation

import (
	"unicode"
)

func UnderLineToCamel(s string) string {
	runes := []rune(s)
	var result []rune

	nextToUpper := false
	for i, r := range runes {
		if r == '_' {
			nextToUpper = true
		} else if nextToUpper {
			result = append(result, unicode.ToUpper(r))
			nextToUpper = false
		} else if i == 0 {
			result = append(result, unicode.ToUpper(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
