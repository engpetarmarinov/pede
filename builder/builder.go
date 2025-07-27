package builder

import (
	"io"
	"log/slog"
	"os"
	"os/exec"

	"github.com/engpetarmarinov/pede/ast"
	"github.com/engpetarmarinov/pede/codegen"
	"github.com/engpetarmarinov/pede/lexer"
	"github.com/engpetarmarinov/pede/parser"
	"github.com/engpetarmarinov/pede/preprocessor"
)

// Preprocess preprocesses the input string, strips comments and empty/unknown lines, and returns a cleaned string.
func Preprocess(input string) (string, error) {
	filtered, err := preprocessor.Preprocess(input, preprocessor.DefaultRules())
	return filtered, err
}

// Lex lexes the input string and returns a lexer instance
func Lex(input string) *lexer.Lexer {
	return lexer.NewLexer(input)
}

// Parse parses the input using the lexer and returns an AST
func Parse(lx *lexer.Lexer) *ast.Program {
	p := parser.NewParser(lx)
	astProgram, err := p.Parse()
	if err != nil {
		slog.Error("builder parse failed", "err", err)
		os.Exit(1)
	}
	return astProgram
}

// Codegen generates LLVM IR from the AST
func Codegen(ast *ast.Program, buildOS, buildARCH string) *codegen.Codegen {
	cg := codegen.NewCodegen(buildOS, buildARCH)
	cg.GenProgram(ast)
	cg.Finish()
	return cg
}

// WriteIR writes the generated IR to a file
func WriteIR(cg *codegen.Codegen, output string) (string, error) {
	irFile := output + ".ll"
	slog.Debug("Writing IR", "file", irFile)
	irf, err := os.Create(irFile)
	if err != nil || irf == nil {
		return "", err
	}
	_, err = cg.WriteTo(irf)
	irf.Close()
	if err != nil {
		return "", err
	}
	return irFile, nil
}

// Link compiles and links the generated IR file into an executable
func Link(cc, irFile, output string) error {
	cmd := exec.Command(cc, irFile, "-o", output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type Options struct {
	OS     string // Target operating system
	ARCH   string // Target architecture
	Input  string // Input .pede file
	Output string // Output binary name
	KeepIR bool   // Whether to keep the generated LLVM IR file
	CC     string // C compiler to use (default: clang)
}

// Build orchestrates the build process
func Build(opts *Options) {
	f, err := os.Open(opts.Input)
	if err != nil {
		slog.Error("failed to open input", "err", err)
		os.Exit(1)
	}
	code, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		slog.Error("failed to read input", "err", err)
		os.Exit(1)
	}
	preprocessed, err := Preprocess(string(code))
	if err != nil {
		slog.Error("preprocess failed", "err", err)
		os.Exit(1)
	}
	lx := Lex(preprocessed)
	program := Parse(lx)
	cg := Codegen(program, opts.OS, opts.ARCH)
	irFile, err := WriteIR(cg, opts.Output)
	if err != nil {
		slog.Error("failed to write IR", "err", err)
		os.Exit(1)
	}
	if err := Link(opts.CC, irFile, opts.Output); err != nil {
		slog.Error("failed to link executable", "err", err)
		os.Exit(1)
	}
	if !opts.KeepIR {
		os.Remove(irFile)
	}
	slog.Info("pede was built", "OS", opts.OS, "ARCH", opts.ARCH)
}
