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
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	urlPtr := flag.String("url", "", "The url of the file")
	destPtr := flag.String("dest", "", "The destination of the downloaded file")

	flag.Parse()

	if urlPtr == nil || destPtr == nil {
		log.Fatalf("Expected url and destination to be filled.")
	}

	url := *urlPtr
	dest := *destPtr

	fmt.Printf("Downloading %q to %q\n", url, dest)

	_, err := os.Stat(dest)
	exists := !os.IsNotExist(err)
	if exists {
		fmt.Println("file already exists, skipping")
		return
	}

	// Setup local files first because it's cheaper to make these errors before
	// ones involving the network.
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		log.Fatalf("Error creating local directory %q: %v\n", filepath.Dir(dest), err)
	}
	out, err := os.Create(dest)
	if err != nil {
		log.Fatalf("Error opening local file %q: %v\n", dest, err)
	}
	defer out.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v\n", err)
	}
	req.Header.Set("User-Agent", "gcp-service-broker")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error getting HTTP resource: %v\n", err)
	}

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Got unexpected HTTP response code: %d\n", response.StatusCode)
	}

	if _, err := io.Copy(out, response.Body); err != nil {
		log.Fatalf("Error copying output: %v\n", err)
	}
}
