package utils

import "strings"

func ConcatStrings(elems ...string) []string {
	return elems
}

func StrSplit(s string, p ...rune) []string {
	//arr := strings.FieldsFunc(s, func(c rune) bool { return c == ',' || c == ';' })

	f := func(c rune) bool {
		for i := 0; i < len(p); i++ {
			if c == p[i] {
				return true
			}
		}
		return false
	}

	return strings.FieldsFunc(s, f)

}
