package tcp

import (
	"github.com/pkg/errors"
)

const (
	defaultKubernetesVersion = "v1.26.1"
)

var (
	ErrKubernetesVersionEmpty = errors.New("Kubernetes version cannot be empty")
	ErrNamespaceEmpty         = errors.New("control plane namespace cannot be empty")
	ErrNameEmpty              = errors.New("control plane name cannot be empty")
	ErrServiceTypeEmpty       = errors.New("control plane service type cannot be empty")
	ErrKubeClientNil          = errors.New("client kubernetes is nil")
	ErrLoggerNil              = errors.New("logger is nil")
)
