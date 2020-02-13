package lint

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

type assignLhs struct {
	Name *regexpStr
}

func (as *assignLhs) Match(n ast.Node) bool {
	stmt := n.(*ast.AssignStmt)

	if as.Name != nil {
		for i := range stmt.Lhs {
			lhs := stmt.Lhs[i]

			// only care if it's an identity, nothing else matters ATM
			if ident, ok := lhs.(*ast.Ident); ok {
				return as.Name.MatchString(ident.Name)
			}
		}
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

	var errstrs []string
	ifr, reqr := ar.If, ar.Require

	if ifr.Lhs != nil && ifr.Lhs.Match(n) {
		buf := new(bytes.Buffer)
		must := reqr.Lhs.Template

		for i := range as.Lhs {
			buf.Reset()

			varname := as.Lhs[i].(*ast.Ident).Name
			idx := ifr.Lhs.Name.FindStringIndex(varname)
			varname = varname[idx[0]:idx[1]]
			must := ifr.Lhs.Name.ReplaceAllString(varname, must)

			tpl, err := template.New("assignRule.VariableName").Funcs(sprig.TxtFuncMap()).Parse(must)
			if err != nil {
				return err
			}

			name := as.Lhs[i].(*ast.Ident).Name
			if err := tpl.Execute(buf, must); err != nil {
				return err
			}

			if name != buf.String() {
				errstrs = append(errstrs, fmt.Sprintf(`expected lhs %s to be named "%s", but it was "%s"`, varname, buf.String(), name))
			}
		}
	}

	if len(errstrs) > 0 {
		return &Error{
			Message: strings.Join(errstrs, "; "),
			Pos:     n.Pos(),
		}
	}

	return nil
}
