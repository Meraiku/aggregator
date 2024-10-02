package cli

import "errors"

var (
	ErrNoArgs                  = errors.New("missing arguments")
	ErrCommandAlreadyRegisterd = errors.New("command already registerd")
	ErrUnknownState            = errors.New("unknown state provided")
	ErrUnknownCommand          = errors.New("unknown command provided")
	ErrInvalidArgumentsCount   = errors.New("invalid arguments number")
)
