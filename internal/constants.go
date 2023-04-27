package internal

import (
	"regexp"
)

var (
	queue = regexp.MustCompile(`^\/queue\/(\S+)$`)
)

const checkerIntervalSec = 1

const (
	ReconnectTimerSec      = 30
	ReconnectAttemptsCount = 5
)
