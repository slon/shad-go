package jokelint

import (
	"go/ast"
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "jokelint",
	Doc:  "check for outdated jokes about go",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	ins := inspector.New(pass.Files)

	// We filter only function calls.
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		visitCall(pass, call)
	})

	return nil, nil
}

func visitCall(pass *analysis.Pass, call *ast.CallExpr) {
	fn, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	def := pass.TypesInfo.Uses[fn.Sel]
	if x, ok := def.(*types.Func); !ok {
		return
	} else if x.Pkg().Path() != "fmt" || x.Name() != "Println" {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg := call.Args[0]
	argTyp := pass.TypesInfo.Types[arg]
	if argTyp.Value != nil && constant.StringVal(argTyp.Value) == "lol no generics" {
		pass.Reportf(call.Pos(), "outdated joke")
	}
}
