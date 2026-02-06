package models_test

import (
	"testing"

	"github.com/brezzgg/delease/internal/models"
)

func TestCommand_ApplyVars(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		vars    map[string]string
		want    string
		wantErr bool
	}{
		{
			name: "similiar vars",
			raw:  "rm ${{log}} ./logs/${{log}}${{log}}",
			vars: map[string]string{
				"log": "file.log",
			},
			want:    "rm file.log ./logs/file.logfile.log",
			wantErr: false,
		},
		{
			name:    "vars len = 0",
			raw:     "hello world",
			vars:    map[string]string{},
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "vars is nil",
			raw:     "hello world",
			vars:    nil,
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "success",
			raw:     "${{var1}} ${{var2}}${{var1}} ${{var3}}${{var2}}",
			vars:    map[string]string{
				"var1": "val1",
				"var2": "val2",
				"var3": "val3",
				"var4": "val4",
			},
			want:    "val1 val2val1 val3val2",
			wantErr: false,
		},
		{
			name:    "inside {}",
			raw:     "{${{var}} }",
			vars:    map[string]string{
				"var": "val",
			},
			want:    "{val }",
			wantErr: false,
		},
		{
			name:    "no required",
			raw:     "cp ${{from}} ${{to}}",
			vars:    map[string]string{
				"from": "file",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &models.Command{
				Raw: tt.raw,
			}
			if err := c.ParseVars(); err != nil {
				t.Errorf("failed to parse vars: %s", err.Error())
			}
			varSrc := &models.VarSource{}
			varSrc.SetSource(tt.vars)
			got, gotErr := c.ApplyVars(varSrc)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ApplyVars() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ApplyVars() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ApplyVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
