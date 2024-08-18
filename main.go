package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"github.com/Terrorknubbel/gitmate/gitrunner"
)

func main() {
	options := []string{
		"Staging",
		"Master",
	}

	headerText := "Select an option to run the corresponding script:"

	cmd := exec.Command("fzf", "--height", "50%", "--ansi", "--reverse", "--pointer", "👉", "--cycle", "--header", headerText)
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	selectedOption := strings.TrimSpace(string(out))

	switch selectedOption {
	case "Staging":
		gitrunner.Staging()
	default:
		fmt.Println("No valid option selected.")
	}
}
