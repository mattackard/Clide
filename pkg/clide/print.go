package clide

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

//prompt returns a string used to emulate a terminal prompt
func printPrompt(cfg Config, typer Typer) (Typer, error) {

	//set colors for user, directory, and primary text
	userColor, err := StringToColor(cfg.ColorScheme.UserText)
	directoryColor, err := StringToColor(cfg.ColorScheme.DirectoryText)
	primaryColor, err := StringToColor(cfg.ColorScheme.PrimaryText)

	//print promt to terminal window
	pos, err := Print(typer, cfg.User, userColor)
	if err != nil {
		return Typer{}, err
	}
	typer.Pos.X += pos.X - 6
	pos, err = Print(typer, ":", primaryColor)
	if err != nil {
		return Typer{}, err
	}
	typer.Pos.X += pos.X - 6
	pos, err = Print(typer, cfg.Directory, directoryColor)
	if err != nil {
		return Typer{}, err
	}
	typer.Pos.X += pos.X - 6
	pos, err = Print(typer, "$ ", primaryColor)
	if err != nil {
		return Typer{}, err
	}
	typer.Pos.X += pos.X - 6
	return typer, nil
}

//writeCommand prints out the given command and emulates a terminal prompt before it
func writeCommand(cmd Command, cfg Config, typer Typer) (Typer, error) {

	//print terminal prompt
	typer.Pos.X = 5
	typer, err := printPrompt(cfg, typer)
	if err != nil {
		return Typer{}, err
	}

	if cmd.WaitForKey || cfg.KeyTriggerAll {
		if len(cfg.TiggerKeys) == 0 {
			pos, err := Print(typer, "WaitForKey or KeyTriggerAll is set, but no TriggerKeys are defined!", sdl.Color{R: 255, G: 0, B: 0, A: 255})
			if err != nil {
				return Typer{}, err
			}

			typer.Pos.Y = pos.Y
			typer.Pos.X = 5
			typer, err = printPrompt(cfg, typer)
			if err != nil {
				return Typer{}, err
			}
		} else {
			ListenForKey(cfg)
		}
	} else {
		//wait before writing the command to the terminal
		time.Sleep(time.Duration(cmd.PreDelay) * time.Millisecond)
	}

	primaryColor, err := StringToColor(cfg.ColorScheme.PrimaryText)
	if err != nil {
		return Typer{}, err
	}

	//type of print command to window
	if cmd.Typed {
		typer.Pos, err = Type(typer, cmd.CmdString, primaryColor)
		if err != nil {
			return Typer{}, err
		}
	} else {
		typer.Pos, err = Print(typer, cmd.CmdString, primaryColor)
		if err != nil {
			return Typer{}, err
		}
	}

	if cmd.WaitForKey || cfg.KeyTriggerAll {
		if len(cfg.TiggerKeys) != 0 {
			ListenForKey(cfg)
		}
	} else {
		//wait before executing the command, but after writing to the terminal
		time.Sleep(time.Duration(cmd.PostDelay) * time.Millisecond)
	}

	return typer, nil
}
