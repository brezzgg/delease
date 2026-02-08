package cmd

import (
	"encoding/json"
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

var (
	verbose bool
	indent  bool
)

func configureLogger(verbose bool) {
	var level lg.LoggerOption
	if verbose {
		level = lg.WithLogLevel(lg.LogLevelDebug)
	} else {
		level = lg.WithLogLevel(lg.LogLevelInfo)
	}

	lg.GlobalLogger = lg.NewLogger(
		level,
		lg.WithPipe(lg.NewPipe(lg.WithSerializer(&serializer{}))),
	)
}

type serializer struct{}

func (s *serializer) Serialize(m lg.Message) string {
	sb := strings.Builder{}

	if m.Level.Level != lg.LogLevelInfo.Level {
		sb.WriteString(strings.ToLower(m.Level.Level))
		sb.WriteString(": ")
	}
	sb.WriteString(m.Text)

	if m.Context != nil {
		sb.WriteString(" ")
		var (
			buf []byte
			err error
		)
		if indent {
			buf, err = json.MarshalIndent(m.Context, "", "  ")
		} else {
			buf, err = json.Marshal(m.Context)
		}
		if err != nil {
			sb.WriteString(`{"serialize_error": "failed to serialize message"}`)
		} else {
			if indent {
				sb.WriteString("\n")
				sb.WriteString(string(buf))
			} else {
				sb.WriteString(string(buf))
			}
		}
	}
	return sb.String()
}

var _ lg.Serializer = (*serializer)(nil)
