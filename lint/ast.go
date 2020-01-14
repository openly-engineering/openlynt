package lint

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func Walk(r *Rule, src string) []error {
	var errs []error

	tfset := token.NewFileSet()
	file, err := parser.ParseFile(tfset, "x.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		err := r.Require.Verify(n)
		if err != nil {
			errs = append(errs, err)
		}

		return true
	})

	return errs
}
