package comment

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestBadCommentTODO(t *testing.T) {
	td := analysistest.TestData()

	analysistest.Run(t, td, Analyzer, "bad")
}
