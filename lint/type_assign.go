package lint

import (
	"bytes"
	"fmt"
	"go/ast"
	"text/template"

	"github.com/Masterminds/sprig"
)

type assignLhs struct {
	Name *regexpStr
}

func (as *assignLhs) Match(n ast.Node) bool {
	stmt := n.(*ast.AssignStmt)

	if as.Name != nil {
		return as.Name.MatchString(stmt.Lhs[0].(*ast.Ident).Name)
	}

	return false
}

type assignRule struct {
	If struct {
		Lhs *assignLhs
		//Rhs *regexpStr
	}

	Require struct {
		Lhs *stringRequirement
		//Rhs *stringRequirement
	}
}

func (ar *assignRule) Verify(n ast.Node) error {
	as, ok := n.(*ast.AssignStmt)
	if !ok {
		return nil
	}

	ifr, reqr := ar.If, ar.Require

	if ifr.Lhs != nil && ifr.Lhs.Match(n) {
		buf := new(bytes.Buffer)
		must := reqr.Lhs.Template

		path := as.Lhs[0].(*ast.Ident).Name
		idx := ifr.Lhs.Name.FindStringIndex(path)
		path = path[idx[0]:idx[1]]

		must = ifr.Lhs.Name.ReplaceAllString(path, must)

		tpl, err := template.New("assignRule.VariableName").Funcs(sprig.TxtFuncMap()).Parse(must)
		if err != nil {
			return err
		}

		name := as.Lhs[0].(*ast.Ident).Name
		if err := tpl.Execute(buf, must); err != nil {
			return err
		}

		if name != buf.String() {
			return &Error{
				Message: fmt.Sprintf(`expected lhs %s to be named "%s", but it was "%s"`, path, buf.String(), name),
				Pos:     n.Pos(),
			}
		}
	}

	return nil
}
