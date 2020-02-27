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
		err := r.Require.Verify(n)
		if err != nil {
			if le, ok := err.(*Violation); ok {
				le.Position = tfset.PositionFor(le.Pos, false)
				le.File = filepath
				le.Rule = r
			}

			if les, ok := err.(*Violations); ok {
				for i := range les.Violations {
					les.Violations[i].Position = tfset.PositionFor(les.Violations[i].Pos, false)
					les.Violations[i].File = filepath
					les.Violations[i].Rule = r
				}
			}

			errs = append(errs, err)
		}

		return true
	})

	return errs
}
