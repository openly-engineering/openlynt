package comment

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `report violations of predefined comment rules

The comment analyzer checks code comments that match a regexp and ensures
that they match another regexp.

This is useful for comments like "TODO", "FIXME", or anything you may use in
your own projects.

For instance, you may require that a TODO comment contains a link to an issue
tracker and additionally require that they have a certain amount of context
that explains the technical information that is relevant to the code.`

var Analyzer = &analysis.Analyzer{
	Name: "comment",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	RunDespiteErrors: true,
}

var _nodeTypes = []ast.Node{
	// we actually use a *ast.File instead of *ast.Comment or
	// *ast.CommentGroup as `go/parser` only returns doc comments from code
	// otherwise.
	(*ast.File)(nil),
}

func run(ap *analysis.Pass) (interface{}, error) {
	// TODO(@chrsm): we need to find a way to pass rules down here.
	// For now, we're going to fake some. see _tmpRules.
	// I *think* what we may need to do port cmd/openlynt to use `checker.Run`,
	// and pass the yaml config values as flags.

	astins := ap.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	astins.Preorder(_nodeTypes, func(n ast.Node) {
		fi := n.(*ast.File)

		for i := range fi.Comments {
			cg := fi.Comments[i]

			for j := range _tmpRules {
				r := _tmpRules[j]
				txt := cg.Text()

				if !r.If.Text.MatchString(txt) {
					continue
				}

				// NOTE(@chrsm): *ast.Comment does not specify whether
				// it's //-style or /*-style and `.Text()` does NOT
				// return the beginning or end. It's not worth it to
				// print the ast and check at this position, so to not
				// fail every /*-style comment based on length, we
				// check for `\n`s in the text.
				// For //-style comments, this should be 0.
				lines := len(cg.List)
				if ncount := strings.Count(txt, "\n"); lines < ncount {
					lines = ncount
				}

				if r.Require.Len > 0 && lines < r.Require.Len {
					ap.Reportf(
						cg.Pos(),
						"expected at least %d lines of comments around TODO, but have %d",
						r.Require.Len,
						lines,
					)
				}

				if !r.Require.Text.MatchString(txt) {
					ap.Reportf(
						cg.Pos(),
						"expected comment to contain /%s/",
						r.Require.Text.String(),
					)
				}
			}
		}
	})

	return nil, nil
}

var _tmpRules = []Rule{
	Rule{
		If: If{
			Text: regexp.MustCompile(`(TODO|FIXME|XXX)`),
		},

		Require: Require{
			Text: regexp.MustCompile(`http:\/\/`),
			Len:  2,
		},
	},
}
