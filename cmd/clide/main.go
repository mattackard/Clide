package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gookit/color"
)

//Config holds the global configuration for a clide demo
type Config struct {
	User         string    `json:"user"`
	Directory    string    `json:"directory"`
	TypeSpeed    int       `json:"typeSpeed"`
	Humanize     float64   `json:"humanize"`
	HideWarnings bool      `json:"hideWarnings"`
	Commands     []Command `json:"commands"`
}

//Command holds a single clide command
type Command struct {
	CmdString string `json:"cmd"`
	Typed     bool   `json:"typed"`
	Window    string `json:"window"`
	PreDelay  int    `json:"predelay"`
	PostDelay int    `json:"postdelay"`
	Timeout   int    `json:"timeout"`
	Hidden    bool   `json:"hidden"`
}

func main() {
	//support missing file extension
	if !strings.HasSuffix(os.Args[1], ".json") {
		os.Args[1] += ".json"
	}

	//unmarshal clide json into config struct
	clide, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(clide)
	cfg := Config{}
	err = json.Unmarshal(bytes, &cfg)

	cfg.validate()

	//run each command in the commands slice
	for _, cmd := range cfg.Commands {
		cmd.validate()
		if cmd.isInstalled() {
			err = command(cmd, cfg)
			if err != nil {
				panic(err)
			}
		} else {
			if !cfg.HideWarnings {
				color.Printf("<yellow>WARNING</>: %s is not installed! Skipping command '%s'.\n", strings.Split(cmd.CmdString, " ")[0], cmd.CmdString)
			}
		}
	}
}

//validate checks for potential issues in a Config and
//adds some default values if they are not present
func (cfg Config) validate() {
	//throw error when no commands are present
	if len(cfg.Commands) == 0 {
		panic("No commands found in provided json file")
	}

	//default directory
	if cfg.Directory == "" {
		var err error
		cfg.Directory, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	//default user
	if cfg.User == "" {
		cfg.User = "demo@clide"
	}

	if !cfg.HideWarnings {
		//check if all commands are installed
		notInstalled := []string{}
		for _, v := range cfg.Commands {
			if !v.isInstalled() {
				notInstalled = append(notInstalled, v.CmdString)
			}
		}

		//comfirm user wants to run program even though uninstalled commands will be skipped
		if len(notInstalled) != 0 {
			color.Println("<yellow>WARNING</>: At least one command is not installed on the system! The following commands will be skipped:")
			for _, badCmd := range notInstalled {
				split := strings.Split(badCmd, " ")
				joined := strings.Join(split[1:], " ")
				color.Printf("\t<green>%s</> %s \n", split[0], joined)
			}
			fmt.Print("Would you still like to run the demo? [y/n]: ")

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answer := scanner.Text()
			fmt.Scan(answer)
			if !strings.HasPrefix(answer, "y") && !strings.HasPrefix(answer, "Y") {
				fmt.Println("Exiting program")
				os.Exit(1)
			}
		}
	}
}

//validate checks for potential issues in a Command
func (cmd Command) validate() {
	//throw an error when no command string is present
	if cmd.CmdString == "" {
		panic("No command string present in command")
	}
}

//isInstalled checks to see if the command is installed on the system
func (cmd Command) isInstalled() bool {
	if _, err := exec.LookPath(strings.Split(cmd.CmdString, " ")[0]); err != nil {
		return false
	}
	return true
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
