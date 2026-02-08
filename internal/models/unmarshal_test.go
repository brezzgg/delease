package models_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/brezzgg/delease/internal/models"
	"gopkg.in/yaml.v3"
)

func Test_UnmarshalTest(t *testing.T) {
	var result *models.Root
	t.Run("decode", func(t *testing.T) {
		decoder := yaml.NewDecoder(bytes.NewReader([]byte(file)))
		decoder.KnownFields(true)
		result = &models.Root{}
		if err := decoder.Decode(result); err != nil {
			t.Fatalf("decode error: %s", err.Error())
		}
	})

	t.Run("root.var", func(t *testing.T) {
		if m := map[string]*models.Var{
			"var1": models.NewVarT("val 1"),
			"var2": models.NewVarT("val 2"),
		}; !reflect.DeepEqual(m, result.Var.GetSource()) {
			t.Errorf("root.Var = %v, want = %v", m, result.Var.GetSource())
		}
	})

	t.Run("root.env", func(t *testing.T) {
		if m := map[string]string{
			"env1": "env 1",
			"env2": "env 2",
		}; !reflect.DeepEqual(m, result.Env.GetSource()) {
			t.Errorf("root.Var = %v, want = %v", m, result.Env.GetSource())
		}
	})

	t.Run("root.include", func(t *testing.T) {
		if a := []string{"inc1", "inc2"}; !reflect.DeepEqual(a, result.Include.GetSource()) {
			t.Errorf("root.Var = %v, want = %v", a, result.Include.GetSource())
		}
	})

	tasks := result.Tasks.GetSource()

	t.Run("tasks len check", func(t *testing.T) {
		if len(tasks) != 5 {
			t.Fatalf("len(root.tasks) = %d, want %d", len(tasks), 5)
		}
	})
	if len(tasks) != 5 {
		return
	}

	t.Run("task 1 check", func(t *testing.T) {
		v := tasks["task1"]
		if m := map[string]*models.Var{
			"var2": models.NewVarT("val 4"),
		}; !reflect.DeepEqual(m, v.Vars.GetSource()) {
			t.Errorf("task1.Var = %v, want = %v", m, v.Vars.GetSource())
		}

		c := make([]*models.Command, 2)
		c[0] = &models.Command{Raw: "cmd1"}
		c[1] = &models.Command{Raw: "cmd2"}
		if !reflect.DeepEqual(c, v.Cmds.GetSource()) {
			t.Errorf("task1.Var = %v, want = %v", v.Cmds.GetSource(), c)
		}
	})

	t.Run("task 2 check", func(t *testing.T) {
		v := tasks["task2"]
		if m := map[string]string{
			"env1": "env 1",
		}; !reflect.DeepEqual(m, v.Env.GetSource()) {
			t.Errorf("task2.Var = %v, want = %v", m, v.Env.GetSource())
		}

		c := make([]*models.Command, 2)
		c[0] = &models.Command{Raw: "cmd1"}
		c[1] = &models.Command{Raw: "cmd2"}
		if !reflect.DeepEqual(c, v.Cmds.GetSource()) {
			t.Errorf("task2.Var = %v, want = %v", v.Cmds.GetSource(), c)
		}
	})

	t.Run("task 3 check", func(t *testing.T) {
		v := tasks["task3"]

		c := make([]*models.Command, 2)
		c[0] = &models.Command{Raw: "cmd1 \\\n  -v\ncmd2\n"}
		if !reflect.DeepEqual(c[0].Raw, v.Cmds.GetSource()[0].Raw) {
			t.Errorf("task3.cmd\nexpc = %v\nwant = %v", v.Cmds.GetSource()[0].Raw, c[0].Raw)
		}
	})

	t.Run("task 4 check", func(t *testing.T) {
		v := tasks["task4"]

		c := make([]*models.Command, 2)
		c[0] = &models.Command{Raw: "cmd1\ncmd2\n"}
		if !reflect.DeepEqual(c[0].Raw, v.Cmds.GetSource()[0].Raw) {
			t.Errorf("task4.cmd\nexpc = %v\nwant = %v", v.Cmds.GetSource()[0].Raw, c[0].Raw)
		}
	})

	t.Run("task 5 check", func(t *testing.T) {
		v := tasks["task5"]

		c := make([]*models.Command, 2)
		c[0] = &models.Command{Raw: "cmd1"}
		c[1] = &models.Command{Raw: "cmd2\n", Os: "windows"}
		if !reflect.DeepEqual(c[0].Raw, v.Cmds.GetSource()[0].Raw) {
			t.Errorf("task5.cmd\nexpc = %v\nwant = %v", v.Cmds.GetSource()[0].Raw, c[0].Raw)
		}
	})
}

var file = `
includes: [inc1, inc2]

vars:
  var1: val 1
  var2: val 2

envs:
  env1: env 1
  env2: env 2

tasks:
  task1:
    vars: 
      var2: val 4
    cmds: [cmd1, cmd2]
  task2:
    envs:
      env1: env 1
    cmds:
      - cmd1
      - cmd2
    before: [task1]
  task3:
    cmds: 
      - run: |
          cmd1 \
            -v
          cmd2
  task4:
    cmds: 
      - run: |
          cmd1
          cmd2
  task5:
    cmds:
      - run: cmd1
      - os: windows
        run: |
          cmd2
      - cmd3
  task6:
`
