package commandrunner

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Terrorknubbel/gitmate/printer"
)

type CommandCheck struct {
	Command      string
	Args         []string
	Output       string
	Expectation  string
	Forbidden    []string
	ErrorMessage string
}

func RunCommands(commands []CommandCheck) error {
	for _, command := range commands {
		err := run(command)
		if err != nil {
			return err
		}
	}

	return nil
}

func run(command CommandCheck) error {
	printer.Info(command.Output)

	output, err := exec.Command(command.Command, command.Args...).Output()
	outputString := strings.TrimSpace(string(output[:]))

	if err != nil {
		fmt.Println(err.Error())
		return errors.New(command.ErrorMessage)
	}

	if command.Expectation != "*" && command.Expectation != outputString {
		return errors.New(command.ErrorMessage)
	}

	for _, v := range command.Forbidden {
		if outputString == v {
			return errors.New(command.ErrorMessage)
		}
	}

	return nil
}
