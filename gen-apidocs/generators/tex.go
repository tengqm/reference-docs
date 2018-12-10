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
	"html"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubernetes-incubator/reference-docs/gen-apidocs/generators/api"
)

type tocItem struct {
	Level int
	Title string
	Link string
	File string
	SubSections []*tocItem
}

type toc struct {
	Title string
	Copyright string
	Sections []*tocItem
}
type TexWriter struct {
	Config *api.Config
	toc toc
	CurrentSection *tocItem
}

func NewTexWriter(config *api.Config, copyright, title string) DocWriter {
	writer := TexWriter{
		Config: config,
		toc: toc {
			Copyright: copyright,
			Title: title,
			Sections: []*tocItem{},
		},
	}
	return &writer
}

func (h *TexWriter) Extension() string {
	return ".tex"
}

func (h *TexWriter) WriteOverview() {
	fn := "_overview.tex"
	writeStaticFile("Overview", fn, h.DefaultStaticContent("Overview"))
	item := tocItem {
		Level: 1,
		Title: "Overview",
		Link: "-strong-api-overview-strong-",
		File: fn,
	}
	h.toc.Sections = append(h.toc.Sections, &item)
	h.CurrentSection = &item
}

func (h *TexWriter) WriteResourceCategory(name, file string) {
	writeStaticFile(name, file + ".tex", h.DefaultStaticContent(name))
	link := strings.Replace(strings.ToLower(name), " ", "-", -1)
	item := tocItem {
		Level: 1,
		Title: strings.ToUpper(name),
		Link: "-strong-" + link + "-strong-",
		File: file + ".tex",
	}
	h.toc.Sections = append(h.toc.Sections, &item)
	h.CurrentSection = &item
}

func (h *TexWriter) DefaultStaticContent(title string) string {
	titleID := strings.ToLower(strings.Replace(title, " ", "-", -1))
	return fmt.Sprintf("\\section{%s}\n\\label{-strong-%s-strong-}\n", titleID, title)
}

func (h *TexWriter) writeOtherVersions(w io.Writer, d *api.Definition) {
	if d.OtherVersions.Len() == 0 {
		return
	}

	// TODO(Qiming): Use TColorbox
	// fmt.Fprint(w, "<DIV class=\"alert alert-success col-md-8\"><I class=\"fa fa-toggle-right\"></I> Other API versions of this object exist:\n")
	fmt.Fprint(w, "Other API versions of this object exist:\n\\begin{itemize}\n")
	for _, v := range d.OtherVersions {
		link, text := v.VersionLinkData()
		fmt.Fprintf(w, "\\item \\href{%s}{%s}\n", link, text)
	}
	fmt.Fprintf(w, "\\end{itemize}")
}

func (h *TexWriter) writeAppearsIn(w io.Writer, d *api.Definition) {
	if d.AppearsIn.Len() != 0 {
		// TODO(Qiming): Use TColorbox
		// fmt.Fprintf(w, "<DIV class=\"alert alert-info col-md-8\"><I class=\"fa fa-info-circle\"></I> Appears In:\n <UL>\n")
		fmt.Fprintf(w, "Appears In:\n\n \\begin{itemize}\n")
		for _, a := range d.AppearsIn {
			link, text := a.FullHrefLinkData()
			fmt.Fprintf(w, "  \\item \\href{%s}{%s}\n", link, text)
		}
		fmt.Fprintf(w, " \\end{itemize}\n\n")
	}
}

func (h *TexWriter) writeFields(w io.Writer, d *api.Definition) {
	// fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Field</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	fmt.Fprintf(w, "\\begin{center}\n\\begin{table}{l|l}\n")
	fmt.Fprintf(w, "\\newcommand{\\tabincell}[2]{\\begin{tabular}{@{}#1@{}}#2\\end{tabular}}\n")
	fmt.Fprintf(w, "\\toprule\nField & Description\\\\\n")
	fmt.Fprintf(w, "\\midrule\n\\endfirsthead\n")
	fmt.Fprintf(w, "\\multicolumn{2}{r}{<Continue>}\\\\\n")
	fmt.Fprintf(w, "\\toprule\nField & Description\\\\\n")
	fmt.Fprintf(w, "\\endhead\n\\bottomrule\n")
	fmt.Fprintf(w, "\\multicolumn{2}{c}{<Continue on next page>}\n")
	fmt.Fprintf(w, "\\endfoot\n\\bottomrule\n\\endlastfoot\n")

	for _, field := range d.Fields {
		fmt.Fprintf(w, "\\tabincell{l}{\\texttt{%s}", field.Name)
		if field.Link() != "" {
			fmt.Fprintf(w, "\\\\ \\textit{%s}", field.FullLink())
		}
		if field.PatchStrategy != "" {
			fmt.Fprintf(w, "\\\\ \\textbf{patch strategy}: \\textit{%s}", field.PatchStrategy)
		}
		if field.PatchMergeKey != "" {
			fmt.Fprintf(w, "\\\\ \\textbf{patch merge key}: \\textit{%s}", field.PatchMergeKey)
		}
		fmt.Fprintf(w, "} ")
		// TODO(Qiming): Fix the HTML escape below
		fmt.Fprintf(w, "& %s \\\\\n", field.DescriptionWithEntities)
	}
	fmt.Fprintf(w, "\\end{longtable}\n\\end{center}\n\n")
}

func (h *TexWriter) WriteDefinitionsOverview() {
	writeStaticFile("Definitions", "_definitions.tex", h.DefaultStaticContent("Definitions"))
	item := tocItem {
		Level: 1,
		Title: "DEFINITIONS",
		Link: "-strong-definitions-strong-",
		File: "_definitions.tex",
	}
	h.toc.Sections = append(h.toc.Sections, &item)
	h.CurrentSection = &item
}

func (h *TexWriter) WriteDefinition(d *api.Definition) {
	fn := "_" + definitionFileName(d) + ".tex"
	path := *api.ConfigDir + "/includes/" + fn
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	nvg := fmt.Sprintf("%s %s %s", d.Name, d.Version, d.GroupDisplayName())
	linkID := getLink(nvg)
	fmt.Fprintf(f, "\\subsection{%s}\n\\label{%s}\n\n", nvg, linkID)
	fmt.Fprintf(f, "\\begin{table}[htbp]\n\\begin{center}\n")
	fmt.Fprintf(f, "\\begin{tabular}{c|c|c}\n\\hline")
	fmt.Fprintf(f, "Group & Version & Kind\\\\\n\\hline")
	fmt.Fprintf(f, "\\texttt{%s} & \\texttt{%s} & \\texttt{%s}\\\\\n",
	            d.GroupDisplayName(), d.Version, d.Name)
	fmt.Fprintf(f, "\\hline\n\\end{tabular}\\end{center}\\end{table}\n\n")

	fmt.Fprintf(f, "%s\n", d.DescriptionWithEntities)
	h.writeOtherVersions(f, d)
	h.writeAppearsIn(f, d)
	h.writeFields(f, d)

	item := tocItem {
		Level: 2,
		Title: nvg,
		Link: linkID,
		File: fn,
	}
	h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &item)
}

func (h *TexWriter) writeSample(w io.Writer, d *api.Definition) {
	if d.Sample.Sample == "" {
		return
	}

	note := d.Sample.Note
	for _, s := range d.GetSamples() {
		sType := strings.Split(s.Tab, ":")[1]
		linkID := sType + "-" + d.LinkID()
		fmt.Fprintf(w, "<BUTTON class=\"btn btn-info\" type=\"button\" data-toggle=\"collapse\"\n")
		fmt.Fprintf(w, "  data-target=\"#%s\" aria-controls=\"%s\"\n", linkID, linkID)
		fmt.Fprintf(w, "  aria-expanded=\"false\">%s example</BUTTON>\n", sType)
	}

	for _, s := range d.GetSamples() {
		sType := strings.Split(s.Tab, ":")[1]
		linkID := sType + "-" + d.LinkID()
		lType := strings.Split(s.Type, ":")[1]
		lang := strings.Split(lType, "_")[1]
		fmt.Fprintf(w, "<DIV class=\"collapse\" id=\"%s\">\n", linkID)
		fmt.Fprintf(w, "  <DIV class=\"panel panel-default\">\n<DIV class=\"panel-heading\">%s</DIV>\n", note)
		fmt.Fprintf(w, "  <DIV class=\"panel-body\">\n<PRE class=\"%s\">", sType)
		fmt.Fprintf(w, "<CODE class=\"lang-%s\">\n", lang)
		// TODO: Add language highlight
		fmt.Fprintf(w, "%s\n</CODE></PRE></DIV></DIV></DIV>\n", html.EscapeString(s.Text))
	}
}

func (h *TexWriter) writeOperationSample(w io.Writer, req bool, op string, examples []api.ExampleText) {
	// e.Tab bdocs-tab:kubectl  | bdocs-tab:curl
	// e.Msg `kubectl` Command  | Output | Response Body | `curl` Command (*requires `kubectl proxy` to be running*)
	// e.Type bdocs-tab:kubectl_shell
	// e.Text <actual command>

	for _, e := range examples {
		eType := strings.Split(e.Tab, ":")[1]
		var sampleID string
		var btnText string
		if req {
			sampleID = "req-" + eType + "-" + op
			btnText = eType + " request"
		} else {
			sampleID = "res-" + eType + "-" + op
			btnText = eType + " response"
		}
		fmt.Fprintf(w, "<BUTTON class=\"btn btn-info\" type=\"button\" data-toggle=\"collapse\"\n")
		fmt.Fprintf(w, "  data-target=\"#%s\" aria-controls=\"%s\"\n", sampleID, sampleID)
		fmt.Fprintf(w, "  aria-expanded=\"false\">%s example</BUTTON>\n", btnText)
	}

	for _, e := range examples {
		eType := strings.Split(e.Tab, ":")[1]
		var sampleID string
		if req {
			sampleID = "req-" + eType + "-" + op
		} else {
			sampleID = "res-" + eType + "-" + op
		}
		msg := e.Msg
		if eType == "curl" && strings.Contains(msg, "proxy") {
			msg = "<CODE>curl</CODE> command (<I>requires <code>kubectl proxy</code> to be running</I>)"
		} else if eType == "kubectl" && strings.Contains(msg, "Command") { // `kubectl` command
			msg = "<CODE>kubectl</CODE> command"
		}
		lType := strings.Split(e.Type, ":")[1]
		lang := strings.Split(lType, "_")[1]
		fmt.Fprintf(w, "<DIV class=\"collapse\" id=\"%s\">\n", sampleID)
		fmt.Fprintf(w, "  <DIV class=\"panel panel-default\">\n<DIV class=\"panel-heading\">%s</DIV>\n", msg)
		fmt.Fprintf(w, "  <DIV class=\"panel-body\">\n<PRE class=\"%s\">", eType)
		fmt.Fprintf(w, "<CODE class=\"lang-%s\">\n", lang)
		// TODO: Add language highlight
		fmt.Fprintf(w, "%s\n</CODE></PRE></DIV></DIV></DIV>\n", e.Text)
	}
}

func (h *TexWriter) writeParams(w io.Writer, title string, params api.Fields) {
	fmt.Fprintf(w, "<H3>%s</H3>\n", title)
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Parameter</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	for _, p := range params {
		fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE>", p.Name)
		if p.Link() != "" {
			fmt.Fprintf(w, "<br /><I>%s</I>", p.FullLink())
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", p.Description)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}

func (h *TexWriter) writeRequestParams(w io.Writer, o *api.Operation) {
	// Operation path params
	if o.PathParams.Len() > 0 {
		h.writeParams(w, "Path Parameters", o.PathParams)
	}

	// operation query params
	if o.QueryParams.Len() > 0 {
		h.writeParams(w, "Query Parameters", o.QueryParams)
	}

	// operation body params
	if o.BodyParams.Len() > 0 {
		h.writeParams(w, "Body Parameters", o.BodyParams)
	}
}

func (h *TexWriter) writeResponseParams(w io.Writer, o *api.Operation) {
	if o.HttpResponses.Len() == 0 {
		return
	}

	fmt.Fprintf(w, "<H3>Response</H3>\n")
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Code</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	for _, p := range o.HttpResponses {
		fmt.Fprintf(w, "<TR><TD>%s", p.Name)
		if p.Field.Link() != "" {
			fmt.Fprintf(w, "<br /><I>%s</I>", p.Field.FullLink())
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", p.Field.Description)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}


func (h *TexWriter) WriteResource(r *api.Resource) {
	fn := "_" + conceptFileName(r.Definition) + ".tex"
	path := *api.ConfigDir + "/includes/" + fn
	w, err := os.Create(path)
	defer w.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	dvg := fmt.Sprintf("%s %s %s", r.Name, r.Definition.Version, r.Definition.GroupDisplayName())
	linkID := getLink(dvg)
	fmt.Fprintf(w, "<H1 id=\"%s\">%s</H1>\n", linkID, dvg)

	h.writeSample(w, r.Definition)

	// GVK
	fmt.Fprintf(w, "<TABLE class=\"col-md-8\">\n<THEAD><TR><TH>Group</TH><TH>Version</TH><TH>Kind</TH></TR></THEAD>\n<TBODY>\n")
	fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD></TR>\n",
	            r.Definition.GroupDisplayName(), r.Definition.Version, r.Name)
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")

	if r.DescriptionWarning != "" {
		fmt.Fprintf(w, "<DIV class=\"alert alert-warning col-md-8\"><P><I class=\"fa fa-warning\"></I> <B>Warning:</B></P><P>%s</P></DIV>\n", r.DescriptionWarning)
	}
	if r.DescriptionNote != "" {
		fmt.Fprintf(w, "<DIV class=\"alert alert-info col-md-8\"><I class=\"fa fa-bullhorn\"></I> %s</DIV>\n", r.DescriptionNote)
	}

	h.writeOtherVersions(w, r.Definition)
	h.writeAppearsIn(w, r.Definition)
	h.writeFields(w, r.Definition)

	// Inline
	if r.Definition.Inline.Len() > 0 {
		for _, d := range r.Definition.Inline {
			fmt.Fprintf(w, "<H3 id=\"%s\">%s %s %s</H3>\n", d.LinkID(), d.Name, d.Version, d.Group)
			h.writeAppearsIn(w, d)
			h.writeFields(w, d)
		}
	}

	item := tocItem {
		Level: 1,
		Title: dvg,
		Link: linkID,
		File: fn,
	}
	h.toc.Sections = append(h.toc.Sections, &item)
	h.CurrentSection = &item

	// Operations
	if len(r.Definition.OperationCategories) == 0 {
		return
	}

	for _, c := range r.Definition.OperationCategories {
		if len(c.Operations) == 0 {
			continue
		}
		catID := strings.Replace(strings.ToLower(c.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
		catID = "-strong-" + catID + "-strong-"
		fmt.Fprintf(w, "\\subsection{%s}\n\\label{%s}\n\n", c.Name, catID)
		ocItem := tocItem {
			Level: 2,
			Title: c.Name,
			Link: catID,
		}
		h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &ocItem)

		for _, o := range c.Operations {
			opID := strings.Replace(strings.ToLower(o.Type.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
			fmt.Fprintf(w, "\\subsection{%s}\n\\label{%s}\n", o.Type.Name, opID)
			opItem := tocItem {
				Level: 2,
				Title: o.Type.Name,
				Link: opID,
			}
			ocItem.SubSections = append(ocItem.SubSections, &opItem)

			// Example requests
			requests := o.GetExampleRequests()
			if len(requests) > 0 {
				h.writeOperationSample(w, true, opID, requests)
			}
			// Example responses
			responses := o.GetExampleResponses()
			if len(responses) > 0 {
				h.writeOperationSample(w, false, opID, responses)
			}

			fmt.Fprintf(w, "%s\n\n", o.Description())
			fmt.Fprintf(w, "\\subsubsection{HTTP Request}\n\n")
			fmt.Fprintf(w, "\\texttt{%s}\\n", o.GetDisplayHttp())

			h.writeRequestParams(w, o)
			h.writeResponseParams(w, o)
		}
	}
}

func (h *TexWriter) WriteOldVersionsOverview() {
	writeStaticFile("Old Versions", "_oldversions.tex", h.DefaultStaticContent("Old Versions"))
	item := tocItem {
		Level: 1,
		Title: "OLD API VERSIONS",
		Link: "-strong-old-api-versions-strong-",
		File: "_oldversions.tex",
	}
	h.toc.Sections = append(h.toc.Sections, &item)
	h.CurrentSection = &item
}

func (h *TexWriter) generateDocument() {
	main, err := os.Create(*api.ConfigDir + "/build/main.tex")
	defer main.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	fmt.Fprintf(main, "\\documentclass[10pt,a4paper]{book}\n")
	fmt.Fprintf(main, "\\include{static_includes/package}\n")
	fmt.Fprintf(main, "\\include{static_includes/format}\n")
	fmt.Fprintf(main, "\\begin{document}\n")
	fmt.Fprintf(main, "\n\\frontmatter\n")
	fmt.Fprintf(main, "\\title{%s}\n", h.toc.Title)
	fmt.Fprintf(main, "\\author{Kubernetes Team}\n")
	fmt.Fprintf(main, "\\maketitle\n\n")
	fmt.Fprintf(main, "\\tableofcontents\n\n")
	fmt.Fprintf(main, "\\newpage\n\n")
	fmt.Fprintf(main, "\\mainmatter\n\n")

	// buffer
	for _, sec := range h.toc.Sections {
		fmt.Printf("Collecting %s ... ", sec.File)
		fn := filepath.Join(*api.ConfigDir, "includes", sec.File)
		_, err := ioutil.ReadFile(fn)
		if err == nil {
			fmt.Fprintf(main, "\\include{includes/%s}\n", sec.File)
		}
		fmt.Printf("OK\n")

		for _, sub := range sec.SubSections {
			if len(sub.File) > 0 {
				subfn := filepath.Join(*api.ConfigDir, "includes", sub.File)
				_, err := ioutil.ReadFile(subfn)
				fmt.Printf("Collecting %s ... ", sub.File)
				if err == nil {
					fmt.Fprintf(main, "\\include{includes/%s}\n", sub.File)
					fmt.Printf("OK\n")
				}
			}

			for _, subsub := range sub.SubSections {
				if len(subsub.File) > 0 {
					subsubfn := filepath.Join(*api.ConfigDir, "includes", subsub.File)
					_, err := ioutil.ReadFile(subsubfn)
					fmt.Printf("Collecting %s ... ", subsub.File)
					if err == nil {
						fmt.Fprintf(main, "\\include{includes/%s}\n", subsub.File)
						fmt.Printf("OK\n")
					}
				}
			}
		}
	}

	fmt.Fprintf(main, "\\end{document}\n")
}

func (h *TexWriter) Finalize() {
	os.MkdirAll(*api.ConfigDir + "/build", os.ModePerm)
	h.generateDocument()
}
