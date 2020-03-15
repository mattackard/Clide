package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
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
	Humanize  float64   `json:"humanize"`
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

	//type the command into the console and wait for it to finish typing before further execution
	written := make(chan bool, 1)
	go writeCommand(cmd, cfg, written)
	<-written

	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}

//writeCommand prints out the given command and emulates a terminal prompt before it
func writeCommand(cmd Command, cfg Config, done chan bool) {

	fmt.Print(prompt(cfg))

	//wait before writing the command to the terminal
	time.Sleep(time.Duration(cmd.PreDelay) * time.Millisecond)

	if cmd.Typed {
		//print each command character using the typespeed and humanize values
		for i, v := range cmd.CmdString {
			time.Sleep(getKeyDelay(cfg))
			fmt.Print(string(v))
			if i == len(cmd.CmdString)-1 {
				//wait before executing the command, but after writing to the terminal
				time.Sleep(time.Duration(cmd.PostDelay) * time.Millisecond)
				fmt.Print("\n")
				done <- true
			}
		}
	} else {
		fmt.Print(cmd.CmdString)

		//wait before executing the command, but after writing to the terminal
		time.Sleep(time.Duration(cmd.PostDelay) * time.Millisecond)
		fmt.Print("\n")
		done <- true
	}
}

//getKeyDelay calculates and returns a time to wait based on type speed and humanization ratio
func getKeyDelay(cfg Config) time.Duration {
	if cfg.Humanize > 0 {
		//set up a seeded random
		rand.Seed(time.Now().UnixNano())

		//calculate speed variance based on humanize field
		variance := (1 - cfg.Humanize - rand.Float64()) * float64(cfg.TypeSpeed)

		return time.Duration(float64(cfg.TypeSpeed)+variance) * time.Millisecond
	}
	return time.Duration(cfg.TypeSpeed) * time.Millisecond
}
