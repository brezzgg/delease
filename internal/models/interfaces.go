package models

type Applier[T any] interface {
	ApplyVars(vars *VarSource) (T, error)
}

type Mergeable[T any] interface {
	Merge(oth T, force bool) T
}

