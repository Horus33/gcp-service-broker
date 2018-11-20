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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	srcPtr := flag.String("src", "", "The source zip file")
	destPtr := flag.String("dest", "", "The destination directory")

	flag.Parse()

	if srcPtr == nil || destPtr == nil {
		log.Fatalf("Expected src and dest to be non-nil")
	}

	src := *srcPtr
	dest := *destPtr
	fmt.Printf("Unzipping %q to %q\n", src, dest)

	zipReader, err := zip.OpenReader(src)
	if err != nil {
		log.Fatalf("Error while unzipping: %v\n", err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		path := filepath.Clean(filepath.Join(dest, f.Name))
		if !strings.HasPrefix(path, dest) {
			log.Fatalf("Error while unzipping: %v\n", fmt.Errorf("Possible vulnerability extracting %q", f.Name))
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				log.Fatalf("Error while unzipping: %v\n", err)
			}
		} else {
			contents, err := f.Open()
			if err != nil {
				log.Fatalf("Error while unzipping: %v\n", err)
			}
			if err := writeFile(contents, path, f.Mode()); err != nil {
				log.Fatalf("Error while unzipping: %v\n", err)
			}
		}
	}

}

func writeFile(src io.ReadCloser, dest string, perm os.FileMode) error {
	defer src.Close()
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
