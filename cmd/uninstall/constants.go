package uninstall

import (
	"time"
)

const (
	commandName             = "uninstall"
	commandShortDescription = "uninstall the Kamaji Kubernetes operator"
)

var (
	timeout = time.Minute
)
