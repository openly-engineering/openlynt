package lint

type Linter struct {
	IgnorePaths []string `yaml:"ignore_paths"`

	Revisions struct {
		From string
		To   string
	}

	Rules map[string]*Rule
}
