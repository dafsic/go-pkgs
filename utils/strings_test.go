package utils

import (
	"strings"
	"testing"
)

func TestStrLink(t *testing.T) {
	s := ConcatStrings("123", "456", "789")

	if strings.Join(s, "|") != "123|456|789" {
		t.Error(s)
	}
}
