package main

import (
	"os"

	"github.com/shimastripe/appdater"
)

func main() {
	cli := &appdater.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
