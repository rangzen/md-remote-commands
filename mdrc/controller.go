package mdrc

import (
	"fmt"
	"net/http"
	"os/exec"
	"path"
	"strconv"

	"github.com/go-logr/logr"
)

type Controller struct {
	logger   logr.Logger
	commands *Commands
	html     *HTML
}

func NewController(l logr.Logger, c *Commands, h *HTML) *Controller {
	return &Controller{
		logger:   l.WithName("controller"),
		commands: c,
		html:     h,
	}
}

func (c *Controller) HandleHTML() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Info("serving", "path", r.URL.Path)
		rendered := c.html.Rendered()
		_, err := fmt.Fprint(w, rendered)
		if err != nil {
			c.logger.Error(err, "unable to write", "html", rendered)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (c *Controller) HandleCommand() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		c.logger.Info("serving", "path", request.URL.Path)
		p := path.Base(request.URL.Path)
		i, err := strconv.Atoi(p)
		if err != nil {
			c.logger.Error(err, "base", p)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !c.commands.Valid(i) {
			c.logger.Error(err, "wrong command index", "index", i)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		command := string(c.commands.Command(i))
		out, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			c.logger.Error(err, "command", command)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(out)
		if err != nil {
			c.logger.Error(err, "write", command)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
