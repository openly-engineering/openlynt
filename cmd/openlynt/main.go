package main

import (
	"flag"
	"io"
	stdlog "log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/golangci/revgrep"
	"github.com/openlyinc/openlynt/lint"
	yaml "gopkg.in/yaml.v3"
)

var (
	srcpath  string
	rulepath string

	revFilter bool

	log = stdlog.New(os.Stderr, "", 0)
)

func main() {
	flag.StringVar(&srcpath, "path", "", "path to parse .go files")
	flag.StringVar(&rulepath, "rules", ".openlynt.yml", "path to yaml config")
	flag.BoolVar(&revFilter, "revfilter", true, "enable revision filter; enabled by default")
	flag.Parse()

	if srcpath == "" {
		// see if argv0 is set
		if flag.Arg(0) == "" {
			log.Fatal("provide -path")
		}

		srcpath = flag.Arg(0)
	}

	linter := &lint.Linter{}
	fp, err := os.Open(rulepath)
	if err != nil {
		log.Fatalf("couldn't open config(%s): %s", rulepath, err)
	}

	if err := yaml.NewDecoder(fp).Decode(&linter); err != nil {
		log.Fatalf("couldn't decode yaml: %s", err)
	}

	var (
		patch    io.Reader
		newfiles []string
	)

	if linter.Revisions.From != "" && revFilter {
		patch, newfiles, err = revgrep.GitPatch(linter.Revisions.From, linter.Revisions.To)
		if err != nil {
			log.Fatalf("[openlynt] revision filter specified, but couldn't generate a git patch: %s", err)
		}
	}

	fail := false
	filepath.Walk(srcpath, func(fpath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return err
		}

		if path.Ext(fi.Name()) != ".go" {
			return err
		}

		for i := range linter.IgnorePaths {
			rxp := regexp.MustCompile(linter.IgnorePaths[i])

			if rxp.MatchString(fpath) {
				return nil
			}
		}

		set := &lint.Violations{}
		for i := range linter.Rules {
			rule := linter.Rules[i]

			errs := lint.Walk(rule, fpath)
			for j := range errs {
				if le, ok := errs[j].(*lint.Violation); ok {
					set.Violations = append(set.Violations, le)
				} else if les, ok := errs[j].(*lint.Violations); ok {
					for k := range les.Violations {
						set.Violations = append(set.Violations, les.Violations[k])
					}
				} else {
					// error in implementation of lint rule
					log.Printf("[openlynt] %s error: %s", rule.Name, errs[j])
				}
			}
		}

		if patch != nil {
			cnt := len(set.Violations)

			set, err = lint.FilterByRevision(patch, newfiles, set)
			if err != nil {
				log.Fatalf("[openlynt] failed to filter by revision: %s", err)
			}

			if cnt > len(set.Violations) {
				log.Printf("[openlynt] filtered %d violations in %s due to revision limits", cnt-len(set.Violations), fpath)
			}
		}

		for i := range set.Violations {
			fail = true
			v := set.Violations[i]

			logViolation(v.Rule, v, v.File)
		}

		return nil
	})

	if fail {
		os.Exit(1)
	}
}

func logViolation(r *lint.Rule, v *lint.Violation, path string) {
	log.Printf(
		"%s:%d violation of %s: %s",
		path, v.Position.Line,
		r.Name,
		v.Error(),
	)
}
