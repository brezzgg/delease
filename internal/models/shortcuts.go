package models

import (
	"errors"
	"maps"
	"reflect"

	"gopkg.in/yaml.v3"
)

func MapSet[K comparable, V any](m map[K]V, key K, val V, force bool) {
	if _, ok := m[key]; ok {
		if force {
			m[key] = val
		}
	} else {
		m[key] = val
	}
}

func MapGet[K comparable, V any](m map[K]V, key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return []K{}
	}
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func MapMerge[K comparable, V any](m, oth map[K]V, force bool) map[K]V {
	if m == nil && oth == nil {
		return make(map[K]V)
	}
	if m == nil && oth != nil {
		return MapCopy(oth)
	}
	if m != nil && oth == nil {
		return MapCopy(m)
	}
	res := make(map[K]V, len(m)+len(oth))
	maps.Copy(res, m)
	for k, v := range oth {
		MapSet(res, k, v, force)
	}
	return res
}

func MapCopy[K comparable, V any](m map[K]V) map[K]V {
	r := make(map[K]V, len(m))
	maps.Copy(r, m)
	return r
}

func MapClean[K comparable, V any](m map[K]V) {
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

func MapUnmarshal[K comparable, V any](dst *map[K]V, value *yaml.Node) error {
	if dst == nil {
		return errors.New("dst is nil")
	}
	if err := value.Decode(dst); err != nil {
		return err
	}
	return nil
}

func SliceUnmarshal[V any](dst *[]V, value *yaml.Node) error {
	if dst == nil {
		return errors.New("dst is nil")
	}
	if err := value.Decode(dst); err != nil {
		return err
	}
	return nil
}

func SliceCopy[V any](s []V) []V {
	r := make([]V, len(s))
	copy(r, s)
	return r
}

// PremergeCheck checks both pointers for nil; if additional logic is required, it returns nil.
func PremergeCheck[T any](left, right *T) *T {
	if left == nil {
		if right == nil {
			var zero T
			return &zero
		}
		return right
	}
	if right == nil {
		return left
	}

	return nil
}
