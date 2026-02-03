package parser_test

import (
	"reflect"
	"testing"

	"github.com/brezzgg/delease/internal/parser"
	"github.com/brezzgg/delease/internal/models"
)

func TestReplaceCmd(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		vars    models.Vars
		want    string
		wantErr bool
	}{
		{
			name:    "no submatch",
			cmd:     "hello world",
			vars:    make(map[string]string),
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "nil vars & no submatch",
			cmd:     "hello world",
			vars:    nil,
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "nil vars & submatch",
			cmd:     "hello {{somevar}}",
			vars:    nil,
			want:    "",
			wantErr: true,
		},
		{
			name: "no var",
			cmd:  "hello {{somevar}}",
			vars: map[string]string{
				"somevar2": "world",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no var 2",
			cmd:  "{{somevar2}} {{somevar}}",
			vars: map[string]string{
				"somevar2": "world",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success 1",
			cmd:  "hello {{somevar}}",
			vars: map[string]string{
				"somevar": "world",
			},
			want:    "hello world",
			wantErr: false,
		},
		{
			name: "success 2",
			cmd:  "hello {{ somevar }}{{ v2 }}",
			vars: map[string]string{
				"somevar": "world",
				"v2":      "!!!",
			},
			want:    "hello world!!!",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parser.ReplaceCmd(tt.cmd, tt.vars)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReplaceCmd() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReplaceCmd() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ReplaceCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testMergeVarsStruct struct {
	m models.Vars
}

func (t *testMergeVarsStruct) GetVars() models.Vars {
	return t.m
}

func TestMergeVariablers(t *testing.T) {
	tests := []struct {
		name  string
		l     models.Variabler
		r     models.Variabler
		force bool
		want  models.Vars
	}{
		{
			name:  "2 nil",
			l:     nil,
			r:     nil,
			force: false,
			want:  models.Vars{},
		},
		{
			name: "l nil",
			l:    nil,
			r: &testMergeVarsStruct{m: map[string]string{
				"1": "11",
				"2": "22",
			}},
			force: false,
			want: models.Vars{
				"1": "11",
				"2": "22",
			},
		},
		{
			name: "r nil",
			r:    nil,
			l: &testMergeVarsStruct{m: map[string]string{
				"1": "11",
				"2": "22",
			}},
			force: false,
			want: models.Vars{
				"1": "11",
				"2": "22",
			},
		},
		{
			name: "success no replace",
			l: &testMergeVarsStruct{m: map[string]string{
				"1": "11",
				"2": "22",
			}},
			r: &testMergeVarsStruct{m: map[string]string{
				"3": "33",
				"4": "44",
			}},
			force: false,
			want: models.Vars{
				"1": "11",
				"2": "22",
				"3": "33",
				"4": "44",
			},
		},
		{
			name: "success force",
			l: &testMergeVarsStruct{m: map[string]string{
				"1": "11",
				"2": "22",
			}},
			r: &testMergeVarsStruct{m: map[string]string{
				"1": "33",
				"2": "44",
			}},
			force: true,
			want: models.Vars{
				"1": "33",
				"2": "44",
			},
		},
		{
			name: "success !force",
			l: &testMergeVarsStruct{m: map[string]string{
				"1": "11",
				"2": "22",
			}},
			r: &testMergeVarsStruct{m: map[string]string{
				"1": "33",
				"2": "44",
			}},
			force: false,
			want: models.Vars{
				"1": "11",
				"2": "22",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.MergeVariablers(tt.l, tt.r, tt.force)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeVariablers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyVars(t *testing.T) {
	tests := []struct {
		name      string
		root      *models.Root
		task      string
		forceVars models.Vars
		wantErr   bool
		wantCmds  []string
	}{
		{
			name: "success",
			root: &models.Root{
				Vars: models.Vars{
					"user": "admin",
				},
				Tasks: map[string]*models.Task{
					"deploy": {
						Vars: models.Vars{
							"env": "prod",
						},
						Cmds: []string{
							"echo {{user}}",
							"deploy --env={{env}}",
						},
					},
				},
			},
			task: "deploy",
			forceVars: models.Vars{
				"version": "1.0.0",
			},
			wantErr: false,
			wantCmds: []string{
				"echo admin",
				"deploy --env=prod",
			},
		},
		{
			name: "forceVars replace",
			root: &models.Root{
				Vars: models.Vars{
					"env": "dev",
				},
				Tasks: map[string]*models.Task{
					"build": {
						Vars: models.Vars{
							"env": "staging",
						},
						Cmds: []string{
							"build --env={{env}}",
						},
					},
				},
			},
			task: "build",
			forceVars: models.Vars{
				"env": "prod",
			},
			wantErr: false,
			wantCmds: []string{
				"build --env=prod",
			},
		},
		{
			name: "task without cmds",
			root: &models.Root{
				Tasks: map[string]*models.Task{
					"empty": {
						Cmds: []string{},
					},
				},
			},
			task:     "empty",
			forceVars: models.Vars{},
			wantErr:  false,
			wantCmds: []string{},
		},
		{
			name: "no tasks",
			root: &models.Root{
				Tasks: map[string]*models.Task{},
			},
			task:      "test",
			forceVars: models.Vars{},
			wantErr:   true,
		},
		{
			name: "task not found",
			root: &models.Root{
				Tasks: map[string]*models.Task{
					"build": {
						Cmds: []string{"go build"},
					},
				},
			},
			task:      "deploy",
			forceVars: models.Vars{},
			wantErr:   true,
		},
		{
			name: "success 2",
			root: &models.Root{
				Vars: models.Vars{
					"user":    "admin",
					"project": "myapp",
				},
				Tasks: map[string]*models.Task{
					"deploy": {
						Vars: models.Vars{
							"env":  "prod",
							"port": "8080",
						},
						Cmds: []string{
							"echo Deploying {{project}}",
							"docker run -p {{port}}:{{port}} {{user}}/{{project}}",
							"echo Deployed to {{env}}",
						},
					},
				},
			},
			task: "deploy",
			forceVars: models.Vars{
				"region": "us-east-1",
			},
			wantErr: false,
			wantCmds: []string{
				"echo Deploying myapp",
				"docker run -p 8080:8080 admin/myapp",
				"echo Deployed",
			},
		},
		{
			name: "task vars replace",
			root: &models.Root{
				Vars: models.Vars{
					"env": "dev",
					"db":  "postgres",
				},
				Tasks: map[string]*models.Task{
					"migrate": {
						Vars: models.Vars{
							"env": "prod",
						},
						Cmds: []string{
							"migrate --env={{env}} --db={{db}}",
						},
					},
				},
			},
			task:      "migrate",
			forceVars: models.Vars{},
			wantErr:   false,
			wantCmds: []string{
				"migrate --env=prod --db=postgres",
			},
		},
		{
			name: "empty vars",
			root: &models.Root{
				Tasks: map[string]*models.Task{
					"simple": {
						Cmds: []string{
							"echo hello world",
							"ls -la",
						},
					},
				},
			},
			task:      "simple",
			forceVars: models.Vars{},
			wantErr:   false,
			wantCmds: []string{
				"echo hello world",
				"ls -la",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.ApplyVars(tt.root, tt.task, tt.forceVars)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ApplyVars() want err, except nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ApplyVars() unexpected error = %v", err)
				return
			}

			if tt.wantCmds != nil {
				gotCmds := tt.root.Tasks[tt.task].Cmds
				if len(gotCmds) != len(tt.wantCmds) {
					t.Errorf("ApplyVars() = %d, want %d", len(gotCmds), len(tt.wantCmds))
					return
				}
				for i, want := range tt.wantCmds {
					if gotCmds[i] != want {
						t.Errorf("ApplyVars() cmd[%d] = %q, want %q", i, gotCmds[i], want)
					}
				}
			}
		})
	}
}
