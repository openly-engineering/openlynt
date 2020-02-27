package main

import (
	"flag"
	stdlog "log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/openlyinc/openlynt/lint"
	yaml "gopkg.in/yaml.v3"
)

var (
	srcpath  string
	rulepath string

	log = stdlog.New(os.Stderr, "", 0)
)

func main() {
	flag.StringVar(&srcpath, "path", "", "path to parse .go files")
	flag.StringVar(&rulepath, "rules", ".openlynt.yml", "path to yaml config")
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
				fail = true

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

		if linter.Revisions.From != "" {
			set, err = lint.FilterByRevision(set, linter.Revisions.From, linter.Revisions.To)
			if err != nil {
				log.Fatalf("failed to filter by revision: %s", err)
			}
		}

		for i := range set.Violations {
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
