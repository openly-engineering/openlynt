package lint

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func Walk(r *Rule, filepath string) []error {
	var errs []error

	tfset := token.NewFileSet()
	file, err := parser.ParseFile(tfset, filepath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		//log.Printf("visiting %#v", n)

		err := r.Require.Verify(n)
		if err != nil {
			if le, ok := err.(*Error); ok {
				le.Position = tfset.PositionFor(le.Pos, false)
			}

			errs = append(errs, err)
		}

		return true
	})

	return errs
}
