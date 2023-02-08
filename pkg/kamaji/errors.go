package kamaji

import (
	"github.com/pkg/errors"
)

var (
	ErrKubeClientNil          = errors.New("client kubernetes is nil")
	ErrLoggerNil              = errors.New("logger is nil")
)
