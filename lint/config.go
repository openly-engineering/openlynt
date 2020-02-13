package lint

type Linter struct {
	IgnorePaths []string `yaml:"ignore_paths"`

	Rules map[string]*Rule
}
