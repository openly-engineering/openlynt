# openlynt

`openlynt` is a lint tool for Go that works on individual ast nodes.
This tool is a work-in-progress - we'll add more as we encode our rules.
More documentation will follow as well once things are solidified.

Supported rules:

- required formatting of named import declarations

## Install

    go get -u github.com/openlyinc/openlynt/cmd/openlynt

## Usage

A sample `openlynt.yml` file can be found at `cmd/openlynt/testdata/openlynt.yml`.

You can provide a path to the `openlynt.yml` file via `-rules` or place it in
the current directory as `.openlynt.yml`.

To provide a path to the source files to parse, you can either pass `-path` OR
run the following: `openlynt path/`

## Example

For this example, we have package `main` that must import several other
packages that fulfill an interface. For clarity and sanity, we require these
packages to be named imports - and to be named following a certain pattern.
In this case, we want them named `pkg<PREFIX>v<version>`.


```go
package main

import (
	"fmt"

	"xyz.org/pkg/prefix/v503"		// unnamed
	INCORRECTv99 "xyz.org/pkg/prefix2/v99"	// wrong name
	pkgOKv1 "xyz.org/pkg/ok/v1"		// correct name
)

func main() {
	fmt.Println("...")
}
```


Our ruleset (`.openlynt.yml`) would be:

```yaml
rule_pkgprefix:
  type: import
  require:
    # any import path that contains "pkg/[a-z0-9]+/v\d+", ie "pkg/something/v1"
    regexp: \/pkg\/(?P<prefix>[a-z0-9]+)\/v(?P<version>[0-9]+)
    name:
      # require the named import to match the following, eg "pkgSOMETHINGv1"
      template: pkg{{ "${prefix}" | upper }}v${version}
```


Running `openlynt` against this source file would result in:

```
$ ./openlynt -path testdata -rules testdata/openlynt.yml && echo "ok"
22:25:51 testdata/incorrect.go: expected import("xyz.org/pkg/prefix/v503") to be named(pkgPREFIXv503), but it was()
22:25:51 testdata/incorrect.go: expected import("xyz.org/pkg/prefix2/v99") to be named(pkgPREFIX2v99), but it was(INCORRECTv99)
```

After fixing:

```
$ ./openlynt -path testdata/correct.go -rules testdata/openlynt.yml && echo "ok"
ok
```

