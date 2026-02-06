package exec

import (
	"bytes"
	"strings"
	"sync"

	"github.com/brezzgg/delease/internal/exec/decoder"
)

type SyncWriter struct {
	mu  sync.Mutex
	log Logger
	t   MsgType
	buf bytes.Buffer
}

func NewSyncWriter(log Logger, t MsgType) *SyncWriter {
	return &SyncWriter{
		log: log,
		t:   t,
	}
}

func (w *SyncWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.buf.Write(w.decodeOutput(p))

	for {
		str := w.buf.String()
		idx := strings.Index(str, "\n")
		if idx == -1 {
			break
		}

		line := str[:idx+1]
		w.log(line, w.t)

		w.buf.Reset()
		w.buf.WriteString(str[idx+1:])
	}

	return len(p), nil
}

func (w *SyncWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.buf.Len() > 0 {
		w.log(w.buf.String(), w.t)
		w.buf.Reset()
	}
}

func (w *SyncWriter) decodeOutput(data []byte) []byte {
	return decoder.Decode(data)
}
