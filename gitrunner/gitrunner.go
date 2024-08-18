package gitrunner

import (
	"os/exec"
	"strings"

	"github.com/Terrorknubbel/gitmate/commandrunner"
	"github.com/Terrorknubbel/gitmate/printer"
)

func Staging() {
	printer.Success("Überprüfe Vorbedingungen")

	err := commandrunner.RunCommands(prerequisiteCommands())

	if err != nil {
		printer.Error(err)
		return
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		printer.Error(err)
	}

	err = commandrunner.RunCommands(branchConditionCommands(currentBranch))

	if err != nil {
		printer.Error(err)
		return
	}

	printer.Success("Bringe den Master und Staging Branch auf den aktuellsten Stand")

	err = commandrunner.RunCommands(branchRebaseCommands())

	if err != nil {
		printer.Error(err)
		return
	}

	err = commandrunner.RunCommands(branchMergeCommands(currentBranch))

	if err != nil {
		printer.Error(err)
		return
	}

	printer.Success("Merge erfolgreich.")
}

func prerequisiteCommands() []commandrunner.CommandCheck {
	return []commandrunner.CommandCheck{
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

func branchConditionCommands(currentBranch string) []commandrunner.CommandCheck {
	return []commandrunner.CommandCheck{
		{Command: "git", Args: []string{"branch", "--show-current"}, Output: "Überprüfe Branch…", Expectation: "*", Forbidden: []string{"master", "main", "staging"}, ErrorMessage: "Du befindest dich in keinem Feature Branch."},
		{Command: "git", Args: []string{"status", "--porcelain"}, Output: "Prüfe auf Änderungen, die nicht zum Commit vorgesehen sind…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Änderungen, die nicht zum Commit vorgesehen sind. Bitte Committe oder Stashe diese vor einem Merge."},
		{Command: "git", Args: []string{"ls-remote", "origin", currentBranch}, Output: "Prüfe auf Remote Branch…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Der Branch " + currentBranch + " existiert nicht auf remote."},
		{Command: "git", Args: []string{"fetch"}, Output: "Aktualisiere Branch Informationen…", Expectation: "*", Forbidden: []string{}, ErrorMessage: ""},
		{Command: "git", Args: []string{"log", "origin/" + currentBranch + ".." + currentBranch}, Output: "Prüfe auf nicht gepushte Commits…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt nicht gepushte Änderungen. Bitte pushe diese vor einem Merge."},
		{Command: "git", Args: []string{"log", currentBranch + ".." + "origin/" + currentBranch}, Output: "Prüfe auf nicht gemergte Commits…", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt nicht gemergte Änderungen. Bitte führe ein 'git pull --rebase' vor einem Merge aus."},
	}
}

func branchRebaseCommands() []commandrunner.CommandCheck {
	return []commandrunner.CommandCheck{
		{Command: "git", Args: []string{"checkout", "master"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Master checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"pull", "--rebase"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Master pull fehlgeschlagen."},
		{Command: "git", Args: []string{"checkout", "staging"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"pull", "--rebase"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging pull fehlgeschlagen."},
	}
}


func branchMergeCommands(featureBranch string) []commandrunner.CommandCheck {
	return []commandrunner.CommandCheck{
		{Command: "git", Args: []string{"checkout", featureBranch}, Output: "Merge Master Branch in " + featureBranch + "…", Expectation: "*", Forbidden: []string{}, ErrorMessage: featureBranch + " checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"merge", "master"}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Merge fehlgeschlagen."},
		{Command: "git", Args: []string{"status", "--porcelain"}, Output: "", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Merge Konflikte. Bitte behebe diese."},
		{Command: "git", Args: []string{"checkout", "staging"}, Output: "Merge " + featureBranch + " in staging…", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Staging checkout fehlgeschlagen."},
		{Command: "git", Args: []string{"merge", featureBranch}, Output: "", Expectation: "*", Forbidden: []string{}, ErrorMessage: "Merge fehlgeschlagen."},
		{Command: "git", Args: []string{"status", "--porcelain"}, Output: "", Expectation: "", Forbidden: []string{}, ErrorMessage: "Es gibt Merge Konflikte. Bitte behebe diese."},
	}
}
