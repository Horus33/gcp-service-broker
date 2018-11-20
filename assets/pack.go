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
	"log"
	"os"
	"path/filepath"

	resources "github.com/omeid/go-resources"
)

func main() {
	pathPtr := flag.String("path", "", "The path of the directory to pack up")
	varPtr := flag.String("var", "", "The name of the variable for this package")

	flag.Parse()

	if pathPtr == nil || varPtr == nil {
		log.Fatalf("Expected path and var to be filled.")
	}

	compilePackage(*varPtr, "assets", *pathPtr)
}

func compilePackage(varname, pkgName, dir string) {
	pkg := resources.New()
	pkg.Config.Pkg = pkgName
	pkg.Config.Var = varname
	pkg.Config.Format = false
	pkg.Config.Declare = true

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if err := pkg.AddFile(path[len(dir):], path); err != nil {
			return err
		}

		fmt.Printf("added file: %q\n", path)
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println("writing output")

	pkg.Write(varname + ".go")
}
