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

	log = stdlog.New(os.Stderr, "", stdlog.Ltime)
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

		for i := range linter.Rules {
			rule := linter.Rules[i]

			errs := lint.Walk(rule, fpath)
			for j := range errs {
				fail = true

				if le, ok := errs[j].(*lint.Error); ok {
					// violation of a lint rule
					log.Printf("%s violation in %s:%d: %s",
						rule.Name, fpath, le.Position.Line, errs[j])
				} else if les, ok := errs[j].(*lint.ErrorCollection); ok {
					for k := range les.Errors {
						log.Printf("%s violation in %s:%d: %s",
							rule.Name, fpath, les.Errors[k].Position.Line, les.Errors[k])
					}
				} else {
					// error in implementation of lint rule
					log.Printf("%s error: %s", rule.Name, errs[j])
				}
			}
		}

		return nil
	})

	if fail {
		os.Exit(1)
	}
}
