package tag

import (
	"fmt"
	"reflect"
	"strings"
)

func GetItemNameByTag(field reflect.StructField, tag string) string {
	var itemName string
	if itemName = field.Tag.Get(tag); itemName != "" {
		s := strings.IndexByte(itemName, ',')
		if s > 0 {
			itemName = itemName[:s]
		}
	} else {
		itemName = field.Name
	}
	return itemName
}

func RangeFields(v interface{}, tag string, fn func(valField reflect.Value, typField reflect.StructField)) {
	vVal := reflect.ValueOf(v)
	if vVal.IsValid() == false {
		return
	}
	vTyp := vVal.Type()
	if vVal.Kind() == reflect.Ptr {
		vTyp = vTyp.Elem()
		if vVal.IsNil() {
			vVal = reflect.Zero(vTyp)
		} else {
			vVal = vVal.Elem()
		}
	}

	if vVal.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s type must be struct", vVal.Type().Name()))
	}
	for i := 0; i < vVal.NumField(); i++ {
		field := vVal.Field(i)
		typField := vTyp.Field(i)
		if GetItemNameByTag(typField, tag) == "-" {
			continue
		}
		if typField.Anonymous {
			if (field.Kind() == reflect.Interface || field.Kind() == reflect.Ptr) && field.IsNil() {
				panic("cannot parse nil point anonymous fields")
			}
			checkField := field
			if checkField.Kind() == reflect.Interface || checkField.Kind() == reflect.Ptr {
				checkField = checkField.Elem()
			}
			if checkField.Kind() == reflect.Struct {
				RangeFields(field.Interface(), tag, fn)
			} else {
				fmt.Println("warning: ignore Anonymous field ", typField.Name, typField.Type.Name())
			}

		} else {
			fn(field, typField)
		}
	}
}
