package lint

import (
	"io"

	"github.com/golangci/revgrep"
)

// similar to golangci-lint, use revgrep to determine what commit a change came
// from.
func FilterByRevision(patch io.Reader, newFiles []string, vs *Violations) (*Violations, error) {
	c := &revgrep.Checker{
		Patch:    patch,
		NewFiles: newFiles,
	}

	if err := c.Prepare(); err != nil {
		return nil, err
	}

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
