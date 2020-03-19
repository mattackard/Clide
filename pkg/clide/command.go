package clide

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/mattackard/Clide/pkg/sdltyper"
)

//Command holds a single clide command
type Command struct {
	CmdString      string `json:"cmd"`
	Typed          bool   `json:"typed"`
	Window         int    `json:"window"`
	PreDelay       int    `json:"predelay"`
	PostDelay      int    `json:"postdelay"`
	Timeout        int    `json:"timeout"`
	Hidden         bool   `json:"hidden"`
	WaitForKey     bool   `json:"waitForKey"`
	ClearBeforeRun bool   `json:"clearBeforeRun"`
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
func (cmd Command) Run(cfg Config, typer sdltyper.Typer) (sdltyper.Typer, error) {
	//clear terminal if set in config or command
	if cmd.ClearBeforeRun || cfg.ClearBeforeAll {
		var err error
		typer, err = clearTerminal(typer)
		if err != nil {
			return typer, err
		}
	}

	//parse program from command string
	split := strings.Split(cmd.CmdString, " ")
	program := split[0]
	command := exec.Command(program, split[1:]...)

	if cmd.Hidden {
		err := command.Run()
		if err != nil {
			return sdltyper.Typer{}, nil
		}
	} else {
		command.Stderr = os.Stderr

		//type the command into the console and wait for it to finish typing before further execution
		var err error
		typer, err = writeCommand(cmd, cfg, typer)
		if err != nil {
			panic(err)
		}
		if cmd.Timeout != 0 {
			//set up a stdout pipe to capture the output
			output, err := command.StdoutPipe()
			if err != nil {
				return sdltyper.Typer{}, nil
			}

			//dont wait for command to finish
			err = command.Start()
			if err != nil {
				return sdltyper.Typer{}, err
			}

			//stream the output from the command in realtime
			//won't block so the sleep timer can run while printing
			go func() {
				scanner := bufio.NewScanner(output)
				for scanner.Scan() {
					line := scanner.Text()
					typer.Pos.X = 5
					typer.Pos, err = sdltyper.Print(typer, line)
				}
			}()

			time.Sleep(time.Duration(cmd.Timeout) * time.Second)
			command.Process.Kill()
		} else {
			output, err := command.Output()
			if err != nil {
				return sdltyper.Typer{}, err
			}
			typer.Pos.X = 5
			typer.Pos, err = sdltyper.Print(typer, string(output))
			if err != nil {
				return sdltyper.Typer{}, err
			}
		}
	}
	return typer, nil
}

func clearTerminal(typer sdltyper.Typer) (sdltyper.Typer, error) {
	var surface *sdl.Surface
	var err error

	//get surface info
	if surface, err = typer.Window.GetSurface(); err != nil {
		return typer, err
	}

	//create a rectangle that fills the screen and make it black
	rect := sdl.Rect{
		X: 0,
		Y: 0,
		W: surface.W,
		H: surface.H,
	}
	surface.FillRect(&rect, 0)

	//draw the rect and update typer position
	typer.Window.UpdateSurface()
	typer.Pos = sdltyper.Position{
		X: 5,
		Y: 5,
		H: 0,
		W: 0,
	}

	// command := exec.Command("clear")
	// command.Stdout = os.Stdout
	// command.Run()

	return typer, nil
}
