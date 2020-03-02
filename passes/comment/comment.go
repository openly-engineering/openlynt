package comment

import "regexp"

type (
	If struct {
		Text *regexp.Regexp
	}

	Require struct {
		Text *regexp.Regexp
		Len  int
	}

	Rule struct {
		If      If
		Require Require
	}
)
