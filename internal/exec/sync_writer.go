package exec

import (
	"strings"
	"sync"
)

type SyncWriter struct {
	mu  sync.Mutex
	log Logger
	t   MsgType
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
	str := string(p)
	if len(strings.TrimSpace(str)) != 0 {
		w.log(str, w.t)
	}
	return len(p), nil
}
