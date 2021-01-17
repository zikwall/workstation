package workstation

import "errors"

var (
	ErrorAsyncProcessAlreadyExists              = errors.New("Error asynchronous process already exists")
	ErrorAsyncProcessNotFoundOrAlreadyCompleted = errors.New("The process was not found or has already been completed")
	ErrorShutdownWithoutGracefulCompletion      = errors.New("Shutdown completed without graceful completion")
	ErrorAsyncProcessChannelAlreadyClosed       = errors.New("The process channel has already been closed")
)
