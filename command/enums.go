package command

//go:generate stringer -type ExitCode -output enums_string.go
type ExitCode int

const (
	ExitCodeOK ExitCode = iota
	ExitCodeError
	ExitCodeParseFlagsError
	ExitCodeBadArgs
	ExitCodeInvalidURL
	ExitCodeProjectNotFound
	ExitCodePageNotFound
	ExitCodeFetchFailure
)
