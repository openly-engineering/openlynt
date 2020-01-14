package lint

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"go/ast"
	"text/template"

	"github.com/Masterminds/sprig"
	yaml "gopkg.in/yaml.v3"
)

const (
	RuleImport   = "import"
	RuleFunction = "func"
)

type regexpStr struct {
	*regexp.Regexp
}

func (r *regexpStr) UnmarshalYAML(n *yaml.Node) error {
	str := ""
	n.Decode(&str)

	r.Regexp = regexp.MustCompile(str)

	return nil
}

type Requirement interface {
	Verify(ast.Node) error
}

type stringRequirement struct {
	Template string
}

type importRequirement struct {
	X string

	Regexp *regexpStr

	Name *stringRequirement
}

func (i *importRequirement) Verify(n ast.Node) error {
	in, ok := n.(*ast.ImportSpec)
	if !ok {
		return nil
	}

	if i.Name != nil {
		sm := i.Regexp.FindStringSubmatch(in.Path.Value)
		if len(sm) == 0 {
			return nil
		}

		buf := new(bytes.Buffer)
		must := i.Name.Template

		path := in.Path.Value
		idx := i.Regexp.FindStringIndex(path)
		path = path[idx[0]:idx[1]]

		must = i.Regexp.ReplaceAllString(path, must)

		tpl := template.Must(
			template.New("replacer").Funcs(sprig.TxtFuncMap()).Parse(must),
		)

		n := ""
		if in.Name != nil {
			n = in.Name.Name
		}

		tpl.Execute(buf, must)
		if n != buf.String() {
			return errors.New(fmt.Sprintf("expected import(%s) to be named(%s), but it was(%s)", in.Path.Value, buf.String(), n))
		}
	}

	return nil
}

type Rule struct {
	Name, Type string

	Require Requirement
}

var _ yaml.Unmarshaler = (*Rule)(nil)

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

	switch typ {
	case RuleImport:
		req := &importRequirement{}

		marshalhack(m["require"], req)

		r.Require = req
	default:
		panic(fmt.Sprintf("node rule type spec(%s) unknown on line %d", typ, n.Line))
	}

	return nil
}

func marshalhack(src, dst interface{}) error {
	b, err := yaml.Marshal(src)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, dst)
}
