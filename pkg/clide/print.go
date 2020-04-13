package clide

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

//prompt returns a string used to emulate a terminal prompt
func printPrompt(cfg Config, typer *Typer) error {

	//set colors for user, directory, and primary text
	userColor, err := StringToColor(cfg.ColorScheme.UserText)
	directoryColor, err := StringToColor(cfg.ColorScheme.DirectoryText)
	primaryColor, err := StringToColor(cfg.ColorScheme.PrimaryText)
	if err != nil {
		return err
	}

	lineY := typer.Pos.Y

	//print promt to terminal window using the user specified color for each section of the prompt
	err = typer.Print(cfg.User, userColor)
	if err != nil {
		return err
	}
	//if the line is scrolled up from the bottom, y pos will be modified in print and needs to be reset a line higher
	if lineY == typer.Pos.Y {
		lineY -= int32(cfg.FontSize) + 2
	}
	typer.Pos.Y = lineY

	err = typer.Print(":", primaryColor)
	if err != nil {
		return err
	}
	if lineY == typer.Pos.Y {
		lineY -= int32(cfg.FontSize) + 2
	}
	typer.Pos.Y = lineY

	err = typer.Print(cfg.Directory, directoryColor)
	if err != nil {
		return err
	}
	if lineY == typer.Pos.Y {
		lineY -= int32(cfg.FontSize) + 2
	}
	typer.Pos.Y = lineY

	err = typer.Print("$ ", primaryColor)
	if err != nil {
		return err
	}
	if lineY == typer.Pos.Y {
		lineY -= int32(cfg.FontSize) + 2
	}
	typer.Pos.Y = lineY

	return nil
}

//writeCommand prints out the given command and emulates a terminal prompt before it
func writeCommand(cmd Command, cfg Config, typer *Typer) error {

	//print terminal prompt
	typer.Pos.X = 5
	err := printPrompt(cfg, typer)
	if err != nil {
		return err
	}

	if cmd.WaitForKey || cfg.KeyTriggerAll {
		if len(cfg.TiggerKeys) == 0 {
			err := typer.Print("WaitForKey or KeyTriggerAll is set, but no TriggerKeys are defined!", sdl.Color{R: 255, G: 0, B: 0, A: 255})
			if err != nil {
				return err
			}

			typer.Pos.X = 5
			err = printPrompt(cfg, typer)
			if err != nil {
				return err
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
		return err
	}

	//type of print command to window
	if cmd.Typed {
		err = typer.Type(cmd.CmdString, primaryColor)
		if err != nil {
			return err
		}
	} else {
		err = typer.Print(cmd.CmdString, primaryColor)
		if err != nil {
			return err
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

	return nil
}
