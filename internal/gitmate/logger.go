package gitmate

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/fatih/color"
)

var defaultLogger atomic.Pointer[Logger]

func DefaultLogger() *Logger { return defaultLogger.Load() }

type Logger struct{}

func (l *Logger) Info(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgBlue)
	defer color.Unset()
	fmt.Println(output)
}

func (l *Logger) SystemInfo(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgMagenta)
	defer color.Unset()
	fmt.Println(output)
}

func (l *Logger) Warning(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgYellow)
	defer color.Unset()
	fmt.Println(output)
}

func (l *Logger) Success(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgGreen)
	defer color.Unset()
	fmt.Println(output)
}

func (l *Logger) Error(err error) {
	color.Set(color.FgRed)
	defer color.Unset()

	fmt.Fprint(os.Stderr, err.Error())
	fmt.Fprintln(os.Stderr, " Abbruch.")
}

func (l *Logger) ErrorString(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgRed)
	defer color.Unset()
	fmt.Println(output)
}
