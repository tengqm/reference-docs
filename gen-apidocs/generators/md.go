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

package generators

import (
	"fmt"
	"io"
	// "io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

type MarkdownWriter struct {
	Config         *api.Config
	TOC            TOC
	CurrentSection *TOCItem
}

func NewMarkdownWriter(config *api.Config, copyright, title string) DocWriter {
	writer := MarkdownWriter{
		Config: config,
		TOC: TOC{
			Copyright: copyright,
			Title:     title,
			Sections:  []*TOCItem{},
		},
	}
	return &writer
}

func (m *MarkdownWriter) Extension() string {
	return ".md"
}

func (m *MarkdownWriter) DefaultStaticContent(title string) string {
	return fmt.Sprintf("## %s\n\n", title)
}

func (m *MarkdownWriter) WriteOverview() {
	fn := filepath.Join(api.SectionsDir, "overview.md")
	os.Link(fn, filepath.Join(api.BuildDir, "overview.md"))
	item := TOCItem{
		Level: 1,
		Title: "Overview",
		Link:  "overview",
		File:  fn,
	}
	m.TOC.Sections = append(m.TOC.Sections, &item)
	m.CurrentSection = &item
}

func (m *MarkdownWriter) WriteAPIGroupVersions(gvs api.GroupVersions) {
	fn := "group_versions.md"
	path := filepath.Join(api.BuildDir, fn)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	fmt.Fprintf(f, "---\ntitle: API Groups and Versions\nweight:20\n---\n\n")
	fmt.Fprintf(f, "The API Groups and their versions are summarized in the following table.\n\n")
	fmt.Fprintf(f, "<TABLE class=\"col-md-8\">\n<THEAD><TR><TH>Group</TH><TH>Version</TH></TR></THEAD>\n<TBODY>\n")

	groups := api.ApiGroups{}
	for group, _ := range gvs {
		groups = append(groups, api.ApiGroup(group))
	}
	sort.Sort(groups)

	for _, group := range groups {
		versionList, _ := gvs[group.String()]
		sort.Sort(versionList)
		var versions []string
		for _, v := range versionList {
			versions = append(versions, v.String())
		}

		fmt.Fprintf(f, "<TR><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD></TR>\n",
			group, strings.Join(versions, ", "))
	}
	fmt.Fprintf(f, "</TBODY>\n</TABLE>\n\n")

	item := TOCItem{
		Level: 1,
		Title: "API Groups",
		Link:  "api-groups",
		File:  fn,
	}
	m.TOC.Sections = append(m.TOC.Sections, &item)
	m.CurrentSection = &item
}

func (m *MarkdownWriter) WriteResourceCategory(name, file string) {
	fn := filepath.Join(api.SectionsDir, file+".md")
	os.Link(fn, filepath.Join(api.BuildDir, file+".md"))
	link := strings.Replace(strings.ToLower(name), " ", "-", -1)
	item := TOCItem{
		Level: 1,
		Title: strings.ToUpper(name),
		Link:  link,
		File:  fn,
	}
	m.TOC.Sections = append(m.TOC.Sections, &item)
	m.CurrentSection = &item
}

func (m *MarkdownWriter) writeFields(w io.Writer, d *api.Definition) {
	fmt.Fprintf(w, "Field        | Description\n------------ | -----------\n")
	for _, field := range d.Fields {
		fmt.Fprintf(w, "`%s`", field.Name)
		if field.Link() != "" {
			fmt.Fprintf(w, "<br /> *%s*", field.Link())
		}
		if field.PatchStrategy != "" {
			fmt.Fprintf(w, "<br /> **patch strategy**: *%s*", field.PatchStrategy)
		}
		if field.PatchMergeKey != "" {
			fmt.Fprintf(w, "<br /> **patch merge key**: *%s*", field.PatchMergeKey)
		}
		fmt.Fprintf(w, " | %s\n", field.DescriptionWithEntities)
	}
}

func (m *MarkdownWriter) writeOtherVersions(f io.Writer, d *api.Definition) {
	if d.OtherVersions.Len() != 0 {
		fmt.Fprintf(f, "### Other API versions:\n\n")
		for _, v := range d.OtherVersions {
			fmt.Fprintf(f, "- %s\n", v.MdLink())
		}
		fmt.Fprintf(f, "\n")
	}
}

func (m *MarkdownWriter) writeAppearsIn(f io.Writer, d *api.Definition) {
	if d.AppearsIn.Len() != 0 {
		fmt.Fprintf(f, "### Appears In:\n\n\n")
		for _, a := range d.AppearsIn {
			fmt.Fprintf(f, "- %s\n", a.MdLink())
		}
		fmt.Fprintf(f, "\n")
	}
}

func (m *MarkdownWriter) WriteDefinitionsOverview() {
	writeStaticFile("Definitions", "definitions.md", m.DefaultStaticContent("Definitions"))
	item := TOCItem{
		Level: 1,
		Title: "DEFINITIONS",
		Link:  "definitions",
		File:  "definitions.md",
	}
	m.TOC.Sections = append(m.TOC.Sections, &item)
	m.CurrentSection = &item
}

func (m *MarkdownWriter) WriteDefinition(d *api.Definition) {
	defname := strings.ToLower(strings.Replace(d.Name, ".", "-", 50))
	fn := fmt.Sprintf("%s-%s-%s.md", defname, d.Version, d.Group)
	dir := filepath.Join(api.BuildDir, "definitions")
	os.MkdirAll(dir, os.ModePerm)
	f, err := os.Create(filepath.Join(dir, fn))
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	nvg := fmt.Sprintf("%s (%s/%s)", d.Name, d.GroupDisplayName(), d.Version)
	linkID := getLink(nvg)
	fmt.Fprintf(f, "## %s {#%s}\n\n", nvg, linkID)
	fmt.Fprintf(f, "Group        | Version    | Kind\n------------ | ---------- | -----------\n")
	fmt.Fprintf(f, "`%s` | `%s` | `%s`\n", d.GroupDisplayName(), d.Version, d.Name)
	fmt.Fprintf(f, "\n")

	fmt.Fprintf(f, "\n%s\n\n", d.DescriptionWithEntities)
	m.writeFields(f, d)

	m.writeOtherVersions(f, d)
	m.writeAppearsIn(f, d)

	item := TOCItem{
		Level: 2,
		Title: nvg,
		Link:  linkID,
		File:  fn,
	}
	m.CurrentSection.SubSections = append(m.CurrentSection.SubSections, &item)
}

func (m *MarkdownWriter) WriteResource(r *api.Resource) {
	d := r.Definition
	defname := strings.ToLower(strings.Replace(d.Name, ".", "-", 50))
	fn := fmt.Sprintf("%s-%s-%s.md", defname, d.Version, d.Group)
	dir := filepath.Join(api.BuildDir, "resources")
	os.MkdirAll(dir, os.ModePerm)
	f, err := os.Create(filepath.Join(dir, fn))
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	dvg := fmt.Sprintf("%s (%s/%s)", r.Name, r.Definition.GroupDisplayName(), r.Definition.Version)
	linkID := getLink(dvg)
	fmt.Fprintf(f, "## %s {#%s}\n\n", dvg, linkID)

	if r.Definition.Sample.Sample != "" {
		note := r.Definition.Sample.Note
		for _, t := range r.Definition.GetSamples() {
			// 't' is a ExampleText
			fmt.Fprintf(f, ">%s %s\n\n", t.Tab, note)
			fmt.Fprintf(f, "```%s\n%s\n```\n\n", t.Type, t.Text)
		}
	}

	// GVK
	fmt.Fprintf(f, "Group        | Version    | Kind\n------------ | ---------- | -----------\n")
	fmt.Fprintf(f, "`%s` | `%s` | `%s`\n\n", r.Definition.GroupDisplayName(), r.Definition.Version, r.Name)

	// TODO: Use shortcode
	if r.DescriptionWarning != "" {
		fmt.Fprintf(f, "<aside class=\"warning\">%s</aside>\n\n", r.DescriptionWarning)
	}
	if r.DescriptionNote != "" {
		fmt.Fprintf(f, "<aside class=\"notice\">%s</aside>\n\n", r.DescriptionNote)
	}

	m.writeOtherVersions(f, r.Definition)
	m.writeAppearsIn(f, r.Definition)
	m.writeFields(f, r.Definition)

	fmt.Fprintf(f, "\n")
	if r.Definition.Inline.Len() > 0 {
		for _, d := range r.Definition.Inline {
			fmt.Fprintf(f, "### %s (%s/%s)\n\n", d.Name, d.Group, d.Version)
			m.writeAppearsIn(f, d)
			m.writeFields(f, d)
		}
	}

	item := TOCItem{
		Level: 1,
		Title: dvg,
		Link:  linkID,
		File:  fn,
	}
	m.TOC.Sections = append(m.TOC.Sections, &item)
	m.CurrentSection = &item

	if len(r.Definition.OperationCategories) == 0 {
		return
	}

	for _, c := range r.Definition.OperationCategories {
		if len(c.Operations) == 0 {
			continue
		}

		catID := strings.Replace(strings.ToLower(c.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
		fmt.Fprintf(f, "## %s {#%s}\n\n", c.Name, catID)
		OCItem := TOCItem{
			Level: 2,
			Title: c.Name,
			Link:  catID,
		}
		m.CurrentSection.SubSections = append(m.CurrentSection.SubSections, &OCItem)

		for _, o := range c.Operations {
			opID := strings.Replace(strings.ToLower(o.Type.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
			fmt.Fprintf(f, "\n### %s {#%s}\n\n", o.Type.Name, opID)

			OPItem := TOCItem{
				Level: 2,
				Title: o.Type.Name,
				Link:  opID,
			}
			OCItem.SubSections = append(OCItem.SubSections, &OPItem)

			// Example requests
			requests := o.GetExampleRequests()
			if len(requests) > 0 {
				m.writeOperationSample(f, true, opID, requests)
			}

			// Example responses
			responses := o.GetExampleResponses()
			if len(responses) > 0 {
				m.writeOperationSample(f, false, opID, responses)
			}

			fmt.Fprintf(f, "%s\n", o.Description())
			fmt.Fprintf(f, "\n#### HTTP Request\n\n`%s`\n\n", o.GetDisplayHttp())

			m.writeRequestParams(f, o)
			m.writeResponseParams(f, o)
		}
	}
}

func (m *MarkdownWriter) writeOperationSample(w io.Writer, req bool, op string, examples []api.ExampleText) {
	for _, e := range examples {
		eType := strings.Split(e.Tab, ":")[1]
		var btnText string
		if req {
			btnText = eType + " request"
		} else {
			btnText = eType + " response"
		}
		fmt.Fprintf(w, "**%s example**\n\n", btnText)

		msg := e.Msg
		if eType == "curl" && strings.Contains(msg, "proxy") {
			msg = "`curl` command (*requires `kubectl proxy` to be running*)"
		} else if eType == "kubectl" && strings.Contains(msg, "Command") { // `kubectl` command
			msg = "`kubectl` command"
		}
		lType := strings.Split(e.Type, ":")[1]
		lang := strings.Split(lType, "_")[1]
		fmt.Fprintf(w, "```%s\n%s\n```\n\n", lang, e.Text)
	}
}

func (m *MarkdownWriter) writeParams(f io.Writer, title string, params api.Fields) {
	fmt.Fprintf(f, "##### %s\n\n", title)
	fmt.Fprintf(f, "<TABLE>\n<THEAD><TR><TH>Parameter</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	for _, p := range params {
		fmt.Fprintf(f, "<TR><TD><CODE>%s</CODE>", p.Name)
		if p.Link() != "" {
			fmt.Fprintf(f, "<br /><I>%s</I>", p.FullLink())
		}
		fmt.Fprintf(f, "</TD><TD>%s</TD></TR>\n", p.Description)
	}
	fmt.Fprintf(f, "</TBODY>\n</TABLE>\n\n")
}

func (m *MarkdownWriter) writeRequestParams(w io.Writer, o *api.Operation) {
	// Operation path params
	if o.PathParams.Len() > 0 {
		m.writeParams(w, "Path Parameters", o.PathParams)
	}

	// operation query params
	if o.QueryParams.Len() > 0 {
		m.writeParams(w, "Query Parameters", o.QueryParams)
	}

	// operation body params
	if o.BodyParams.Len() > 0 {
		m.writeParams(w, "Body Parameters", o.BodyParams)
	}
}

func (m *MarkdownWriter) writeResponseParams(f io.Writer, o *api.Operation) {
	if o.HttpResponses.Len() == 0 {
		return
	}

	fmt.Fprintf(f, "#### Response\n\n")
	fmt.Fprintf(f, "<TABLE>\n<THEAD><TR><TH>Code</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	responses := o.HttpResponses
	sort.Slice(responses, func(i, j int) bool {
		return strings.Compare(responses[i].Name, responses[j].Name) < 0
	})
	for _, p := range responses {
		fmt.Fprintf(f, "<TR><TD>%s", p.Name)
		if p.Field.Link() != "" {
			fmt.Fprintf(f, "<br /><I>%s</I>", p.Field.FullLink())
		}
		fmt.Fprintf(f, "</TD><TD>%s</TD></TR>\n", p.Field.Description)
	}
	fmt.Fprintf(f, "</TBODY>\n</TABLE>\n\n")
}

func (m *MarkdownWriter) WriteOldVersionsOverview() {
	writeStaticFile("Old Versions", "oldversions.md", m.DefaultStaticContent("Old Versions"))

	item := TOCItem{
		Level: 1,
		Title: "Old API Versions",
		Link:  "old-versions",
		File:  "oldversions.md",
	}
	m.TOC.Sections = append(m.TOC.Sections, &item)
	m.CurrentSection = &item
}

func (m *MarkdownWriter) generateNavContent() string {
	/*
		nav := ""
		for _, sec := range h.TOC.Sections {
			// class for level-1 navigation item
			nav += "<UL>\n"
			if strings.Contains(sec.Link, "strong") {
				nav += fmt.Sprintf(" <LI class=\"nav-level-1 strong-nav\"><A href=\"#%s\" class=\"nav-item\"><STRONG>%s</STRONG></A></LI>\n", sec.Link, sec.Title)
			} else {
				nav += fmt.Sprintf(" <LI class=\"nav-level-1\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n",
					sec.Link, sec.Title)
			}

			// close H1 items which have no subsections or strong navs
			if len(sec.SubSections) == 0 || (sec.Level == 1 && strings.Contains(sec.Link, "strong")) {
				nav += "</UL>\n"
			}

			// short circuit to next if no sub-sections
			if len(sec.SubSections) == 0 {
				continue
			}

			// wrapper1
			nav += fmt.Sprintf(" <UL id=\"%s-nav\" style=\"display: none;\">\n", sec.Link)
			for _, sub := range sec.SubSections {
				nav += "  <UL>\n"
				if strings.Contains(sub.Link, "strong") {
					nav += fmt.Sprintf("   <LI class=\"nav-level-%d strong-nav\"><A href=\"#%s\" class=\"nav-item\"><STRONG>%s</STRONG></A></LI>\n",
						sub.Level, sub.Link, sub.Title)
				} else {
					nav += fmt.Sprintf("   <LI class=\"nav-level-%d\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n",
						sub.Level, sub.Link, sub.Title)
				}
				// close this H1/H2 if possible
				if len(sub.SubSections) == 0 {
					nav += " </UL>\n"
					continue
				}

				// 3rd level
				// another wrapper
				nav += fmt.Sprintf("   <UL id=\"%s-nav\" style=\"display: none;\">\n", sub.Link)
				for _, subsub := range sub.SubSections {
					nav += fmt.Sprintf("    <LI class=\"nav-level-%d\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n", subsub.Level, subsub.Link, subsub.Title)
					if len(subsub.SubSections) == 0 {
						continue
					}

					fmt.Printf("*** found third level!\n")
					nav += fmt.Sprintf("   <UL id=\"%s-nav\" style=\"display: none;\">\n", subsub.Link)
					for _, subsubsub := range subsub.SubSections {
						nav += fmt.Sprintf("    <LI class=\"nav-level-%d\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n",
							subsubsub.Level, subsubsub.Link, subsubsub.Title)
					}
					nav += "   </UL>\n"
				}
				// end wrapper2
				nav += "   </UL>\n"
				nav += "  </UL>\n"
			}
			// end wrapper1
			nav += " </UL>\n"
			// end top UL
			nav += "</UL>\n"
		}
	*/
	return ""
}

func (m *MarkdownWriter) Finalize() {
	os.MkdirAll(api.BuildDir, os.ModePerm)

	buf := ""
	data, err := ioutil.ReadFile(filepath.Join(api.SectionsDir, "index.md"))
	if err == nil {
		buf += string(data)
		fmt.Println(OK)
	} else {
		fmt.Println(NOT_FOUND)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05 (MST)")
	buf = strings.Replace(buf, "TIMESTAMP", timestamp, -1)

	pos := strings.LastIndex(m.Config.SpecVersion, ".")
	release := m.Config.SpecVersion[1:pos]
	buf = strings.Replace(buf, "RELEASE", release, -1)

	ioutil.WriteFile(filepath.Join(api.BuildDir, "_index.md"), buf, 0)

	const OK = "\033[32mOK\033[0m"
	const NOT_FOUND = "\033[31mNot found\033[0m"

	for _, sec := range m.TOC.Sections {
		fmt.Printf("Processing %s ... ", sec.File)

		if err == nil {
			fmt.Println(OK)
		} else {
			fmt.Println(NOT_FOUND)
			continue
		}

		for _, sub := range sec.SubSections {
			if len(sub.File) > 0 {
				fmt.Printf("Writing %s ... ", sub.File)
				if err == nil {
					fmt.Println(OK)
				} else {
					fmt.Println(NOT_FOUND)
					continue
				}
			}

			for _, subsub := range sub.SubSections {
				if len(subsub.File) > 0 {
					fmt.Printf("Processing %s ...", subsub.File)
					if err == nil {
						fmt.Println(OK)
					} else {
						fmt.Println(NOT_FOUND)
					}
				}
			}
		}
	}
}
