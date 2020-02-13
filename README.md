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
rules:

  rule_pkgprefix:
    type: import
    name: Named Import Rule
    if:
      # any import path that contains "pkg/[a-z0-9]+/v\d+", ie "pkg/something/v1"
      path: regexp: \/pkg\/(?P<prefix>[a-z0-9]+)\/v(?P<version>[0-9]+)

    require:
      # require the named import to match the following, eg "pkgSOMETHINGv1"
      name: pkg{{ "${prefix}" | upper }}v${version}
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

## Supported Rules

At the moment, the following rule types are supported.

### `import`

```yaml
rule_short_name:
  type: import
  name: Friendly Rule Name

  if:
    path: REGEXP

  require:
    name: "custom-import-name-template"
```

The `import` rule supports `text/template` replacements with your regexp, as
well as [sprig](https://github.com/Masterminds/sprig) functions - eg, piping
a named-regex-match into `upper`, `lower`, `plural`, etc.

### `comment_group`

```yaml
rule_short_name:
  type: comment_group
  name: Friendly Rule Name

  if:
    text: REGEXP

  require:
    # optional
    text: REGEXP

    # optional
    len: int
```

The `comment_group` rule supports regexp matches and requirements for the
entire comment group text. It will match both multiple lines of `//` in
succession as well as `/*`.

The `len` parameter specifies an optional required minimum-line-count for a
matching comment. As an example, if you require a `TODO` comment to contain a
link to an issue tracker _and_ have at least 2 lines of information about the
problem:

```
rule_todo:
  type: comment_group
  name: TODO Requirement
  if:
    text: TODO

  require:
    text: "https?://github.com/your/project/issues/\\d+"
    len: 2
```

A comment like this would fail to pass:

```go
// TODO: I'm not saying anything useful or linking to the issue.
```
