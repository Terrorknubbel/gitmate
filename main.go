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

	headerText := "Wähle den Ziel-Branch aus:"

	cmd := exec.Command("fzf", "--height", "50%", "--ansi", "--reverse", "--pointer", "👉", "--cycle", "--header", headerText)
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))

	fmt.Println("Merge Automatisierung mit GitMate 🪄")

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	selectedOption := strings.TrimSpace(string(out))

	switch selectedOption {
	case "Staging":
		gitrunner.Staging()
	default:
		fmt.Println("Diese Option wird aktuell nicht unterstützt.")
	}
}
