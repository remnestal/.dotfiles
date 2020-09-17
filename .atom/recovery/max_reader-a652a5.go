package iohelper

import (
	"errors"
	"io"
)

// ErrReadQuotaExceeded indicates that a max reader was requested to read more
// than it was permitted to read.
var ErrReadQuotaExceeded error = errors.New("read quota exceeded")

// MaxReader is a reader that reads until either it's underlying data stream has
// been exhausted, the requested amount of bytes were read or it's read-quota
// has been exceeded. This is useful for setting an upper bound on reading to
// avoid costly processing of entities that are too large.
type MaxReader struct {
	R io.Reader // Underlying reader.
	N int64     // Read quota remaining.
}

func (mr *MaxReader) Read(p []byte) (int, error) {
	if mr.N <= 0 {
		return 0, ErrReadQuotaExceeded
	}
	if int64(len(p)) > mr.N {
		p = p[0:mr.N]
	}
	n, err := mr.R.Read(p)
	mr.N -= int64(n)
	return n, err
}
