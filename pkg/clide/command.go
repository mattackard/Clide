package clide

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

//Command holds a single clide command
type Command struct {
	CmdString  string `json:"cmd"`
	Typed      bool   `json:"typed"`
	Window     string `json:"window"`
	PreDelay   int    `json:"predelay"`
	PostDelay  int    `json:"postdelay"`
	Timeout    int    `json:"timeout"`
	Hidden     bool   `json:"hidden"`
	WaitForKey bool   `json:"waitForKey"`
}

//Validate checks for potential issues in a Command
func (cmd Command) Validate() {
	//throw an error when no command string is present
	if cmd.CmdString == "" {
		panic("No command string present in command")
	}
}

//IsInstalled checks to see if the command is installed on the system
func (cmd Command) IsInstalled() bool {
	program := strings.Split(cmd.CmdString, " ")[0]
	if _, err := exec.LookPath(program); err != nil {
		return false
	}
	return true
}

//Run runs a cli command with options to wait before and after execution
func (cmd Command) Run(cfg Config) error {
	split := strings.Split(cmd.CmdString, " ")
	program := split[0]
	command := exec.Command(program, split[1:]...)

	if !cmd.Hidden {
		command.Stderr = os.Stderr
		command.Stdout = os.Stdout

		//type the command into the console and wait for it to finish typing before further execution
		written := make(chan bool, 1)
		go writeCommand(cmd, cfg, written)
		<-written
	}

	if cmd.Timeout != 0 {
		err := command.Start()
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(cmd.Timeout) * time.Second)
		command.Process.Kill()
	} else {
		err := command.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
