package utils

import "strings"

// ConcatStrings 将多个字符串合成一个字符串数组
func ConcatStrings(elems ...string) []string {
	return elems
}

// StrSplit 使用多个字符对字符进行切割
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
