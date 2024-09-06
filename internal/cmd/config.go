package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Terrorknubbel/gitmate/internal/gitmate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	logger *gitmate.Logger

	infrastructureAbsPath gitmate.AbsPath
	commandDirAbsPath gitmate.AbsPath
}

func NewConfig() (*Config, error) {
	c := &Config{
		logger: gitmate.DefaultLogger(),
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/gitmate")
	viper.SetEnvPrefix("gitmate")
	viper.BindEnv("infrastructure_path")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			c.infrastructureAbsPath = gitmate.AbsPath(viper.GetString("INFRASTRUCTURE_PATH"))
		} else {
			return c, err
		}
	} else {
		c.infrastructureAbsPath = gitmate.AbsPath(viper.GetString("infrastructure_path"))
	}

	if c.infrastructureAbsPath.Empty() {
		c.infrastructureAbsPath = gitmate.AbsPath(viper.GetString("INFRASTRUCTURE_PATH"))
	}

	c.commandDirAbsPath = gitmate.AbsPath(wd)

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

func (c *Config) run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	cmd.Dir = c.commandDirAbsPath.String()

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	outputString := strings.TrimSpace(string(output[:]))

	return outputString, nil
}

// withinDir changes the current working directory to the given directory and
// executes the given function. It restores the current working directory after
// the function has finished.
func (c *Config) withinDir(dir gitmate.AbsPath, f func(c *Config) error) error {
	oldDir := c.commandDirAbsPath
	defer func() {
		c.commandDirAbsPath = oldDir
		c.logger.SystemInfo("Wechsle ins Verzeichnis " + oldDir.String())
	}()

	c.logger.SystemInfo("Wechsle ins Verzeichnis " + dir.String())
	c.commandDirAbsPath = dir

	return f(c)
}
