package models

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type YamlMapSource[V any] struct {
	data map[string]V
}

func (e *YamlMapSource[V]) SetSource(v map[string]V) {
	e.data = v
}

func (e *YamlMapSource[V]) Get(key string) (V, bool) {
	return MapGet(e.data, key)
}

func (e *YamlMapSource[V]) Keys() []string {
	return MapKeys(e.data)
}

func (e *YamlMapSource[V]) GetMap() map[string]V {
	return e.data
}

func (e *YamlMapSource[V]) GetMapCopy() map[string]V {
	return MapCopy(e.data)
}

func (e *YamlMapSource[V]) Merge(oth YamlMapSource[V], force bool) YamlMapSource[V] {
	return YamlMapSource[V]{data: MapMerge(e.data, oth.GetMap(), force)}
}

func (e *YamlMapSource[V]) Len() int {
	if e == nil {
		return 0
	}
	return len(e.data)
}

func (e *YamlMapSource[V]) Set(key string, val V, force bool) {
	MapSet(e.data, key, val, force)
}

func (e *YamlMapSource[V]) UnmarshalYAML(value *yaml.Node) error {
	e.data = make(map[string]V, 0)
	return MapUnmarshal(&e.data, value)
}

func (e *YamlMapSource[V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.data)
}

var _ Source = (*YamlMapSource[any])(nil)
