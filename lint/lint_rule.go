package lint

import (
	"fmt"

	"go/ast"

	yaml "gopkg.in/yaml.v3"
)

var (
	_ yaml.Unmarshaler = (*Rule)(nil)
	_ Requirement      = (*importRule)(nil)
)

const (
	RuleImport   = "import"
	RuleAssign   = "assignment"
	RuleFunction = "func"
)

type Requirement interface {
	Verify(ast.Node) error
}

type Rule struct {
	Name, Type string

	Require Requirement
}

func (r *Rule) UnmarshalYAML(n *yaml.Node) error {
	m := map[string]interface{}{}

	if err := n.Decode(&m); err != nil {
		panic(err)
	}

	typ, ok := m["type"].(string)
	if !ok {
		panic("node missing type spec")
	}

	r.Type = typ

	name, ok := m["name"].(string)
	if !ok {
		panic("node missing name")
	}

	r.Name = name

	switch typ {
	case RuleImport:
		req := &importRule{}
		marshalhack(m, req)

		r.Require = req
	case RuleAssign:
		req := &assignRule{}
		marshalhack(m, req)
		r.Require = req
	default:
		panic(fmt.Sprintf("node rule type spec(%s) unknown on line %d", typ, n.Line))
	}

	return nil
}
