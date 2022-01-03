# Markdown Remote Commands (mdrc)

Utility tool to expose through a web server some commands defined in a Markdown file.

## Installation

```shell
go install github.com/rangzen/md-remote-commands/cmd/mdrc@latest
```

## Usage

* Get an example from [the examples' directory](https://github.com/rangzen/md-remote-commands/tree/main/examples).
  E.g., `system.md` with `wget https://raw.githubusercontent.com/rangzen/md-remote-commands/main/examples/system.md`.
* Run `mdrc system.md`.
* Open your navigator to the system with the correct port (1234 by default).