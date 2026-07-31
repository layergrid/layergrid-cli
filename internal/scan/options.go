package scan

type Options struct {
	Format     string
	Output     string
	FailOn     string
	NoColor    bool
	Verbose    bool
	Quiet      bool
	Frameworks []string
	Rules      []string
}
