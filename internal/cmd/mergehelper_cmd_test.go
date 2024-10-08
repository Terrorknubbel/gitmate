package cmd

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LocalRepoSuite struct {
	suite.Suite
	originalDir   string
	localRepoDir  string
	remoteRepoDir string
}

var stagingArgs = []string{
	"merge",
	"staging",
}

func CaptureStderr(f func(args []string) error, args []string) string {
	// Suppress stdout
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)

	err := f(args)
	if err != nil {
		return err.Error()
	}

	os.Stdout = oldStdout
	return ""
}

func (suite *LocalRepoSuite) SetupTest() {
	localRepoDir, _ := os.MkdirTemp("", "localRepoDir")

	originalDir, _ := os.Getwd()
	os.Chdir(localRepoDir)

	suite.localRepoDir = localRepoDir
	suite.originalDir = originalDir
}

func (suite *LocalRepoSuite) prepareLocalAndRemoteRepoWithFeatureBranch() {
	remoteRepo, _ := os.MkdirTemp("", "remoteRepo")
	os.Chdir(remoteRepo)
	exec.Command("git", "init").Run()

	os.Chdir(suite.localRepoDir)
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "--add", "remote.origin.url", remoteRepo).Run()
	exec.Command("git", "config", "--add", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*").Run()
	exec.Command("git", "checkout", "-b", "feature_branch").Run()

	suite.remoteRepoDir = remoteRepo
}

func (suite *LocalRepoSuite) TearDownTest() {
	os.RemoveAll(suite.localRepoDir)
	os.Chdir(suite.originalDir)
}

func (suite *LocalRepoSuite) Test_PrerequisiteCommands_NotInAGitDir() {
	output := CaptureStderr(suite.testConfig().execute, stagingArgs)

	suite.Contains(output, "Du befindest dich in keinem Git Verzeichnis.")
}

func (suite *LocalRepoSuite) Test_BranchConditionCommands_NotInAFeatureBranch() {
	exec.Command("git", "init").Run()

	output := CaptureStderr(suite.testConfig().execute, stagingArgs)
	suite.Contains(output, "Du befindest dich in keinem Feature Branch.")
}

func (suite *LocalRepoSuite) Test_BranchConditionCommands_CurrentBranchIsStaging() {
	exec.Command("git", "init").Run()
	exec.Command("git", "checkout", "-b", "staging").Run()

	output := CaptureStderr(suite.testConfig().execute, stagingArgs)

	suite.Contains(output, "Du befindest dich in keinem Feature Branch.")
}

func (suite *LocalRepoSuite) Test_BranchConditionCommands_UnstashedChanges() {
	suite.prepareLocalAndRemoteRepoWithFeatureBranch()

	outfile, _ := os.Create("./file1.txt")
	defer outfile.Close()

	outfile.WriteString("Commit 1")

	output := CaptureStderr(suite.testConfig().execute, stagingArgs)

	suite.Contains(output, "Es gibt Änderungen, die nicht zum Commit vorgesehen sind. Bitte Committe oder Stashe diese vor einem Merge.")
}

func (suite *LocalRepoSuite) Test_BranchConditionCommands_NoRemoteBranch() {
	exec.Command("git", "init").Run()
	exec.Command("git", "checkout", "-b", "feature_branch").Run()

	outfile, _ := os.Create("./file1.txt")
	defer outfile.Close()

	outfile.WriteString("Commit 1")

	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "commit1").Run()

	output := CaptureStderr(suite.testConfig().execute, stagingArgs)

	suite.Contains(output, "Der Branch feature_branch existiert nicht auf remote.")
}

func (suite *LocalRepoSuite) Test_BranchConditionCommands_UnstagedChanges() {
	suite.prepareLocalAndRemoteRepoWithFeatureBranch()

	exec.Command("git", "push", "--set-upstream", "origin", "feature_branch").Run()

	outfile, _ := os.Create("./file2.txt")
	defer outfile.Close()

	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "commit1").Run()

	output := CaptureStderr(suite.testConfig().execute, stagingArgs)

	suite.Contains(output, "Es gibt nicht gepushte Änderungen. Bitte pushe diese vor einem Merge.")
}

func (suite *LocalRepoSuite) Test_BranchConditionCommands_UnmergedChanges() {
	suite.prepareLocalAndRemoteRepoWithFeatureBranch()

	outfile, _ := os.Create("./localFile.txt")
	defer outfile.Close()

	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "localCommit").Run()
	exec.Command("git", "push", "--set-upstream", "origin", "feature_branch").Run()

	os.Chdir(suite.remoteRepoDir)
	exec.Command("git", "checkout", "feature_branch").Output()

	outfile, _ = os.Create("./remoteFile.txt")
	defer outfile.Close()

	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "remoteCommit").Run()

	os.Chdir(suite.localRepoDir)

	output := CaptureStderr(suite.testConfig().execute, stagingArgs)

	suite.Contains(output, "Es gibt nicht gemergte Änderungen. Bitte führe ein 'git pull --rebase' vor einem Merge aus.")
}

// TODO: Test_MergeMasterIntoFeatureBranch

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(LocalRepoSuite))
}

func (suite *LocalRepoSuite) testConfig() *Config {
	return newTestConfig(&suite.Suite)
}
