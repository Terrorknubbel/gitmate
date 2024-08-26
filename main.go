package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/Terrorknubbel/gitmate/gitrunner"
	"github.com/fatih/color"
	fzf "github.com/junegunn/fzf/src"
)

func main() {
	inputChan := make(chan string)
	outputChan := make(chan string)
	var wg sync.WaitGroup

	go func() {
		defer close(inputChan)
		for _, s := range []string{"Staging", "Master"} {
			inputChan <- s
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for selectedOption := range outputChan {
			switch selectedOption {
			case "Staging":
				gitrunner.Run(gitrunner.Staging)
			case "Master":
				gitrunner.Run(gitrunner.Master)
			default:
				fmt.Println("Diese Option wird aktuell nicht unterstützt.")
			}
		}
	}()

	options, err := fzf.ParseOptions(
		true,
		[]string{
			"--height=50%",
			"--ansi",
			"--reverse",
			"--pointer=👉",
			"--cycle",
			"--header=Wähle den Ziel-Branch aus:",
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(fzf.ExitError)
	}

	options.Input = inputChan
	options.Output = outputChan

	fmt.Print("Merge Automatisierung mit ")
	color.Set(color.FgYellow)
	fmt.Print("GitMate 🪄")
	color.Unset()

	code, err := fzf.Run(options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(code)
	}

	close(outputChan)
	wg.Wait()
}
