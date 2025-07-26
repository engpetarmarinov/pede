package cli

import (
	"flag"
	"log/slog"
	"os"

	"github.com/engpetarmarinov/pede/cli/cmds/build"
)

type Options struct {
	Cmd      string
	LogLevel string
}

func Usage() {
	slog.Info(`pede - a simple programming language

Usage: pede [options] <command> [arguments]

Options:
  -log <level>         Set log level (DEBUG, INFO, WARN, ERROR; default: DEBUG)

Commands:
  build <input.pede>   Build the specified .pede file
  help                 Show this help message
`)
}

func Parse() *Options {
	var opts Options
	flag.StringVar(&opts.LogLevel, "log", "DEBUG", "log level: DEBUG, INFO, WARN, ERROR")
	flag.Usage = Usage
	flag.Parse()
	opts.Cmd = flag.Arg(0)
	return &opts
}

func Run(opts *Options) {
	switch opts.Cmd {
	case "build":
		buildOpts := build.Parse(flag.Args()[1:])
		build.Run(buildOpts)
	case "help":
		Usage()
	default:
		Usage()
		os.Exit(1)
	}
}
