package cmd

// Build command builds the program at the specified path.
type Build struct{}

// Execute runs the command.
func (b Build) Execute(_ []string) error { return nil }
