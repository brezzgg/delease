package models

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type YamlSliceSource[V any] struct {
	data []V
}

func (s *YamlSliceSource[V]) GetSource() []V {
	return s.data
}

func (s *YamlSliceSource[V]) GetSourceCopy() []V {
	return SliceCopy(s.data)
}

func (s *YamlSliceSource[V]) SetSource(src []V) {
	s.data = src
}

func (s *YamlSliceSource[V]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

func (s *YamlSliceSource[V]) UnmarshalYAML(value *yaml.Node) error {
	s.data = make([]V, 0)
	return SliceUnmarshal(&s.data, value)
}

func (s *YamlSliceSource[V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

var _ Source[[]any] = (*YamlSliceSource[any])(nil)
