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
	"io"

	"github.com/go-logr/logr"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

type Command []byte

type Commands struct {
	logger   logr.Logger
	commands []Command
}

// NewCommands creates a Commands structure from extracted commands from data.
func NewCommands(l logr.Logger, data []byte) *Commands {
	f := &Commands{
		logger:   l.WithName("commands"),
		commands: make([]Command, 0),
	}
	f.extract(data)
	return f
}

// extract extracts commands from a Markdown text.
func (f *Commands) extract(data []byte) {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: f.extractShellCodeBlock(),
	}
	renderer := html.NewRenderer(opts)
	markdown.ToHTML(data, nil, renderer)
	f.logger.Info("extracted", "count", len(f.commands))
}

// extractShellCodeBlock is a special RenderNodeHook that intercepts commands
// during "fake" rendering of the Markdown text.
func (f *Commands) extractShellCodeBlock() func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if _, ok := node.(*ast.CodeBlock); !ok {
			return ast.GoToNext, false
		}
		codeBlock := node.(*ast.CodeBlock)
		if bytes.Equal(codeBlock.Info, []byte("mdrc")) {
			f.commands = append(f.commands, codeBlock.Literal)
		}
		return ast.GoToNext, true
	}
}

// Valid returns true if the index i is a correct index for commands.
func (f *Commands) Valid(i int) bool {
	return len(f.commands) != 0 && i >= 0 && i < len(f.commands)
}

// Command returns the command at index i.
func (f *Commands) Command(i int) Command {
	return f.commands[i]
}
