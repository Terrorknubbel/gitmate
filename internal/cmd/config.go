package cmd

import (
	"github.com/Terrorknubbel/gitmate/internal/gitmate"
	"github.com/spf13/cobra"
)

type Config struct {
	logger *gitmate.Logger
}

func NewConfig() *Config {
	return &Config{
		logger: gitmate.DefaultLogger(),
	}
}

// newRootCmd returns a new root github.com/spf13/cobra.Command.
func (c *Config) newRootCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:           "gitmate",
		Short:         "Run your daily git commands automagically.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	for _, cmd := range []*cobra.Command{
		c.newMenuViewPreviewCmd(),
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
