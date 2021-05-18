package clients

// AcErrorType
type AcErrorType int

const (
	ConnectionError       AcErrorType = iota
	InstanceCreationError AcErrorType = iota
	TokenGenerationError  AcErrorType = iota
)

// String
func (e AcErrorType) String() string {
	return [...]string{
		"ConnectionError",
		"InstanceCreationError",
		"TokenGenerationError",
	}[e]
}

// AcError
type AcError struct {
	Type AcErrorType
	Err  error
}
