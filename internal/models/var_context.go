package models

type VarContext struct {
	parent *VarContext
	source *VarSource
}

func NewRootVarContext(sources ...*VarSource) *VarContext {
	var merged *VarSource
	for _, src := range sources {
		if src == nil || src.Len() == 0 {
			continue
		}
		if merged == nil {
			merged = src
		} else {
			merged = merged.Merge(src, true)
		}
	}

	if merged == nil {
		merged = &VarSource{}
	}

	return &VarContext{
		parent: nil,
		source: merged,
	}
}

func (c *VarContext) Child(source *VarSource) *VarContext {
	return &VarContext{
		parent: c,
		source: source,
	}
}

func (c *VarContext) Get(key string) (string, bool) {
	if c.source != nil {
		if val, ok := c.source.Get(key); ok {
			return val, true
		}
	}

	if c.parent != nil {
		return c.parent.Get(key)
	}

	return "", false
}

func (c *VarContext) Flatten() *VarSource {
	if c.parent == nil {
		if c.source == nil {
			return &VarSource{}
		}
		return c.source
	}

	parentFlat := c.parent.Flatten()

	if c.source == nil || c.source.Len() == 0 {
		return parentFlat
	}

	return parentFlat.Merge(c.source, true)
}

func (c *VarContext) GetAllSource() map[string]string {
	flat := c.Flatten()
	return flat.GetSource()
}
