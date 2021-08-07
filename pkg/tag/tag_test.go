package tag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type TestBDto struct {
	Id    string `json:"_id,omitempty"`
	Name  string `json:"name"`
	CTime int64  `json:"ctime"` // 创建时间
	MTime int64  `json:"mtime"` // 修改时间
}

type TestADto struct {
	Id        string    `json:"_id,omitempty"`
	Data      *TestBDto `json:"data"`
	Name      string    `json:"name"`
	TestFloat float64   `json:"test_float"`
	cTime     int64     `json:"ctime"` // 创建时间
	MTime     int64     `json:"mtime"` // 修改时间
}

var res = make(map[string]string)
var fieldProcess = func(valField reflect.Value, typField reflect.StructField) {
	key := GetItemNameByTag(typField, "json")
	println(key)
	value := ""
	switch valField.Kind() {
	case reflect.Invalid:
		value = "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		value = strconv.FormatInt(valField.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value = strconv.FormatUint(valField.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		value = strconv.FormatBool(valField.Bool())
	case reflect.String:
		value = strconv.Quote(valField.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		//value = valField.Type().String() + " 0x" +
		//	strconv.FormatUint(uint64(valField.Pointer()), 16)
		value = JsonFmt(valField.Interface())
	default: // reflect.Array, reflect.Struct, reflect.Interface
		//value = valField.Type().String() + " value"
		value = JsonFmt(valField.Interface())
	}
	res[key] = value
}

func TestRangeFields(t *testing.T) {

	test := &TestADto{
		Id: "111",
		Data: &TestBDto{
			Id:    "222",
			Name:  "222name",
			CTime: 7777,
			MTime: 7777,
		},
		Name:  "111name",
		cTime: 666,
		MTime: 666,
	}

	RangeFields(test, "json", fieldProcess)

	fmt.Println(res)
}

func JsonFmt(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}
