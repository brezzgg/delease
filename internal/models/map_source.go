package models

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type YamlMapSource[V any] struct {
	data map[string]V
}

func (s *YamlMapSource[V]) GetSource() map[string]V {
	return s.data
}

func (s *YamlMapSource[V]) GetSourceCopy() map[string]V {
	return MapCopy(s.data)
}

func (s *YamlMapSource[V]) SetSource(v map[string]V) {
	s.data = v
}

func (s *YamlMapSource[V]) Get(key string) (V, bool) {
	return MapGet(s.data, key)
}

func (s *YamlMapSource[V]) Keys() []string {
	return MapKeys(s.data)
}

func (s *YamlMapSource[V]) Merge(oth YamlMapSource[V], force bool) YamlMapSource[V] {
	return YamlMapSource[V]{data: MapMerge(s.data, oth.GetSource(), force)}
}

func (s *YamlMapSource[V]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

func (s *YamlMapSource[V]) Set(key string, val V, force bool) {
	MapSet(s.data, key, val, force)
}

func (s *YamlMapSource[V]) UnmarshalYAML(value *yaml.Node) error {
	s.data = make(map[string]V, 0)
	return MapUnmarshal(&s.data, value)
}

func (s *YamlMapSource[V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

var _ Source[map[string]any] = (*YamlMapSource[any])(nil)
