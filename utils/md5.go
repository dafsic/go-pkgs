package utils

import (
	"crypto/md5"
	"fmt"
)

func StrToMD5(v string) string {
	if v == "" {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(v)))
}
