package main

import (
	"github.com/engpetarmarinov/pede/cli"
	"github.com/engpetarmarinov/pede/logutil"
)

func main() {
	opts := cli.Parse()
	logutil.Setup(opts.LogLevel)
	cli.Run(opts)
}
