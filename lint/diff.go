package lint

import (
	"github.com/golangci/revgrep"
)

// similar to golangci-lint, use revgrep to determine what commit a change came
// from.
func FilterByRevision(vs *Violations, from, to string) (*Violations, error) {
	gp, _, err := revgrep.GitPatch(from, to)
	if err != nil {
		return nil, err
	}

	c := &revgrep.Checker{
		RevisionFrom: from,
		RevisionTo:   to,
		Patch:        gp,
	}

	if err := c.Prepare(); err != nil {
		return nil, err
	}

	filtered := &Violations{
		Violations: make([]*Violation, 0, len(vs.Violations)),
	}

	for i := range vs.Violations {
		v := vs.Violations[i]

		if _, isNew := c.IsNewIssue(v); isNew {
			filtered.Violations = append(filtered.Violations, v)
		}
	}

	return filtered, nil
}
