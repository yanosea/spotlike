package proxy

import (
	"os"
)

// Os is an interface that provides a proxy of the methods of os.
type Os interface {
	Pipe() (File, File, error)
}

// osProxy is a proxy struct that implements the Os interface.
type osProxy struct{}

// NewOs returns a new instance of the Os interface.
func NewOs() Os {
	return &osProxy{}
}

// Pipe creates a synchronous in-memory pipe.
func (osProxy) Pipe() (File, File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return &fileProxy{read}, &fileProxy{write}, nil
}

// File is an interface that provides a proxy of the methods of os.File.
type File interface {
	AsOsFile() *os.File
	Close() error
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

// fileProxy is a proxy struct that implements the File interface.
type fileProxy struct {
	file *os.File
}

// This method allows fileProxy to be type asserted to *os.File
func (f *fileProxy) AsOsFile() *os.File {
	return f.file
}

// Close closes the File, rendering it unusable for I/O.
func (f *fileProxy) Close() error {
	return f.file.Close()
}

// Read reads up to len(p) bytes into p.
func (f *fileProxy) Read(b []byte) (n int, err error) {
	return f.file.Read(b)
}

// Write writes len(p) bytes from p to the underlying data stream.
func (f *fileProxy) Write(b []byte) (n int, err error) {
	return f.file.Write(b)
}
