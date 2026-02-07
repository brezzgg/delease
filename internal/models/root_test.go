package models_test

import (
	"reflect"
	"testing"

	"github.com/brezzgg/delease/internal/models"
)

func TestIncludeSource_Merge(t *testing.T) {
	tests := []struct {
		name  string
		left  []string
		right []string
		force bool
		want  []string
	}{
		{
			name: "success",
			left: []string{
				"inc1",
				"inc2",
			},
			right: []string{
				"inc2",
				"inc3",
			},
			force: true,
			want: []string{
				"inc3",
				"inc1",
				"inc2",
			},
		},
		{
			name: "success 2",
			left: []string{
				"inc1",
				"inc2",
			},
			right: []string{
				"inc2",
				"inc3",
			},
			force: false,
			want: []string{
				"inc1",
				"inc2",
				"inc3",
			},
		},
		{
			name: "success 3",
			left: nil,
			right: []string{
				"inc2",
				"inc3",
			},
			force: false,
			want: []string{
				"inc2",
				"inc3",
			},
		},
		{
			name: "success 4",
			left: []string{
				"inc1",
				"inc2",
			},
			right: []string{},
			force: false,
			want: []string{
				"inc1",
				"inc2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lsrc := &models.IncludeSource{}
			lsrc.SetSource(tt.left)
			rsrc := &models.IncludeSource{}
			rsrc.SetSource(tt.right)
			got := lsrc.Merge(rsrc, tt.force)

			if !reflect.DeepEqual(got.GetSource(), tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
