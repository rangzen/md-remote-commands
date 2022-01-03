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
	"fmt"
	"io"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

const (
	head = `
<head>
	<title>Markdown Remote Commands</title>
</head>
`
	javascript = `
<script type="text/javascript" >
	function Run(num) {
		fetch('/run/'+num).then(response =>{
    		return response.text();
		}).then(data =>{
			document.getElementById('command-'+num).innerHTML=data;
		})
	}
</script>
`
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

func (h *HTML) render(data []byte) {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: h.renderShellCodeBlock(),
	}
	renderer := html.NewRenderer(opts)
	sb := strings.Builder{}
	sb.WriteString("<html>")
	sb.WriteString(head)
	sb.WriteString("<body>")
	sb.WriteString(javascript)
	sb.Write(markdown.ToHTML(data, nil, renderer))
	sb.WriteString(`
</body>
</html>
`)
	h.html = sb.String()
	h.logger.Info("rendered", "size", len(h.html))
}

func (h *HTML) renderShellCodeBlock() func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	original := html.NewRenderer(html.RendererOptions{})
	var cmdCount int
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if _, ok := node.(*ast.CodeBlock); !ok {
			return ast.GoToNext, false
		}
		codeBlock := node.(*ast.CodeBlock)
		original.CodeBlock(w, codeBlock)
		if bytes.Equal(codeBlock.Info, []byte("mdrc")) {
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

func (h *HTML) Rendered() string {
	return h.html
}
