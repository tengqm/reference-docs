/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators"
	"github.com/pkg/errors"
)

var flFormat = flag.String("f", "markdown", "format for output, one of 'html' and 'markdown'.")

func main() {
	flag.Parse()
	if *flFormat == "html" || *flFormat == "markdown" {
		generators.GenerateFiles(*flFormat)
	} else {
		panic(errors.Errorf("unsupported format '%s' specified", *flFormat))
	}
}
