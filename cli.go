package appdater

import (
	"io"
)

type CLI struct {
	OutStream, ErrStream io.Writer
}

func (c *CLI) Run(args []string) int {
	return 0
}
