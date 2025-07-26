package build

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/engpetarmarinov/pede/builder"
)

type Options struct {
	OS     string
	ARCH   string
	Input  string
	Output string
	KeepIR bool
	CC     string
}

func Usage() {
	slog.Info(`pede build - Build a .pede file into a native executable

Usage:
  pede build [options] <input.pede>

Options:
  -o <output>     Output binary name (default: input filename without extension)
  --keep-ir       Keep the generated LLVM IR file (default: delete after linking)
  --cc <compiler> Use specified C compiler (clang or gcc, default: clang)
  --os <os>       Operating system target (default: current OS)
  --arch <arch>   Architecture target (default: current architecture)
  --log <level>   Set log level (DEBUG, INFO, WARN, ERROR; default: DEBUG)

Note: pede depends on clang by default to link the generated LLVM IR to a native executable.
`)
}

func Run(opts *Options) {
	builderOpts := &builder.Options{
		OS:     opts.OS,
		ARCH:   opts.ARCH,
		Input:  opts.Input,
		Output: opts.Output,
		KeepIR: opts.KeepIR,
		CC:     opts.CC,
	}
	builder.Build(builderOpts)
}

func Parse(args []string) *Options {
	var opts Options
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	fs.StringVar(&opts.Output, "o", "", "output binary name")
	fs.BoolVar(&opts.KeepIR, "keep-ir", false, "keep the generated LLVM IR file")
	fs.StringVar(&opts.OS, "os", "", "target operating system (default: current OS)")
	fs.StringVar(&opts.ARCH, "arch", "", "target architecture (default: current architecture)")
	fs.StringVar(&opts.CC, "cc", "clang", "C compiler to use (clang or gcc)")
	fs.Usage = Usage
	err := fs.Parse(args)
	if err != nil {
		slog.Error("Error parsing flags", "err", err)
		Usage()
	}

	if fs.NArg() < 1 {
		slog.Error("Usage: pede build [options] <input.pede>")
		Usage()
		os.Exit(1)
	}
	opts.Input = flag.Arg(len(flag.Args()) - 1)
	if opts.Input == "" {
		slog.Error("No input file specified. Use <input.pede> to specify the input file.")
		Usage()
		os.Exit(1)
	}
	if opts.Output == "" {
		base := filepath.Base(opts.Input)
		opts.Output = strings.TrimSuffix(base, filepath.Ext(base))
		if opts.Output == "" {
			slog.Error("Could not determine output file name. Use -o to specify output.")
			Usage()
			os.Exit(1)
		}
	}
	return &opts
}
