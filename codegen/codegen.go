package codegen

import (
	"fmt"
	"io"

	"github.com/engpetarmarinov/pede/ast"
	"github.com/engpetarmarinov/pede/lexer"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Codegen struct {
	mod           *ir.Module
	block         *ir.Block
	vars          map[string]*ir.InstAlloca
	fmtStrGlobal  *ir.Global            // cache for float format string global
	fmtStrSGlobal *ir.Global            // cache for string format string global
	strGlobals    map[string]*ir.Global // cache for string literals
}

// NewCodegen initializes a new Codegen instance with a module and entry block.
func NewCodegen(os, arch string) *Codegen {
	mod := ir.NewModule()
	mod.TargetTriple = getTargetTriple(os, arch)
	mainFn := mod.NewFunc("main", types.Void)
	entry := mainFn.NewBlock("entry")
	return &Codegen{
		mod:        mod,
		block:      entry,
		vars:       make(map[string]*ir.InstAlloca),
		strGlobals: make(map[string]*ir.Global),
	}
}

func (cg *Codegen) WriteTo(w io.Writer) (n int64, err error) {
	n, err = cg.mod.WriteTo(w)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// GenStmt dispatches codegen for statements
func (cg *Codegen) GenStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.Assignment:
		cg.GenAssign(s)
	case *ast.PrintStmt:
		cg.GenPrint(s)
	default:
		panic("unsupported statement type")
	}
}

// GenPrint emits code for print(x)
func (cg *Codegen) GenPrint(p *ast.PrintStmt) {
	printf := cg.getOrDeclarePrintf()
	val := cg.genExpr(p.Expr)

	switch val.Type().String() {
	case types.Double.String():
		if cg.fmtStrGlobal == nil {
			cg.fmtStrGlobal = cg.mod.NewGlobalDef(".fmtstr", constant.NewCharArrayFromString("%f\n\x00"))
		}
		arrayType := cg.fmtStrGlobal.Init.(*constant.CharArray).Typ
		fmtPtr := cg.block.NewGetElementPtr(arrayType, cg.fmtStrGlobal, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		cg.block.NewCall(printf, fmtPtr, val)
	case types.I8Ptr.String():
		if cg.fmtStrSGlobal == nil {
			cg.fmtStrSGlobal = cg.mod.NewGlobalDef(".fmtstr_s", constant.NewCharArrayFromString("%s\n\x00"))
		}
		arrayType := cg.fmtStrSGlobal.Init.(*constant.CharArray).Typ
		fmtPtr := cg.block.NewGetElementPtr(arrayType, cg.fmtStrSGlobal, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		cg.block.NewCall(printf, fmtPtr, val)
	default:
		panic("unsupported print type: " + val.Type().String())
	}
}

func (cg *Codegen) GenAssign(a *ast.Assignment) {
	exprVal := cg.genExpr(a.Expr)
	var alloca *ir.InstAlloca
	switch exprVal.Type().String() {
	case types.Double.String():
		alloca = cg.block.NewAlloca(types.Double)
		cg.block.NewStore(exprVal, alloca)
	case types.I8Ptr.String():
		alloca = cg.block.NewAlloca(types.I8Ptr)
		cg.block.NewStore(exprVal, alloca)
	default:
		panic("unsupported assignment type: " + exprVal.Type().String())
	}
	cg.vars[a.Name] = alloca
}

func (cg *Codegen) genExpr(e ast.Expr) value.Value {
	switch n := e.(type) {
	case *ast.Number:
		return constant.NewFloat(types.Double, n.Value)
	case *ast.String:
		if g, ok := cg.strGlobals[n.Value]; ok {
			// Reuse global if already created
			arrayType := g.Init.(*constant.CharArray).Typ
			return cg.block.NewGetElementPtr(arrayType, g, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		}
		g := cg.mod.NewGlobalDef(".str."+fmt.Sprintf("%x", len(cg.strGlobals)), constant.NewCharArrayFromString(n.Value+"\x00"))
		cg.strGlobals[n.Value] = g
		arrayType := g.Init.(*constant.CharArray).Typ
		return cg.block.NewGetElementPtr(arrayType, g, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	case *ast.Variable:
		ptr := cg.vars[n.Name]
		return cg.block.NewLoad(ptr.ElemType, ptr)
	case *ast.Binary:
		lhs := cg.genExpr(n.Left)
		rhs := cg.genExpr(n.Right)
		switch n.Op {
		case lexer.TokenPlus:
			if lhs.Type().String() == types.Double.String() && rhs.Type().String() == types.Double.String() {
				return cg.block.NewFAdd(lhs, rhs)
			}
			// Optionally, support string concatenation here
			panic("string concatenation not supported yet")
		case lexer.TokenStar:
			return cg.block.NewFMul(lhs, rhs)
		default:
			panic("unsupported operator: " + n.Op)
		}
	default:
		panic("unknown expression node")
	}
}

// Add helper to get or declare printf
func (cg *Codegen) getOrDeclarePrintf() *ir.Func {
	for _, fn := range cg.mod.Funcs {
		if fn.Name() == "printf" {
			return fn
		}
	}
	printf := cg.mod.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
	printf.Sig.Variadic = true
	return printf
}

func (cg *Codegen) Finish() {
	cg.block.NewRet(nil)
}

// GenProgram emits code for a program (list of statements)
func (cg *Codegen) GenProgram(prog *ast.Program) {
	for _, stmt := range prog.Stmts {
		cg.GenStmt(stmt)
	}
}
