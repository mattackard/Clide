package clide

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//set up env

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

	testCfg := Config{
		Commands: []Command{
			testCmd,
		},
	}

	newCfg, err := testCfg.Validate()
	if err != nil {
		t.Errorf("Validate config failed with valid config: %s", err.Error())
	}

	//test for default user
	if newCfg.User == "" {
		t.Errorf("Default config user not applied")
	}

	//check for default directory
	currentDir, _ := os.Getwd()
	if newCfg.Directory != currentDir {
		t.Errorf("Default config directory not applied. Expected %s to equal %s", newCfg.Directory, currentDir)
	}

	//check for default window
	if len(newCfg.Windows) == 0 {
		t.Error("Default config window not applied")
	}

}
