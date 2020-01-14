package main

import (
	"bytes"
	"flag"
	stdlog "log"
	"os"
	"path"
	"path/filepath"

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

	var rules map[string]*lint.Rule
	fp, err := os.Open(rulepath)
	if err != nil {
		log.Fatalf("couldn't open config(%s): %s", rulepath, err)
	}

	if err := yaml.NewDecoder(fp).Decode(&rules); err != nil {
		log.Fatalf("couldn't decode yaml: %s", err)
	}

	buf := new(bytes.Buffer)
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

		buf.Reset()
		fp, err := os.Open(fpath)
		if err != nil {
			return err
		}

		_, err = buf.ReadFrom(fp)
		if err != nil {
			return err
		}

		src := buf.String()
		for i := range rules {
			errs := lint.Walk(rules[i], src)
			for i := range errs {
				fail = true

				log.Printf("%s: %s\n", fpath, errs[i])
			}
		}

		return nil
	})

	if fail {
		os.Exit(1)
	}
}
