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
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/rangzen/md-remote-commands/mdrc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// defaultMarkdown is the Markdown text used when mdrc is run without args.
const defaultMarkdown = `# Welcome in [mdrc](https://github.com/rangzen/md-remote-commands)

This is the default Markdown text because you run ` + "`mdrc`" + ` without argument.

If you need Markdown example files, you can use the commands below.  
Click on the run button and re-run mdrc with the downloaded file.

Thank you for using this tool!

## Learn the ` + "`ls`" + ` command

` + "```mdrc" + `
curl -o learn-ls.md https://raw.githubusercontent.com/rangzen/md-remote-commands/main/examples/learn-ls.md
` + "```" + `

then re-run mdrc with the file:
` + "`mdrc learn-ls`" + `

## System commands
` + "```mdrc" + `
curl -o system.md https://raw.githubusercontent.com/rangzen/md-remote-commands/main/examples/system.md
` + "```" + `

then re-run mdrc with the file:
` + "`mdrc learn-ls`" + `
`

var port string

var rootCmd = &cobra.Command{
	Use:   "mdrc [flags] file.md",
	Short: "Remote commands in Markdown",
	Long:  `Remote commands in Markdown with web GUI.`,
	Run: func(cmd *cobra.Command, args []string) {
		App(args)
	},
}

func App(args []string) {
	log := configureLogger()

	var data []byte
	var err error
	if len(args) == 0 {
		log.Info("no argument, using default text")
		data = []byte(defaultMarkdown)
	} else {
		filepath := args[0]
		log.Info("reading...", "file", filepath)
		data, err = ioutil.ReadFile(filepath)
		if err != nil {
			log.Error(err, "unable to read markdown")
			return
		}
	}

	commands := mdrc.NewCommands(log, data)
	html := mdrc.NewHTML(log, data)
	c := mdrc.NewController(log, commands, html)
	s := mdrc.NewServer(log, c)
	s.Serve(port)
}

// configureLogger configures the Zap implementation of logr.Logger.
func configureLogger() logr.Logger {
	var log logr.Logger

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)
	return log
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&port, "port", "p", "1234", "listening port for web server")
}
