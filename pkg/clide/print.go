package clide

import (
	"time"

	"github.com/eiannone/keyboard"
)

//prompt returns a string used to emulate a terminal prompt
func prompt(cfg Config) string {
	return cfg.User + ":" + cfg.Directory + "> "
}

//writeCommand prints out the given command and emulates a terminal prompt before it
func writeCommand(cmd Command, cfg Config, typer Typer) (Typer, error) {

	//print terminal prompt
	typer.Pos.X = 5
	pos, err := Print(typer, prompt(cfg))
	if err != nil {
		return Typer{}, err
	}

	if cmd.WaitForKey || cfg.KeyTriggerAll {
		ListenForKey(cfg)
	} else {
		//wait before writing the command to the terminal
		time.Sleep(time.Duration(cmd.PreDelay) * time.Millisecond)
	}

	//set typer x position after command prompt
	typer.Pos = Position{
		X: pos.X,
		Y: typer.Pos.Y,
		H: pos.H,
		W: pos.W,
	}

	//type of print command to window
	if cmd.Typed {
		typer.Pos, err = Type(typer, cmd.CmdString)
		if err != nil {
			return Typer{}, err
		}
	} else {
		typer.Pos, err = Print(typer, cmd.CmdString)
		if err != nil {
			return Typer{}, err
		}
	}

	if cmd.WaitForKey || cfg.KeyTriggerAll {
		ListenForKey(cfg)
	} else {
		//wait before executing the command, but after writing to the terminal
		time.Sleep(time.Duration(cmd.PostDelay) * time.Millisecond)
	}

	return typer, nil
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
