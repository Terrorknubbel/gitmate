package cmd

import (
	"errors"
)

type GitCommand struct {
	Args         []string
	Output       string
	Expectation  string
	Forbidden    []string
	ErrorMessage string
}

type GitCommands []GitCommand

func (c *Config) RunGitCommands(commands GitCommands) error {
	for _, command := range commands {
		err := c.runCommand(command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) runCommand(command GitCommand) error {
	c.logger.Info(command.Output)

	output, err := c.run("git", command.Args...)
	if err != nil {
		return errors.New(command.ErrorMessage)
	}

	if command.Expectation != "*" && command.Expectation != output {
		return errors.New(command.ErrorMessage)
	}

	for _, v := range command.Forbidden {
		if output == v {
			return errors.New(command.ErrorMessage)
		}
	}

	return nil
}
