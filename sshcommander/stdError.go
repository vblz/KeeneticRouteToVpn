package sshcommander

// StdErrError represents error while execution when exit code is 0, but stderr is not null
type StdErrError struct {
	StdErrorText string
}

func (err StdErrError) Error() string {
	return err.StdErrorText
}
