package cmd

import (
	"github.com/Terrorknubbel/gitmate/internal/gitmate"
	"github.com/spf13/cobra"
)

func (c *Config) newInfrastructureCmd() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:     "infrastructure",
		Aliases: []string{"i"},
		Short:   "Bringt das Infrastruktur Repo auf den aktuellsten Stand und stellt den korrekten Zustand sicher",
		RunE:    c.infrastructure,
	}

	return gitCmd
}

func (c *Config) infrastructure(cmd *cobra.Command, args []string) error {
	infrastructureAbsPath := "/home/terrorknubbel/Projekte/go/chezmoi"

	// output, err := c.run(gitmate.AbsPath(infrastructureAbsPath), "pwd")
	c.run(gitmate.AbsPath(infrastructureAbsPath), "git", "checkout", "-b", "tmp")
	output, err := c.run(gitmate.AbsPath(infrastructureAbsPath), "git", "branch", "--show-current")

	if err != nil {
		return err
	}

	c.logger.Info(output)

	return nil
}
