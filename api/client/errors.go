package client

import "errors"

// ErrConnection
var ErrConnection = errors.New("failed to connect to the buttress server")

// ErrAppInstance
var ErrAppInstance = errors.New("failed to instantiate the buttress application")

// ErrAuthInstance
var ErrAuthInstance = errors.New("failed to instantiate auth")
