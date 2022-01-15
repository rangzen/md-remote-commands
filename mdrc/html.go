/*
Copyright © 2022 Cédric L’homme <public@l-homme.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package mdrc

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"github.com/go-logr/logr"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

// content holds our static web server content.
//go:embed template
var content embed.FS

const (
	mdrcCodeBlockIdentifier = "mdrc"
)

type HTML struct {
	logger logr.Logger
	html   string
}

func NewHTML(l logr.Logger, data []byte) *HTML {
	logger := l.WithName("html")
	h := &HTML{
		logger: logger,
	}
	h.render(data)
	return h
}

// render renders the HTML code from the Markdown text.
func (h *HTML) render(data []byte) {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: h.renderShellCodeBlock(),
	}
	renderer := html.NewRenderer(opts)
	m := string(markdown.ToHTML(data, nil, renderer))

	tmpl, err := template.ParseFS(content, "template/root.html")
	if err != nil {
		log.Fatalf("parsing template: %v\n", err)
	}

	sb := &strings.Builder{}
	d := struct {
		Markdown string
	}{
		Markdown: m,
	}
	err = tmpl.Execute(sb, d)
	if err != nil {
		log.Fatalf("executing template: %v\n", err)
	}

	h.html = sb.String()
	h.logger.Info("rendered", "size", len(h.html))
}

// renderShellCodeBlock injects the mdrc HTML code to deal with the commands
// inside CodBlock with mdrc identifier (ast.CodeBlock.Info).
func (h *HTML) renderShellCodeBlock() func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	original := html.NewRenderer(html.RendererOptions{})
	var cmdCount int
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if _, ok := node.(*ast.CodeBlock); !ok {
			return ast.GoToNext, false
		}
		codeBlock := node.(*ast.CodeBlock)
		original.CodeBlock(w, codeBlock)
		if bytes.Equal(codeBlock.Info, []byte(mdrcCodeBlockIdentifier)) {
			_, err := w.Write(
				[]byte(
					fmt.Sprintf(
						"<button type=\"submit\" onclick=\"Run(%d)\">Run</button>\n<pre><code class=\"language-shell\" id=\"command-%d\"></code></pre>\n",
						cmdCount,
						cmdCount,
					)))
			if err != nil {
				h.logger.Error(err, "codeBlock", "info", codeBlock.Info, "literal", codeBlock.Literal)
			}
			cmdCount++
		}
		return ast.GoToNext, true
	}
}

// Rendered returns the final HTML rendered with mdrc injections.
func (h *HTML) Rendered() string {
	return h.html
}
