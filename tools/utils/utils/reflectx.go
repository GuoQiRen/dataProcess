package utils

import (
	"encoding/json"
	"reflect"
	"strconv"
)

/**
反射获取struct各个field的名称，并拷贝赋值
*/
func ReflectStructAssign(from, to interface{}) {
	f := reflect.ValueOf(from).Elem()
	t := reflect.ValueOf(to).Elem()
	fType := f.Type()
	for i := 0; i < f.NumField(); i++ {
		fName := fType.Field(i).Name
		if ok := t.FieldByName(fName).IsValid(); ok {
			t.FieldByName(fName).Set(reflect.ValueOf(f.Field(i).Interface()))
		}
	}
}

/**
反射获取json的tags
*/
func ReflectGetJsonTags(stu interface{}) (tags []string) {
	t := reflect.TypeOf(stu).Elem()
	for i := 0; i < t.NumField(); i++ {
		tags = append(tags, t.Field(i).Tag.Get("json"))
	}
	return
}

/**
通过tags去赋值
*/
func ReflectSetStructValue(stu interface{}, tags []string, values []string) {
	var kid reflect.Kind

	valOf := reflect.ValueOf(stu).Elem()
	valType := valOf.Type()

	for ind := range tags {
		name := valType.Field(ind).Name
		field := valOf.FieldByName(name)
		if ok := field.IsValid(); ok {
			kid = field.Kind()
		}

		switch kid {
		case reflect.String:
			strVal := values[ind]
			field.SetString(strVal)
		case reflect.Float32, reflect.Float64:
			float64Val, _ := strconv.ParseFloat(values[ind], 64)
			field.SetFloat(float64Val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			int64Val, _ := strconv.ParseInt(values[ind], 10, 64)
			field.SetInt(int64Val)
		case reflect.Bool:
			boolVal, _ := strconv.ParseBool(values[ind])
			field.SetBool(boolVal)
		case reflect.Array, reflect.Slice:
			var sliceObj []interface{}
			_ = json.Unmarshal([]byte(values[ind]), &sliceObj)
			field.Set(reflect.ValueOf(sliceObj))
		default:
			continue
		}
	}
	return
}
