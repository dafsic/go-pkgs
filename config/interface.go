package config

import (
	"reflect"
	"strings"
)

type Config interface {
	GetItem(e string) interface{}
	SetItem(k string, v interface{})
}

func GetItem(conf *interface{}, e string) interface{} {
	var cfg interface{}
	rt := reflect.TypeOf(*conf)
	rv := reflect.ValueOf(*conf)

	fieldNum := rt.NumField()
	for i := 0; i < fieldNum; i++ {
		if strings.ToUpper(rt.Field(i).Name) == strings.ToUpper(e) {
			cfg = rv.FieldByName(rt.Field(i).Name).Interface()
			break
		}
	}
	return cfg
}

// TODO: 待测试
// SetItem 设置conf中字段名字为k的值为v
func SetItem(conf *interface{}, k string, v interface{}) {
	rt := reflect.TypeOf(conf)
	rv := reflect.ValueOf(conf).Elem()

	fieldNum := rt.NumField()
	for i := 0; i < fieldNum; i++ {
		if strings.ToUpper(rt.Field(i).Name) == strings.ToUpper(k) {
			rv.Field(i).Set(reflect.ValueOf(v))
			break
		}
	}
}
