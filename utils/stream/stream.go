// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stream

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Source func() (io.ReadCloser, error)
type Dest func() (io.WriteCloser, error)

// Copy copies data from a source stream to a destination stream.
func Copy(src Source, dest Dest) error {
	readCloser, err := src()
	if err != nil {
		return fmt.Errorf("copy couldn't open source: %v", err)
	}
	defer readCloser.Close()

	writeCloser, err := dest()
	if err != nil {
		return fmt.Errorf("copy couldn't open destination: %v", err)
	}
	defer writeCloser.Close()

	if _, err := io.Copy(writeCloser, readCloser); err != nil {
		return fmt.Errorf("copy couldn't copy data: %v", err)
	}

	if err := readCloser.Close(); err != nil {
		return fmt.Errorf("copy couldn't close source: %v", err)
	}

	if err := writeCloser.Close(); err != nil {
		return fmt.Errorf("copy couldn't close destination: %v", err)
	}

	return nil
}

// FromYaml converts the interface to a stream of Yaml.
func FromYaml(v interface{}) Source {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return FromError(err)
	}

	return FromBytes(bytes)
}

// FromBytes streams the given bytes as a buffer.
func FromBytes(b []byte) Source {
	return FromReader(bytes.NewReader(b))
}

// FromBytes streams the given bytes as a buffer.
func FromString(s string) Source {
	return FromBytes([]byte(s))
}

// FromError returns a nil ReadCloser and the error passed when called.
func FromError(err error) Source {
	return func() (io.ReadCloser, error) {
		return nil, err
	}
}

// FromFile joins the segments of the path and reads from it.
func FromFile(path ...string) Source {
	return FromReadCloserError(os.Open(filepath.Join(path...)))
}

// FromReadCloserError reads the contents of the readcloser and takes ownership of closing it.
// If err is non-nil, it is returned as a source.
func FromReadCloserError(rc io.ReadCloser, err error) Source {
	return func() (io.ReadCloser, error) {
		return rc, err
	}
}

// FromReader converts a Reader to a ReadCloser and takes ownership of closing it.
func FromReader(rc io.Reader) Source {
	return func() (io.ReadCloser, error) {
		return ioutil.NopCloser(rc), nil
	}
}

// ToFile concatenates the given path segments with filepath.Join, creates any
// parent directoreis if needed, and writes the file.
func ToFile(path ...string) Dest {
	return ToModeFile(0600, path...)
}

// ToModeFile is like ToFile, but sets the permissions on the created file.
func ToModeFile(mode os.FileMode, path ...string) Dest {
	return func() (io.WriteCloser, error) {
		outputPath := filepath.Join(path...)
		if err := os.MkdirAll(filepath.Dir(outputPath), 0700|os.ModeDir); err != nil {
			return nil, err
		}

		return os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	}
}

// ToBuffer buffers the contents of the stream and on Close() calls the callback
// returning its error.
func ToBuffer(closeCallback func(*bytes.Buffer) error) Dest {
	return func() (io.WriteCloser, error) {
		return &bufferCallback{closeCallback: closeCallback}, nil
	}
}

// ToYaml unmarshals the contents of the stream as YAML to the given struct.
func ToYaml(v interface{}) Dest {
	return ToBuffer(func(buf *bytes.Buffer) error {
		return yaml.Unmarshal(buf.Bytes(), v)
	})
}

// ToError returns an error when Dest is initialized.
func ToError(err error) Dest {
	return func() (io.WriteCloser, error) {
		return nil, err
	}
}

// ToDiscard discards any data written to it.
func ToDiscard() Dest {
	return ToWriter(ioutil.Discard)
}

// ToWriter forwards data to the given writer, this function WILL NOT close the
// underlying stream so it is safe to use with things like stdout.
func ToWriter(writer io.Writer) Dest {
	return ToWriteCloser(NopWriteCloser(writer))
}

// ToWriteCloser forwards data to the given WriteCloser which will be closed
// after the copy finishes.
func ToWriteCloser(w io.WriteCloser) Dest {
	return func() (io.WriteCloser, error) {
		return w, nil
	}
}

// bufferCallback buffers the results and on close calls the callback.
type bufferCallback struct {
	bytes.Buffer
	closeCallback func(*bytes.Buffer) error
}

// Close implements io.Closer
func (b *bufferCallback) Close() error {
	return b.closeCallback(&b.Buffer)
}

// NopWriteCloser works like io.NopCloser, but for writers.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return errWriteCloser{Writer: w, CloseErr: nil}
}

type errWriteCloser struct {
	io.Writer
	CloseErr error
}

func (w errWriteCloser) Close() error {
	return w.CloseErr
}