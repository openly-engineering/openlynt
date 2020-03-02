package namedimport

import "regexp"

type (
	If struct {
		Path *regexp.Regexp
	}

	Require struct {
		Name string
	}

	Rule struct {
		If      If
		Require Require
	}
)
