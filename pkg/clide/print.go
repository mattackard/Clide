package clide

import (
	"fmt"
	"math/rand"
	"time"
)

//prompt returns a string used to emulate a terminal prompt
func prompt(cfg Config) string {
	return cfg.User + ":" + cfg.Directory + "> "
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
