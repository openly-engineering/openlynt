package lint

import (
	"regexp"

	yaml "gopkg.in/yaml.v3"
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

type stringRequirement struct {
	Template string
}

func (s *stringRequirement) UnmarshalYAML(n *yaml.Node) error {
	str := ""
	n.Decode(&str)

	s.Template = str

	return nil
}
