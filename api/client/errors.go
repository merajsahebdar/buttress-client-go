package client

// ClientErrorType
type ClientErrorType int

const (
	ConnectionError       ClientErrorType = iota
	InstanceCreationError ClientErrorType = iota
	TokenGenerationError  ClientErrorType = iota
)

// String
func (e ClientErrorType) String() string {
	return [...]string{
		"ConnectionError",
		"InstanceCreationError",
		"TokenGenerationError",
	}[e]
}

// ClientError
type ClientError struct {
	Type ClientErrorType
	Err  error
}
