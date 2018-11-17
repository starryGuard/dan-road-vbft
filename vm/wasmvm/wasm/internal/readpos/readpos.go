package readpos

import (
	"io"
)

// ReadPos implements io.Reader and stores the current number of bytes read from
// the reader
type ReadPos struct {
	R      io.Reader
	CurPos int64
}

// Read implements the io.Reader interface
func (r *ReadPos) Read(p []byte) (int, error) {
	n, err := r.R.Read(p)
	r.CurPos += int64(n)
	return n, err
}

// ReadByte implements the io.ByteReader interface
func (r *ReadPos) ReadByte() (byte, error) {
	p := make([]byte, 1)
	_, err := r.R.Read(p)
	return p[0], err
}
