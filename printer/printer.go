package printer

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Info(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgBlue)
	defer color.Unset()
	fmt.Println(output)
}

func Warning(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgYellow)
	defer color.Unset()
	fmt.Println(output)
}

func Success(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgGreen)
	defer color.Unset()
	fmt.Println(output)
}

func Error(err error) {
	color.Set(color.FgRed)
	defer color.Unset()

	fmt.Fprint(os.Stderr, err.Error())
	fmt.Fprintln(os.Stderr, " Abbruch.")
}

func ErrorString(output string) {
	if output == "" {
		return
	}

	color.Set(color.FgRed)
	defer color.Unset()
	fmt.Println(output)
}
