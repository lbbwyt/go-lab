package array

import "reflect"

func InArray(obj interface{}, target interface{}) bool {
	if target == nil {
		return false
	}

	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// ReverseArrayString reverse the slice's element from tail to head.
func ReverseArrayString(t []string) []string {
	if len(t) == 0 {
		return t
	}
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
	return t
}
