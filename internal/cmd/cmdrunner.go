package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Command struct {
	Command      string
	Args         []string
	Output       string
	Expectation  string
	Forbidden    []string
	ErrorMessage string
}

func (c *Config) RunCommands(commands []Command) error {
	for _, command := range commands {
		err := c.runCommand(command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) runCommand(command Command) error {
	c.logger.Info(command.Output)

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
