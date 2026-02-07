package models

type Mergeable[T any] interface {
	Merge(oth T, force bool) T
}
