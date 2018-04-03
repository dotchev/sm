package postgres

import (
	"fmt"
)

type Logger struct {
}

func (Logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (Logger) Verbose() bool {
	return true
}
