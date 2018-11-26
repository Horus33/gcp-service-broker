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

// +build ignore

package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s packagedest packagename dirtozip", os.Args[0])
	}

	packageDest := os.Args[1]
	packageName := os.Args[2]
	dirToZip := os.Args[3]

	// create package
	pkgRoot := filepath.Join(packageDest, packageName)
	os.RemoveAll(pkgRoot)
	if err := os.MkdirAll(pkgRoot, os.ModeDir|0777); err != nil {
		log.Fatalf("couldn't create package %q, %v", pkgRoot, err)
	}

	// zip directory up in it
	// zipFile := filepath.Join(pkgRoot, "packed.zip")
	// defer os.Remove(zipFile)
	buf, err := zipDirectory(dirToZip, pkgRoot)
	if err != nil {
		log.Fatalf("couldn't create zip from %q %v", dirToZip, err)
	}

	chunkCount, err := writeChunks(pkgRoot, packageName, buf)
	if err != nil {
		log.Fatalf("couldn't create %d chunks %v", chunkCount, err)
	}

	// create index file

	// create chunks

}

func zipDirectory(dir, to string) (*bytes.Buffer, error) {
	f := new(bytes.Buffer)
	// f, err := os.Create(to)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		zipPath := strings.Join(filepath.SplitList(path[len(dir)+1:]), "/")
		log.Printf("Packing %q as %q\n", path, zipPath)
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		return readFile(path, writer)
	})

	if err != nil {
		return f, err
	}

	return f, nil
}

func readFile(src string, dest io.Writer) error {
	fd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(dest, fd)
	return err
}

// writeChunks divides the buffer up into portions and writes them under
// the given directory. It returns the number of chunks written as well as an
// error if it was encountered.
func writeChunks(dir, packageName string, buf *bytes.Buffer) (int, error) {
	count := 0
	for {
		chunk := buf.Next(1024 * 1024)
		if len(chunk) == 0 {
			return count, nil
		}
		chunkid := filepath.Join(dir, fmt.Sprintf("chunk%d.go", count))
		contents := fmt.Sprintf("package %s\n\nvar embedded%d = %#v\n", packageName, count, chunk)
		err := ioutil.WriteFile(chunkid, []byte(contents), 0644)
		if err != nil {
			return count, err
		}
		count += 1
	}
}

var indexTemplate = `
package chunked

import (
	"archive/zip"
	"bytes"
	"io"
)

func NewZipReader() (*zip.Reader, error) {
	fd := File{
		embedded0,
	}

	return zip.NewReader(fd, fd.Length())
}

type chunk []byte
type File []chunk

func (f File) ReadAt(p []byte, off int64) (n int, err error) {
	iOff := int(off)
	slice := f.slice(iOff, iOff+len(p))

	for i, b := range slice {
		p[i] = b
	}

	n = len(slice)
	if n < len(p) {
		err = io.EOF
	}

	return
}

func (f File) Length() int64 {
	l := 0
	for _, chunk := range f {
		l += len(chunk)
	}
	return int64(l)
}

func (f File) slice(start, end int) []byte {
	buf := new(bytes.Buffer)

	chunkStart := 0
	for _, chunk := range f {
		chunkEnd := chunkStart + len(chunk)

		overlapStart := imin(chunkEnd, imax(chunkStart, start))
		overlapEnd := imax(chunkStart, imin(chunkEnd, end))

		buf.Write(chunk[overlapStart:overlapEnd])
		chunkStart = chunkEnd
	}

	return buf.Bytes()
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
``
