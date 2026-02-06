package parser

import (
	"bytes"
	"strings"

	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	def, cwd string
	visited  map[string]bool
}

type RootWithPath struct {
	root *models.Root
	path string
}

func NewRootWithPath(root *models.Root, path string) (*RootWithPath, error) {
	return &RootWithPath{
		root: root,
		path: path,
	}, nil
}

func (h *RootWithPath) Equal(oth *RootWithPath) bool {
	return strings.EqualFold(h.path, oth.path)
}

func (h *RootWithPath) Root() *models.Root {
	return h.root
}

func New(def, cwd string) *Parser {
	return &Parser{
		def: def,
		cwd: cwd,
	}
}

func (p *Parser) Parse() (*models.Root, error) {
	p.visited = make(map[string]bool)
	return p.parse(p.def, p.cwd)
}

func (p *Parser) parse(def, cwd string) (*models.Root, error) {
	hroot, err := p.decode(def, cwd)
	if err != nil {
		return nil, lg.Ef("decode error: %w", err)
	}

	if p.visited[hroot.path] {
		return nil, lg.Ef("recursive include detected in file %s", def)
	}
	p.visited[hroot.path] = true

	currentRoot := hroot.Root()

	if currentRoot.Include != nil && currentRoot.Include.Len() > 0 {
		for _, includePath := range currentRoot.Include.GetCopy() {
			includedRoot, err := p.parse(includePath, cwd)
			if err != nil {
				return nil, err
			}
			currentRoot = currentRoot.Merge(includedRoot, false)
		}
	}

	return currentRoot, nil
}

func (p *Parser) decode(def, cwd string) (*RootWithPath, error) {
	b, abs, err := FindConfig(def, cwd)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(b))
	decoder.KnownFields(true)

	root := &models.Root{}

	if err = decoder.Decode(root); err != nil {
		return nil, err
	}

	hr, err := NewRootWithPath(root, abs)
	if err != nil {
		return nil, err
	}

	return hr, nil
}
