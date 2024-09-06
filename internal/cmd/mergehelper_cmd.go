package cmd

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	Staging Action = "staging"
	Master  Action = "master"
)

type Action string

func (c *Config) newMergehelperCmd() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:   "merge",
		Short: "Automatisierter Merge vom Feature Branch in Staging oder Master",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || args[0] != string(Staging) && args[0] != string(Master) {
				//lint:ignore ST1005 Capitalization is fine…
				return errors.New("Ungültiges Argument. Bitte 'staging' oder 'master' angeben.")
			}
			return nil
		},
		RunE: c.RunMergehelper,
	}

	return gitCmd
}

func (c *Config) RunMergehelper(cmd *cobra.Command, args []string) error {
	c.logger.Success("Überprüfe Vorbedingungen")

	err := c.RunGitCommands(prerequisiteCommands())
	if err != nil {
		return err
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	err = c.RunGitCommands(branchConditionCommands(currentBranch))
	if err != nil {
		return err
	}

	action := Action(args[0])

	switch action {
	case Staging:
		c.logger.Success("Bringe den Master und Staging Branch auf den aktuellsten Stand")
	case Master:
		c.logger.Success("Bringe den Master Branch auf den aktuellsten Stand")
	}

	err = c.RunGitCommands(branchRebaseCommands(action))
	if err != nil {
		return err
	}

	err = c.RunGitCommands(branchMergeCommands(currentBranch, action))
	if err != nil {
		return err
	}

	c.logger.Success("Merge erfolgreich.")

	handleFinalPush(c)

	return nil
}

func prerequisiteCommands() GitCommands {
	return GitCommands{
		{Args: []string{"--version"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Git ist nicht installiert."},
		{Args: []string{"rev-parse", "--is-inside-work-tree"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Du befindest dich in keinem Git Verzeichnis."},
	}
}

func getCurrentBranch() (string, error) {
	output, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func branchConditionCommands(currentBranch string) GitCommands {
	return GitCommands{
		{Args: []string{"branch", "--show-current"}, Output: "Überprüfe Branch…", Expectation: "*", Forbidden: []string{"master", "main", "staging"}, ErrorMessage: "Du befindest dich in keinem Feature Branch."},
		{Args: []string{"status", "--porcelain"}, Output: "Prüfe auf Änderungen, die nicht zum Commit vorgesehen sind…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Änderungen, die nicht zum Commit vorgesehen sind. Bitte Committe oder Stashe diese vor einem Merge."},
		{Args: []string{"ls-remote", "origin", currentBranch}, Output: "Prüfe auf Remote Branch…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Der Branch " + currentBranch + " existiert nicht auf remote."},
		{Args: []string{"fetch"}, Output: "Aktualisiere Branch Informationen…", Expectation: "*", Forbidden: []string{}, ErrorMessage: ""},
		{Args: []string{"log", "origin/" + currentBranch + ".." + currentBranch}, Output: "Prüfe auf nicht gepushte Commits…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt nicht gepushte Änderungen. Bitte pushe diese vor einem Merge."},
		{Args: []string{"log", currentBranch + ".." + "origin/" + currentBranch}, Output: "Prüfe auf nicht gemergte Commits…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt nicht gemergte Änderungen. Bitte führe ein 'git pull --rebase' vor einem Merge aus."},
	}
}

func branchRebaseCommands(action Action) GitCommands {
	// TODO: Determine master vs main. (git ls-remote --symref origin HEAD | awk '/^ref:/ {print substr($2, 12)}')
	masterCommands := GitCommands{
		{Args: []string{"checkout", "master"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Master checkout fehlgeschlagen."},
		{Args: []string{"pull", "--rebase"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Master pull fehlgeschlagen."},
	}

	stagingCommands := append(masterCommands, GitCommands{
		{Args: []string{"checkout", "staging"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging checkout fehlgeschlagen."},
		{Args: []string{"pull", "--rebase"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging pull fehlgeschlagen."},
	}...)

	if action == Staging {
		return stagingCommands
	}

	return masterCommands
}

func branchMergeCommands(featureBranch string, action Action) GitCommands {
	return GitCommands{
		{Args: []string{"checkout", featureBranch}, Output: "Merge Master Branch in " + featureBranch + "…", Expectation: "*", Forbidden: []string{}, ErrorMessage: featureBranch + " checkout fehlgeschlagen."},
		{Args: []string{"merge", "master"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Merge fehlgeschlagen."},
		{Args: []string{"status", "--porcelain"}, Output: "", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Merge Konflikte. Bitte behebe diese."},
		{Args: []string{"checkout", string(action)}, Output: "Merge " + featureBranch + " in " + string(action) + "…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging checkout fehlgeschlagen."},
		{Args: []string{"merge", featureBranch}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Merge fehlgeschlagen."},
		{Args: []string{"status", "--porcelain"}, Output: "", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Merge Konflikte. Bitte behebe diese."},
	}
}

func finalPushCommands() GitCommands {
	return GitCommands{
		{Args: []string{"push"}, Output: "Pushe zum Remote Branch…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Push fehlgeschlagen."},
	}
}

func handleFinalPush(c *Config) {
	reader := bufio.NewReader(os.Stdin)
	c.logger.Warning("Merge zum Remote Branch pushen? (y/n)")

	input, err := reader.ReadString('\n')
	if err != nil {
		c.logger.Error(err)
		return
	}

	// Clean the input
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		err = c.RunGitCommands(finalPushCommands())

		if err != nil {
			c.logger.Error(err)
		}
	} else if input == "n" || input == "no" {
		c.logger.ErrorString("Keine Änderungen gepusht.")
	} else {
		c.logger.ErrorString("Ungültige Eingabe. Bitte gebe 'y' oder 'n' ein.")
		handleFinalPush(c)
	}
}
