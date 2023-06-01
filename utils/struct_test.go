package utils

import "testing"

func TestStrctUpdate(t *testing.T) {

	type tmp struct {
		Id   int
		Name string
	}

	var a, b tmp

	a.Id = 1
	b.Name = "libai"

	StructUpdate(&a, &b)

	if a.Id != 1 || a.Name != "libai" {
		t.Error(a)
	}

}
