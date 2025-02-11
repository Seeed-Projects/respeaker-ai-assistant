package main

import (
	"bytes"
	"fmt"
	"io"
)

type nopSeek struct {
	*bytes.Buffer
}

func (n *nopSeek) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart && offset == 0 {
		return 0, nil
	}
	return 0, fmt.Errorf("seek not implemented")
}
