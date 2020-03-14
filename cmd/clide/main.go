package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

//Config holds the global configuration for a clide demo
type Config struct {
	User      string    `json:"user"`
	Directory string    `json:"directory"`
	TypeSpeed int       `json:"typeSpeed"`
	Humanize  float32   `json:"humanize"`
	Commands  []Command `json:"commands"`
}

//Command holds a single clide command
type Command struct {
	CmdString string `json:"cmd"`
	Typed     bool   `json:"typed"`
	Window    string `json:"window"`
	PreDelay  int    `json:"predelay"`
	PostDelay int    `json:"postdelay"`
}

func main() {
	//unmarshal clide json into clide struct
	clide, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(clide)

	cfg := Config{}
	err = json.Unmarshal(bytes, &cfg)

	//run each command in the clide slice
	for _, cmd := range cfg.Commands {
		err = command(cmd, cfg)
		if err != nil {
			panic(err)
		}
	}

}

//prompt returns a string used to emulate a terminal prompt
func prompt(cfg Config) string {
	return cfg.User + ":" + cfg.Directory + "> "
}

//command runs a cli command with options to wait before and after execution
func command(cmd Command, cfg Config) error {
	split := strings.Split(cmd.CmdString, " ")
	program := split[0]
	command := exec.Command(program, split[1:]...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	//wait before writing the command to the terminal
	time.Sleep(time.Duration(cmd.PreDelay) * time.Millisecond)
	writeCommand(cmd.CmdString, cfg)

	//wait before executing the command, but after writing to the terminal
	time.Sleep(time.Duration(cmd.PostDelay) * time.Millisecond)
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}

//writeCommand prints out the given command and emulates a terminal prompt before it
func writeCommand(cmd string, cfg Config) {
	fmt.Println(prompt(cfg) + cmd)
}
