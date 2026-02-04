package models

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type YamlSliceSource[V any] struct {
	data []V
}

func (s *YamlSliceSource[V]) SetSource(v []V) {
	s.data = v
}

func (s *YamlSliceSource[V]) Get() []V {
	return s.data
}

func (s *YamlSliceSource[V]) GetCopy() []V {
	return SliceCopy(s.data)
}

func (s *YamlSliceSource[V]) Len() int {
	return len(s.data)
}

func (s *YamlSliceSource[V]) UnmarshalYAML(value *yaml.Node) error {
	s.data = make([]V, 0)
	return SliceUnmarshal(&s.data, value)
}

func (s *YamlSliceSource[V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

var _ Source = (*YamlSliceSource[any])(nil)
