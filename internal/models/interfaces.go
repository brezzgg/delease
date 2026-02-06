package models

type Applier[T any] interface {
	ApplyVars(vars *VarSource) (T, error)
	Applied() bool
}

type Mergeable[T any] interface {
	Merge(oth T, force bool) T
}
