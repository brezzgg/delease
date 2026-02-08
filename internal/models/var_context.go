package models

type VarContext struct {
	parent       *VarContext
	source       *VarSource
	example bool
}

func NewVarContext(src ...*VarSource) *VarContext {
	var merged *VarSource
	for _, src := range src {
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
		parent:       nil,
		source:       merged,
		example: false,
	}
}

func NewExampleVarContext(src ...*VarSource) *VarContext {
	c := NewVarContext()
	c.example = true
	return c
}

func (c *VarContext) Child(src *VarSource) *VarContext {
	return &VarContext{
		parent: c,
		source: src,
	}
}

func (c *VarContext) Get(key string, varType VarType) (string, bool) {
	if c.example && (varType == VarTypeDynamic || varType == VarTypeOs)  {
		return key, true
	}

	if c.source != nil {
		if val, ok := c.source.Get(key); ok {
			if val.IsType(varType) {
				return val.Content, true
			}
		}
	}

	if c.parent != nil {
		return c.parent.Get(key, varType)
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

func (c *VarContext) GetAllSource() map[string]*Var {
	flat := c.Flatten()
	return flat.GetSource()
}
