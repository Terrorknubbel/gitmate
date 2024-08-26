package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/Terrorknubbel/gitmate/gitrunner"
	"github.com/fatih/color"
	fzf "github.com/junegunn/fzf/src"
)

func handleFzfPreview() {
	if len(os.Args) > 1 && os.Args[1] == "--preview" {
		option := os.Args[2]

		var updateLine string

		if option == "Staging" {
			updateLine = "Aktualisiert den Feature, Staging & Master Branch"
		} else {
			updateLine = "Aktualisiert den Feature und Master Branch"
		}

		output :=
`Bringt den aktuellen Feature Branch in den %s Branch.

Folgende Schritte werden durchgeführt:
 ↳ Überprüfung auf sauberen Branch Status (alle Commits gepusht, remote Status, …)
 ↳ %s
 ↳ Merged Master in den Feature Branch
 ↳ Merged den Feature Branch in den %s Branch

Im Falle von Merge Konflikten oder anderen Fehlern wird sofortig abgebrochen.
Es wird NICHT automatisch gepusht. Dies passiert erst nach einer manuellen Bestätigung.
`
		s := fmt.Sprintf(output, option, updateLine, option)
		fmt.Print(s)

		os.Exit(0)
	}
}

func main() {
	handleFzfPreview()

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
			"--height=95%",
			"--ansi",
			"--reverse",
			"--pointer=👉",
			"--cycle",
			"--header=Wähle den Ziel-Branch aus:",
			"--preview=gitmate --preview {}",
			"--preview-window=wrap",
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
