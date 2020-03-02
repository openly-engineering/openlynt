package namedimport

import (
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `report violations of predefined named import rules

The namedimport analyzer checks import statements that match a regexp and
ensures that they are named properly.`

var Analyzer = &analysis.Analyzer{
	Name: "namedimport",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	RunDespiteErrors: true,
}

var _nodeTypes = []ast.Node{
	(*ast.ImportSpec)(nil),
}

func run(ap *analysis.Pass) (interface{}, error) {
	// TODO(@chrsm): we need to find a way to pass rules down here.
	// For now, we're going to fake some. see _tmpRules.
	// I *think* what we may need to do port cmd/openlynt to use `checker.Run`,
	// and pass the yaml config values as flags.
	astins := ap.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	astins.Preorder(_nodeTypes, func(n ast.Node) {
		imp := n.(*ast.ImportSpec)

		for i := range _tmpRules {
			r := _tmpRules[i]

			if !r.If.Path.MatchString(imp.Path.Value) {
				continue
			}

			checkimp(ap, r, imp)
		}
	})

	return nil, nil
}

func checkimp(ap *analysis.Pass, r Rule, n *ast.ImportSpec) {
	// TODO(@chrsm): we need to support the same templating scheme
	// currently implemented in master.
	expect := r.Require.Name

	if n.Name == nil {
		ap.Reportf(n.Pos(), "expected %s, but import has no name", expect)

		return
	}

	if n.Name.Name != expect {
		ap.Reportf(n.Name.Pos(), "expected %s, but import was named %s", expect, n.Name.Name)

		return
	}
}

var _tmpRules = []Rule{
	Rule{
		If: If{
			Path: regexp.MustCompile(`/pkgxyz\/v(?P<version>\d+)`),
		},
		Require: Require{
			Name: "pkgxyzV1",
		},
	},
}
