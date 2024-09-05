package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Terrorknubbel/gitmate/internal/gitmate"
	"github.com/spf13/cobra"
)

type Config struct {
	logger *gitmate.Logger

	commandDirAbsPath string
}

func NewConfig() (*Config, error) {
	c := &Config{
		logger: gitmate.DefaultLogger(),
	}

	wd, err := os.Getwd()
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	c.commandDirAbsPath = wd

	return c, nil
}

// newRootCmd returns a new root github.com/spf13/cobra.Command.
func (c *Config) newRootCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:           "gitmate",
		Short:         "Automagische Git Befehle",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	for _, cmd := range []*cobra.Command{
		c.newMenuViewPreviewCmd(),
		c.newMergehelperCmd(),
		c.newInfrastructureCmd(),
	} {
		if cmd != nil {
			rootCmd.AddCommand(cmd)
		}
	}

	return rootCmd, nil
}

// execute creates a new root command and executes it with args.
func (c *Config) execute(args []string) error {
	rootCmd, err := c.newRootCmd()
	if err != nil {
		return err
	}
	rootCmd.SetArgs(args)

	return rootCmd.Execute()
}

func (c *Config) run(dir gitmate.AbsPath, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	if !dir.Empty() {
		cmd.Dir = string(dir)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	outputString := strings.TrimSpace(string(output[:]))
	return outputString, nil
}
