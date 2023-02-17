/*
Copyright 2018 The Kubernetes Authors.

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
package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

type Doc struct {
	Filename string `json:"filename,omitempty"`
}

type DocWriter interface {
	Extension() string
	DefaultStaticContent(title string) string
	WriteOverview()
	WriteAPIGroupVersions(gvs api.GroupVersions)
	WriteResourceCategory(name, file string)
	WriteResource(r *api.Resource)
	WriteDefinitionsOverview()
	WriteDefinition(d *api.Definition)
	WriteOldVersionsOverview()
	Finalize()
}

func GenerateFiles(format string) {
	// Load the yaml config
	config := api.NewConfig()
	PrintInfo(config)
	ensureIncludeDir()

	var copyright_tmpl string
	if format == "html" {
		copyright_tmpl = "<a href=\"https://github.com/kubernetes/kubernetes\">Copyright 2016-%s The Kubernetes Authors.</a>"
	} else {
		copyright_tmpl = "[Copyright 2016-%s The Kubernetes Authors](https://github.com/kubernetes/kubernetes)"
	}

	now := time.Now().Format("2006")
	copyright := fmt.Sprintf(copyright_tmpl, now)
	var title string
	if !*api.BuildOps {
		title = "Kubernetes Resource Reference Docs"
	} else {
		title = "Kubernetes API Reference Docs"
	}
	var writer DocWriter
	if format == "html" {
		writer = NewHTMLWriter(config, copyright, title)
	} else {
		writer = NewMarkdownWriter(config, copyright, title)
	}
	writer.WriteOverview()

	// Write API groups
	writer.WriteAPIGroupVersions(config.Definitions.GroupVersions)

	// Write resource definitions
	for _, c := range config.ResourceCategories {
		writer.WriteResourceCategory(c.Name, c.Include)
		for _, r := range c.Resources {
			if r.Definition == nil {
				fmt.Printf("Warning: Missing definition for item in TOC %s\n", r.Name)
				continue
			}
			writer.WriteResource(r)
		}
	}

	writer.WriteDefinitionsOverview()
	// Add other definition imports
	definitions := api.SortDefinitionsByName{}
	for _, d := range config.Definitions.All {
		// Don't add definitions for top level resources in the toc or inlined resources
		if d.InToc || d.IsInlined || d.IsOldVersion {
			continue
		}
		definitions = append(definitions, d)
	}
	sort.Sort(definitions)
	for _, d := range definitions {
		writer.WriteDefinition(d)
	}

	writer.WriteOldVersionsOverview()
	oldversions := api.SortDefinitionsByName{}
	for _, d := range config.Definitions.All {
		// Don't add definitions for top level resources in the toc or inlined resources
		if d.IsOldVersion {
			oldversions = append(oldversions, d)
		}
	}
	sort.Sort(oldversions)
	for _, d := range oldversions {
		// Skip Inlined definitions
		if d.IsInlined {
			continue
		}
		r := &api.Resource{Definition: d, Name: d.Name}
		writer.WriteResource(r)
	}

	writer.Finalize()
}

func ensureIncludeDir() {
	if _, err := os.Stat(api.BuildDir); os.IsNotExist(err) {
		os.Mkdir(api.BuildDir, os.FileMode(0700))
	}
	if _, err := os.Stat(api.IncludesDir); os.IsNotExist(err) {
		os.Mkdir(api.IncludesDir, os.FileMode(0700))
	}
}

func definitionFileName(d *api.Definition) string {
	name := "generated_" + strings.ToLower(strings.Replace(d.Name, ".", "_", 50))
	return fmt.Sprintf("%s_%s_%s_definition", name, d.Version, d.Group)
}

func conceptFileName(d *api.Definition) string {
	name := "generated_" + strings.ToLower(strings.Replace(d.Name, ".", "_", 50))
	return fmt.Sprintf("%s_%s_%s_concept", name, d.Version, d.Group)
}

func getLink(s string) string {
	tmp := strings.Replace(s, ".", "-", -1)
	return strings.ToLower(strings.Replace(tmp, " ", "-", -1))
}

func writeStaticFile(title, location, defaultContent string) {
	fn := filepath.Join(api.SectionsDir, location)
	to := filepath.Join(api.IncludesDir, location)
	_, err := os.Stat(fn)
	if err == nil {
		// copy the file if it exists
		os.Link(fn, to)
		return
	}

	if !os.IsNotExist(err) {
		panic(fmt.Sprintf("Could not stat file %s %v", fn, err))
	}
	fmt.Printf("Creating file %s\n", to)
	file, err := os.Create(to)
	if err != nil {
		panic(err)
	}
	file.Close()

	file, err = os.OpenFile(to, os.O_WRONLY, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(file, defaultContent)
	file.Close()
}
