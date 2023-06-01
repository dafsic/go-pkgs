package utils

import "reflect"

// StructUpdate 将new结构体中字段不为0值的更新到old结构中，前提0值不能是有意义的
func StructUpdate(old, new interface{}) {
	valueOfOld := reflect.ValueOf(old).Elem()
	valueOfNew := reflect.ValueOf(new).Elem()

	for i := 0; i < valueOfNew.NumField(); i++ {
		if valueOfNew.Field(i).Interface() != reflect.Zero(valueOfNew.Field(i).Type()).Interface() {
			valueOfOld.Field(i).Set(valueOfNew.Field(i))
		}
	}
}
