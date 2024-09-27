package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func (c *Config) newInfrastructureCmd() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:     "infrastructure",
		Aliases: []string{"i"},
		Short:   "Bringt das Infrastructure Repo auf den aktuellsten Stand und stellt den korrekten Zustand sicher",
		RunE:    c.infrastructure,
	}

	return gitCmd
}

func (c *Config) infrastructure(cmd *cobra.Command, args []string) error {
	if c.infrastructureAbsPath.Empty() {
		//lint:ignore ST1005 Capitalization is fine…
		return errors.New("Kein Infrastructure Verzeichnis angegeben. Setze die Umgebungsvariable GITMATE_INFRASTRUCTURE_PATH oder setze den Wert für 'infrastructure_path' in '~/.config/gitmate/config.json'.")
	}

	if !c.infrastructureAbsPath.IsAbs() {
		return errors.New("Das Infrastructure Verzeichnis muss absolut sein. Aktuell gesetztes Verzeichnis '" + c.infrastructureAbsPath.String() + "'.")
	}

	err := c.withinDir(c.infrastructureAbsPath, func(c *Config) error {
		err := c.RunGitCommands(infrastructureCommands())
		if err == nil {
			c.logger.Success("Infrastructure ist nun auf dem aktuellsten Stand.")
		}
		return err
	})

	return err
}

func infrastructureCommands() GitCommands {
	return GitCommands{
		{Args: []string{"--version"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Git ist nicht installiert."},
		{Args: []string{"rev-parse", "--is-inside-work-tree"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Du befindest dich in keinem Git Verzeichnis."},
		{Args: []string{"branch", "--show-current"}, Output: "Überprüfe Branch…", Expectation: "master", Forbidden: []string{}, ErrorMessage: "Du befindest dich nicht im master Branch."},
		{Args: []string{"status", "--porcelain"}, Output: "Prüfe auf Änderungen, die nicht zum Commit vorgesehen sind…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Änderungen, die nicht zum Commit vorgesehen sind. Bitte Committe oder Stashe diese vor einem Merge."},
		{Args: []string{"ls-remote", "origin", "master"}, Output: "Prüfe auf Remote Branch…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Der master Branch existiert nicht auf remote."},
		{Args: []string{"fetch"}, Output: "Aktualisiere Branch Informationen…", Expectation: "*", Forbidden: []string{}, ErrorMessage: ""},
		{Args: []string{"log", "origin/master..master"}, Output: "Prüfe auf nicht gepushte Commits…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt nicht gepushte Änderungen. Bitte pushe diese vor einem Merge."},
		{Args: []string{"pull", "--rebase"}, Output: "Bringe den master Branch auf den aktuellsten Stand", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Master pull fehlgeschlagen."},
	}
}
