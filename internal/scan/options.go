package scan

type Options struct {
	Format     string
	Output     string
	FailOn     string
	ConfigPath string
	NoColor    bool
	Verbose    bool
	Quiet      bool
	Include    []string
	Exclude    []string
	Frameworks []string
	Rules      []string
}
