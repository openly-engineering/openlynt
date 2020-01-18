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

func TestAssign(t *testing.T) {
	yml := `
type: assignment
name: Assignment Rule
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
	if len(errs) != 1 {
		t.Fatalf("expected 1 errors, got %d", len(errs))
	}
}
