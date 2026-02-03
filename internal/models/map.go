package models

import "reflect"

func cleanMap[K comparable, V any](m map[K]V) {
	if m == nil {
		return
	}
	var keys []K
	for k, v := range m {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
			if rv.IsNil() {
				keys = append(keys, k)
			}
		}
	}
	for _, key := range keys {
		delete(m, key)
	}
}
