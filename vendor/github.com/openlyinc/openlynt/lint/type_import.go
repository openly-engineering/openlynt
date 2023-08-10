package lint

import (
	"bytes"
	"fmt"
	"go/ast"
	"text/template"

	"github.com/Masterminds/sprig"
)

type importRule struct {
	If struct {
		Path *regexpStr
	}

	Require struct {
		Name *stringRequirement
	}
}

func (i *importRule) Verify(n ast.Node) error {
	in, ok := n.(*ast.ImportSpec)
	if !ok {
		return nil
	}

	ifr, reqr := i.If, i.Require

	if reqr.Name != nil {
		sm := ifr.Path.FindStringSubmatch(in.Path.Value)
		if len(sm) == 0 {
			return nil
		}

		buf := new(bytes.Buffer)
		must := reqr.Name.Template

		path := in.Path.Value
		idx := ifr.Path.FindStringIndex(path)
		path = path[idx[0]:idx[1]]

		must = ifr.Path.ReplaceAllString(path, must)

		tpl, err := template.New("importRule.Name").Funcs(sprig.TxtFuncMap()).Parse(must)
		if err != nil {
			return err
		}

		name := ""
		if in.Name != nil {
			name = in.Name.Name
		}

		if err := tpl.Execute(buf, must); err != nil {
			return err
		}

		if name != buf.String() {
			return &Violation{
				Message: fmt.Sprintf(`expected import %s to be named "%s", but it was "%s"`, in.Path.Value, buf.String(), name),
				Pos:     n.Pos(),
			}
		}
	}

	return nil
}
