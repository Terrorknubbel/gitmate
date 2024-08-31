package cmd

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

const (
	Staging Action = iota
	Master
)

type Action int

func (c *Config) RunMergehelper(action Action) {
	c.logger.Success("Überprüfe Vorbedingungen")

	err := c.RunCommands(prerequisiteCommands())
	if err != nil {
		c.logger.Error(err)
		return
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		c.logger.Error(err)
		return
	}

	err = c.RunCommands(branchConditionCommands(currentBranch))
	if err != nil {
		c.logger.Error(err)
		return
	}

	switch action {
	case Staging:
		c.logger.Success("Bringe den Master und Staging Branch auf den aktuellsten Stand")
	case Master:
		c.logger.Success("Bringe den Master Branch auf den aktuellsten Stand")
	}

	err = c.RunCommands(branchRebaseCommands(action))
	if err != nil {
		c.logger.Error(err)
		return
	}

	err = c.RunCommands(branchMergeCommands(currentBranch, action))
	if err != nil {
		c.logger.Error(err)
		return
	}

	c.logger.Success("Merge erfolgreich.")

	handleFinalPush(c)
}

func prerequisiteCommands() []Command {
	return []Command{
		{Command: "git", Args: []string{"--version"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Git ist nicht installiert."},
		{Command: "git", Args: []string{"rev-parse", "--is-inside-work-tree"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Du befindest dich in keinem Git Verzeichnis."},
	}
}

func getCurrentBranch() (string, error) {
	output, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func branchConditionCommands(currentBranch string) []Command {
	return []Command{
		{Command: "git", Args: []string{"branch", "--show-current"}, Output: "Überprüfe Branch…", Expectation: "*", Forbidden: []string{"master", "main", "staging"}, ErrorMessage: "Du befindest dich in keinem Feature Branch."},
		{Command: "git", Args: []string{"status", "--porcelain"}, Output: "Prüfe auf Änderungen, die nicht zum Commit vorgesehen sind…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Änderungen, die nicht zum Commit vorgesehen sind. Bitte Committe oder Stashe diese vor einem Merge."},
		{Command: "git", Args: []string{"ls-remote", "origin", currentBranch}, Output: "Prüfe auf Remote Branch…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Der Branch " + currentBranch + " existiert nicht auf remote."},
		{Command: "git", Args: []string{"fetch"}, Output: "Aktualisiere Branch Informationen…", Expectation: "*", Forbidden: []string{}, ErrorMessage: ""},
		{Command: "git", Args: []string{"log", "origin/" + currentBranch + ".." + currentBranch}, Output: "Prüfe auf nicht gepushte Commits…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt nicht gepushte Änderungen. Bitte pushe diese vor einem Merge."},
		{Command: "git", Args: []string{"log", currentBranch + ".." + "origin/" + currentBranch}, Output: "Prüfe auf nicht gemergte Commits…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt nicht gemergte Änderungen. Bitte führe ein 'git pull --rebase' vor einem Merge aus."},
	}
}

func branchRebaseCommands(action Action) []Command {
	// TODO: Determine master vs main. (git ls-remote --symref origin HEAD | awk '/^ref:/ {print substr($2, 12)}')
	masterCommands := []Command{
		{Command: "git", Args: []string{"checkout", "master"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Master checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"pull", "--rebase"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Master pull fehlgeschlagen."},
	}

	stagingCommands := append(masterCommands, []Command{
		{Command: "git", Args: []string{"checkout", "staging"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"pull", "--rebase"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging pull fehlgeschlagen."},
	}...)

	if action == Staging {
		return stagingCommands
	}

	return masterCommands
}

func branchMergeCommands(featureBranch string, action Action) []Command {
	var finalBranch string
	if action == Staging {
		finalBranch = "staging"
	} else {
		finalBranch = "master"
	}

	return []Command{
		{Command: "git", Args: []string{"checkout", featureBranch}, Output: "Merge Master Branch in " + featureBranch + "…", Expectation: "*", Forbidden: []string{}, ErrorMessage: featureBranch + " checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"merge", "master"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Merge fehlgeschlagen."},
		{Command: "git", Args: []string{"status", "--porcelain"}, Output: "", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Merge Konflikte. Bitte behebe diese."},
		{Command: "git", Args: []string{"checkout", finalBranch}, Output: "Merge " + featureBranch + " in staging…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"merge", featureBranch}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Merge fehlgeschlagen."},
		{Command: "git", Args: []string{"status", "--porcelain"}, Output: "", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Merge Konflikte. Bitte behebe diese."},
	}
}

func finalPushCommands() []Command {
	return []Command{
		{Command: "git", Args: []string{"push"}, Output: "Pushe zum Remote Branch…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Push fehlgeschlagen."},
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
		err = c.RunCommands(finalPushCommands())

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
