package lint

import (
	"fmt"
	"go/ast"
	"strings"
)

// commentGroupRule operates on a group of comments.
// For instance, _this_ comment is composed of 5 *ast.Comment in a single
// *ast.CommentGroup.
// There may be a case for grouped and non-grouped implementations, but at the
// moment we have no need for it.
type commentGroupRule struct {
	If struct {
		Text *regexpStr
		// could add a "type" here for `//` vs `/* */`
		// rare to see `/*` being used in Go
	}

	Require struct {
		Text *regexpStr
		Len  int
	}
}

func (r *commentGroupRule) Verify(n ast.Node) error {
	// To get _all_ comments, we have to fetch it from *ast.File.
	// specifying *ast.CommentGroup results in only doc comments
	// being parsed. /shrug
	astfi, ok := n.(*ast.File)
	if !ok {
		return nil
	}

	errs := &ErrorCollection{}
	for i := range astfi.Comments {
		cg := astfi.Comments[i]

		if r.Require.Text != nil {
			if !r.If.Text.MatchString(cg.Text()) {
				continue
			}

			// NOTE(@chrsm): *ast.Comment does not specify whether
			// it's //-style or /*-style and `.Text()` does NOT
			// return the beginning or end. It's not worth it to
			// print the ast and check at this position, so to not
			// fail every /*-style comment based on length, we
			// check for `\n`s in the text.
			// For //-style comments, this should be 0.
			lines := len(cg.List)
			if ncount := strings.Count(cg.Text(), "\n"); lines != ncount {
				lines = ncount
			}

			if r.Require.Len > 0 && lines < r.Require.Len {
				errs.Errors = append(errs.Errors, &Error{
					Message: fmt.Sprintf(`expected at least %d lines of context but have %d`, r.Require.Len, lines),
					Pos:     cg.Pos(),
				})

			}

			if !r.Require.Text.MatchString(cg.Text()) {
				errs.Errors = append(errs.Errors, &Error{
					Message: fmt.Sprintf(`expected comment to match /%s/`, r.Require.Text.String()),
					Pos:     cg.Pos(),
				})
			}
		}
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}
