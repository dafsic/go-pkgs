package config

import (
	"reflect"
	"strings"
)

type Config interface {
	GetElem(e string) interface{}
}

func GetElem(conf interface{}, e string) interface{} {
	var cfg interface{}
	rt := reflect.TypeOf(conf)
	rv := reflect.ValueOf(conf)

	fieldNum := rt.NumField()
	for i := 0; i < fieldNum; i++ {
		if strings.ToUpper(rt.Field(i).Name) == strings.ToUpper(e) {
			cfg = rv.FieldByName(rt.Field(i).Name).Interface()
			break
		}
	}
	return cfg
}
