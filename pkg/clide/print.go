package clide

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

//prompt returns a string used to emulate a terminal prompt
func prompt(cfg Config) string {
	return cfg.User + ":" + cfg.Directory + "> "
}

//writeCommand prints out the given command and emulates a terminal prompt before it
func writeCommand(cmd Command, cfg Config, done chan bool) {

	fmt.Print(prompt(cfg))

	if cmd.WaitForKey {
		err := keyboard.Open()
		if err != nil {
			panic(err)
		}
		defer keyboard.Close()

		//keep looping until key from TriggerKeys is pressed
		for {
			if keyPressed(cfg.TiggerKeys) {
				break
			}
		}
	} else {
		//wait before writing the command to the terminal
		time.Sleep(time.Duration(cmd.PreDelay) * time.Millisecond)
	}

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

//keyPressed returns whether or not the key pressed in any []string
func keyPressed(keys []string) bool {
	char, _, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		//convert key(string) into rune slice
		runes := []rune(key)

		//check all runes in slice against keyboard char
		for _, r := range runes {
			if r == char {
				return true
			}
		}
	}
	return false
}
