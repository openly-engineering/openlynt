package lint

import (
	"context"

	"github.com/golangci/revgrep"
)

// NewRevisionChecker builds a revgrep checker for the given revision range and
// loads the underlying git patch.
func NewRevisionChecker(ctx context.Context, from, to string) (*revgrep.Checker, error) {
	c := &revgrep.Checker{
		RevisionFrom: from,
		RevisionTo:   to,
	}

	if err := c.Prepare(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

// similar to golangci-lint, use revgrep to determine what commit a change came
// from.
func FilterByRevision(c *revgrep.Checker, vs *Violations) (*Violations, error) {
	filtered := &Violations{
		Violations: make([]*Violation, 0, len(vs.Violations)),
	}

	for i := range vs.Violations {
		v := vs.Violations[i]

		// files that are "new" - no history, whether unstaged or
		// otherwise - will _not_ be filtered out. revgrep will not detect
		// these in by nature of how .GitPatch pulls changes
		noHist := false
		for k := range c.NewFiles {
			if v.File == c.NewFiles[k] {
				filtered.Violations = append(filtered.Violations, v)

				noHist = true
				break
			}
		}

		if noHist {
			continue
		}

		if _, isNew := c.IsNewIssue(v); isNew {
			filtered.Violations = append(filtered.Violations, v)
		}
	}

	return filtered, nil
}
