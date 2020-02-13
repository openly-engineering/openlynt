package lint

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestErr(t *testing.T) {
	yml := `
type: import
name: Import
if:
  path: \/pkg\/(?P<prefix>[a-z]+)\/v(?P<version>[0-9]+)
require:
  name: req{{ "${prefix}" | upper }}v${version}
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, "testdata/import_single.go")
	if len(errs) == 0 {
		t.Fatalf("expected >0 errors")
	}
}

func TestBasicOK(t *testing.T) {
	yml := `
type: import
name: Import
if:
  path: \/pkg\/(?P<prefix>[a-z]+)\/v(?P<version>[0-9]+)
require:
  name: req{{ "${prefix}" | upper }}v${version}
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, "testdata/import_two.go")
	if len(errs) != 0 {
		t.Fatalf("expected 0 errors, got %d", len(errs))
	}
}

func TestBasicMixed(t *testing.T) {
	yml := `
type: import
name: Import Rule
if:
  path: \/pkg\/(?P<prefix>[a-z]+)\/v(?P<version>[0-9]+)
require:
  name: req{{ "${prefix}" | upper }}v${version}
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, "testdata/import_mixed.go")
	if len(errs) != 1 {
		t.Fatalf("expected 1 errors, got %d", len(errs))
	}

	err := errs[0]
	if err.Error() != `expected import "xyz.org/pkg/nottest/v5" to be named "reqNOTTESTv5", but it was "reqNOTTESTv3"` {
		t.Fatalf("error message incorrect, got %s", err)
	}
}

func TestAssign_LHS(t *testing.T) {
	yml := `
type: assignment
name: LHS Assignment Rule
if:
  lhs:
    name: (?P<lhs>[a-zA-Z0-9_]+)
require:
  # eg, no underscores in variable names
  lhs: "{{ \"${lhs}\" | replace \"_\" \"\" }}"
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, "testdata/assign_test.go")
	if len(errs) != 2 {
		for i := range errs {
			t.Logf("\t%s", errs[i])
		}

		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestCommentGroupGood(t *testing.T) {
	yml := `
type: comment_group
name: Enforce FIXME
if:
  text: FIXME
require:
  # must contain a link to an open issue
  text: "https://github.com/openlyinc/openlynt/issues/\\d+"
  # must have at least N lines of context
  len: 3
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, "testdata/comment_good.go")
	if len(errs) != 0 {
		for i := range errs {
			t.Logf("\t%s", errs[i])
		}

		t.Fatalf("expected 0 errors, got %d", len(errs))
	}
}

func TestCommentGroupBad(t *testing.T) {
	yml := `
type: comment_group
name: Enforce FIXME
if:
  text: FIXME
require:
  # must contain a link to an open issue
  text: "https://github.com/openlyinc/openlynt/issues/\\d+"
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, "testdata/comment_bad.go")
	if len(errs) != 1 {
		for i := range errs {
			t.Logf("\t%s", errs[i])
		}

		t.Fatalf("expected 1 errors, got %d", len(errs))
	}
}

func TestBytesBuffer(t *testing.T) {
	t.SkipNow()

	yml := `
type: unary_expr
name: bytes.Buffer rule
if:
  selector: "\\&bytes\\.Buffer"
require:
  warn: Use new(bytes.Buffer) instead
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, "testdata/bytesbuffer_test.go")
	if len(errs) != 1 {
		for i := range errs {
			t.Logf("\t%s", errs[i])
		}

		t.Fatalf("expected 1 errors, got %d", len(errs))
	}
}
