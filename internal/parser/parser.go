package parser

import (
	"bytes"

	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

func Parse(b []byte) (*models.Root, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(b))
	decoder.KnownFields(true)

	root := &models.Root{}

	if err := decoder.Decode(root); err != nil {
		return nil, lg.Ef("unmarshal error: %w", err)
	}

	return root, nil
}
