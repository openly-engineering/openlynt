package lint

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestErr(t *testing.T) {
	yml := `
type: import
require:
  x: abcd
  # eg, /pkg/(1: prefix)/v(2: 12)\
  regexp: \/pkg\/(?P<prefix>[a-z]+)\/v(?P<version>[0-9]+)

  # a property of the "ast.ImportSpec" struct.
  name:
    template: req{{ "${prefix}" | upper }}v${version}
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	src := `
package main

import (
	reqLULv1 "xyz.org/pkg/test/v1"
)

`

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, src)
	if len(errs) == 0 {
		t.Fatalf("expected >0 errors")
	}
}

func TestBasicOK(t *testing.T) {
	yml := `
type: import
require:
  regexp: \/pkg\/(?P<prefix>[a-z]+)\/v(?P<version>[0-9]+)
  name:
    template: req{{ "${prefix}" | upper }}v${version}
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	src := `
package main

import (
	reqTESTv1 "xyz.org/pkg/test/v1"

	reqTESTv2 "xyz.org/pkg/test/v2"
)

`

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, src)
	if len(errs) != 0 {
		t.Fatalf("expected 0 errors, got %d", len(errs))
	}
}

func TestBasicMixed(t *testing.T) {
	yml := `
type: import
require:
  regexp: \/pkg\/(?P<prefix>[a-z]+)\/v(?P<version>[0-9]+)
  name:
    template: req{{ "${prefix}" | upper }}v${version}
`

	r := &Rule{}
	if err := yaml.Unmarshal([]byte(yml), r); err != nil {
		t.Fatalf("error while unmarshalling: %s", err)
	}

	src := `
package main

import (
	reqTESTv1 "xyz.org/pkg/test/v1"
	reqNOTTESTv3 "xyz.org/pkg/nottest/v5"
	reqTESTv2 "xyz.org/pkg/test/v2"
)

`

	defer func() {
		if x := recover(); x != nil {
			t.Logf("recovered from panic: %v", x)
		}
	}()

	errs := Walk(r, src)
	if len(errs) != 1 {
		t.Fatalf("expected 1 errors, got %d", len(errs))
	}

	err := errs[0]
	if err.Error() != "expected import(\"xyz.org/pkg/nottest/v5\") to be named(reqNOTTESTv5), but it was(reqNOTTESTv3)" {
		t.Fatalf("error message incorrect, got %s", err)
	}
}
