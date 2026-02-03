package exec_test

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/brezzgg/delease/internal/exec"
)

// Mock Logger для тестов
type mockLogger struct {
	mu       sync.Mutex
	messages []logMessage
}

type logMessage struct {
	text    string
	msgType exec.MsgType
}

func (m *mockLogger) Log(text string, msgType exec.MsgType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, logMessage{text: text, msgType: msgType})
}

func (m *mockLogger) GetMessages() []logMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]logMessage{}, m.messages...)
}

func (m *mockLogger) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = nil
}

func setupSh(t *testing.T, lines, env []string) (*exec.Sh, *mockLogger) {
	t.Helper()

	sh := &exec.Sh{}
	logger := &mockLogger{}

	err := sh.Setup("/tmp", lines, env, logger.Log)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	return sh, logger
}

func TestSh_Setup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wd      string
		lines   []string
		env     []string
		wantErr bool
	}{
		{
			name:    "valid setup",
			wd:      "/tmp",
			lines:   []string{"echo hello", "pwd"},
			env:     []string{"VAR=value"},
			wantErr: false,
		},
		{
			name:    "empty lines",
			wd:      "/tmp",
			lines:   []string{},
			env:     []string{},
			wantErr: false,
		},
		{
			name:    "with multiple env vars",
			wd:      ".",
			lines:   []string{"echo $VAR1 $VAR2"},
			env:     []string{"VAR1=hello", "VAR2=world"},
			wantErr: false,
		},
		{
			name:    "relative working directory",
			wd:      "..",
			lines:   []string{"pwd"},
			env:     []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sh := &exec.Sh{}
			logger := &mockLogger{}

			err := sh.Setup(tt.wd, tt.lines, tt.env, logger.Log)

			if (err != nil) != tt.wantErr {
				t.Errorf("Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSh_RunLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		lines        []string
		lineToRun    int
		env          []string
		wantCode     int
		wantError    bool
		wantStdout   string
		wantStderr   string
		checkMessage bool
	}{
		{
			name:         "successful echo",
			lines:        []string{"echo hello"},
			lineToRun:    0,
			wantCode:     0,
			wantError:    false,
			wantStdout:   "hello\n",
			checkMessage: true,
		},
		{
			name:      "exit with code 5",
			lines:     []string{"exit 5"},
			lineToRun: 0,
			wantCode:  5,
			wantError: false,
		},
		{
			name:      "false command",
			lines:     []string{"false"},
			lineToRun: 0,
			wantCode:  1,
			wantError: false,
		},
		{
			name:      "true command",
			lines:     []string{"true"},
			lineToRun: 0,
			wantCode:  0,
			wantError: false,
		},
		{
			name:      "line out of bounds - negative",
			lines:     []string{"echo test"},
			lineToRun: -1,
			wantCode:  0,
			wantError: true,
		},
		{
			name:      "line out of bounds - too large",
			lines:     []string{"echo test"},
			lineToRun: 10,
			wantCode:  0,
			wantError: true,
		},
		{
			name:       "stderr output",
			lines:      []string{"echo error >&2"},
			lineToRun:  0,
			wantCode:   0,
			wantError:  false,
			wantStderr: "error\n",
		},
		{
			name:         "environment variable",
			lines:        []string{"echo $TEST_VAR"},
			lineToRun:    0,
			env:          []string{"TEST_VAR=hello_world"},
			wantCode:     0,
			wantError:    false,
			wantStdout:   "hello_world\n",
			checkMessage: true,
		},
		{
			name:      "invalid syntax",
			lines:     []string{"echo 'unclosed"},
			lineToRun: 0,
			wantError: true,
		},
		{
			name:         "multiline command",
			lines:        []string{"echo line1\necho line2"},
			lineToRun:    0,
			wantCode:     0,
			wantError:    false,
			wantStdout:   "line1",
			checkMessage: true,
		},
		{
			name:      "command not found",
			lines:     []string{"/nonexistent_command_xyz_12345"},
			lineToRun: 0,
			wantCode:  127,
			wantError: false,
		},
		{
			name:         "pipe commands",
			lines:        []string{"echo hello world | grep hello"},
			lineToRun:    0,
			wantCode:     0,
			wantError:    false,
			wantStdout:   "hello world\n",
			checkMessage: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sh, logger := setupSh(t, tt.lines, tt.env)

			result := sh.RunLine(context.Background(), tt.lineToRun)

			if (result.Error != nil) != tt.wantError {
				t.Errorf("RunLine() error = %v, wantError %v", result.Error, tt.wantError)
			}

			if !tt.wantError && result.Code != tt.wantCode {
				t.Errorf("RunLine() code = %d, want %d", result.Code, tt.wantCode)
			}

			if tt.checkMessage {
				messages := logger.GetMessages()
				found := false
				for _, msg := range messages {
					if strings.Contains(msg.text, strings.TrimSpace(tt.wantStdout)) {
						found = true
						if msg.msgType != exec.MsgTypeStdout {
							t.Errorf("Expected stdout message, got type %v", msg.msgType)
						}
						break
					}
				}
				if !found && tt.wantStdout != "" {
					t.Errorf("Expected stdout containing %q, got messages: %+v", tt.wantStdout, messages)
				}
			}

			if tt.wantStderr != "" {
				messages := logger.GetMessages()
				found := false
				for _, msg := range messages {
					if strings.Contains(msg.text, strings.TrimSpace(tt.wantStderr)) &&
						msg.msgType == exec.MsgTypeStderr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected stderr containing %q, got messages: %+v", tt.wantStderr, messages)
				}
			}
		})
	}
}

func TestSh_Run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		lines          []string
		env            []string
		wantFinalCode  int
		wantFinalError bool
	}{
		{
			name:           "all lines succeed",
			lines:          []string{"echo line1", "echo line2", "echo line3"},
			wantFinalCode:  0,
			wantFinalError: false,
		},
		{
			name:           "stop on error",
			lines:          []string{"echo line1", "exit 5", "echo line3"},
			wantFinalCode:  5,
			wantFinalError: false,
		},
		{
			name:           "stop on false",
			lines:          []string{"true", "false", "echo should not run"},
			wantFinalCode:  1,
			wantFinalError: false,
		},
		{
			name:           "empty lines",
			lines:          []string{},
			wantFinalCode:  0,
			wantFinalError: false,
		},
		{
			name:           "single line success",
			lines:          []string{"echo hello"},
			wantFinalCode:  0,
			wantFinalError: false,
		},
		{
			name:           "single line failure",
			lines:          []string{"exit 42"},
			wantFinalCode:  42,
			wantFinalError: false,
		},
		{
			name: "complex script",
			lines: []string{
				"export VAR=test",
				"echo $VAR",
				"true",
			},
			wantFinalCode:  0,
			wantFinalError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sh, _ := setupSh(t, tt.lines, tt.env)

			resChan := make(chan exec.Result, 1)
			ctx := context.Background()

			go sh.Run(ctx, resChan)

			result := <-resChan

			if (result.Error != nil) != tt.wantFinalError {
				t.Errorf("Run() error = %v, wantError %v", result.Error, tt.wantFinalError)
			}

			if !tt.wantFinalError && result.Code != tt.wantFinalCode {
				t.Errorf("Run() final code = %d, want %d", result.Code, tt.wantFinalCode)
			}
		})
	}
}
