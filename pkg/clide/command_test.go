package clide

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()

	//tear down env
	os.Exit(exitCode)
}

func TestIsInstalled(t *testing.T) {
	uninstalled := Command{
		CmdString: "uninstalled command",
	}
	installed := Command{
		CmdString: "echo Testing",
	}

	if uninstalled.IsInstalled() {
		t.Error("Uninstalled command should return false")
	}

	if !installed.IsInstalled() {
		t.Error("Installed command should return true")
	}
}

// TestValidate tests command.Validate() and config.Validate()
func TestValidate(t *testing.T) {
	noCmd := Command{
		CmdString: "",
	}
	testCmd := Command{
		CmdString: "echo Testing",
	}

	err := noCmd.Validate()
	if err == nil {
		t.Error("No error thrown using empty CmdString")
	}

	err = testCmd.Validate()
	if err != nil {
		t.Error("Error thrown for valid CmdString")
	}
}
