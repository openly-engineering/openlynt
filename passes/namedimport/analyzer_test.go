package namedimport

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestBadUnderscoredImport(t *testing.T) {
	td := analysistest.TestData()

	analysistest.Run(t, td, Analyzer, "bad")
}
