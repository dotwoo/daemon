package daemon

import (
	"os"
	"sync"
)

// File Handler
type FileHandler struct {
	fp    *os.File
	guard sync.Mutex
	path  string
	mode  os.FileMode
}

func NewFileHandler(path string, mode os.FileMode) *FileHandler {
	handler := &FileHandler{
		path: path,
		mode: mode,
	}
	handler.Reopen()

	return handler
}

func (h *FileHandler) Close() error {
	h.guard.Lock()
	defer h.guard.Unlock()

	return h.fp.Close()
}

func (h *FileHandler) Write(p []byte) (n int, err error) {
	h.guard.Lock()
	defer h.guard.Unlock()

	return h.fp.Write(p)
}

func (h *FileHandler) Reopen() {
	h.guard.Lock()
	defer h.guard.Unlock()

	if h.fp != nil {
		h.fp.Close()
	}

	fp, err := os.OpenFile(
		h.path,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		h.mode)

	if err != nil {
		panic(err)
	}

	h.fp = fp
}
