package models

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Source[T any] interface {
	Len() int
	SetSource(src T)
	GetSource() T
	GetSourceCopy() T
	yaml.Unmarshaler
	json.Marshaler
}
