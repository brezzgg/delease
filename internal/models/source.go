package models

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Source interface {
	Len() int
	yaml.Unmarshaler
	json.Marshaler
}
