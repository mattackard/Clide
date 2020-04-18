package clide

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Command holds a single clide command
type Command struct {
	CmdString      string   `json:"cmd"`
	Typed          bool     `json:"typed"`
	Window         string   `json:"window"`
	PreDelay       int      `json:"predelay"`
	PostDelay      int      `json:"postdelay"`
	Timeout        int      `json:"timeout"`
	Hidden         bool     `json:"hidden"`
	WaitForKey     bool     `json:"waitForKey"`
	ClearBeforeRun bool     `json:"clearBeforeRun"`
	Async          bool     `json:"async"`
	HideWindow     bool     `json:"hideWindow"`
	ResizeWindows  []Window `json:"resizeWindows"`
}

const defaultTimeout = 10

// Validate checks for potential issues in a Command
func (cmd Command) Validate() error {
	// throw an error when no command string is present
	if cmd.CmdString == "" {
		err := errors.New("No command string present in command")
		return err
	}
	return nil
}

// IsInstalled checks to see if the command is installed on the system
func (cmd Command) IsInstalled() bool {
	program := strings.Split(cmd.CmdString, " ")[0]
	if _, err := exec.LookPath(program); err != nil {
		if program != "cd" {
			return false
		}
	}
	return true
}

// Run runs a cli command with options to wait before and after execution
func (cmd Command) Run(cfg *Config, typer *Typer, exitChan chan bool) error {

	err := cmd.Validate()
	if err != nil {
		typer.Print(err.Error(), sdl.Color{R: 255, G: 0, B: 0, A: 255})
	}

	if !cmd.IsInstalled() {
		if !cfg.HideWarnings {
			warning := fmt.Sprintf("WARNING: %s is not installed! Skipping command: '%s'.\n", strings.Split(cmd.CmdString, " ")[0], cmd.CmdString)
			err = typer.Print(warning, sdl.Color{R: 255, G: 0, B: 0, A: 255})
			if err != nil {
				panic(err)
			}
		}
	}

	// if resize windows is set, get each window and resize it as set
	if len(cmd.ResizeWindows) > 0 {
		for _, win := range cmd.ResizeWindows {
			win.Window.Hide()

			target := cfg.GetWindow(win.Name)
			target.SetPosition(win.X, win.Y)
			target.SetSize(win.Width, win.Height)

			cfg.ClearAllWindows()
		}
	}
	typer.Window.Show()

	// clear terminal if set in config or command
	if cmd.ClearBeforeRun || cfg.ClearBeforeAll {
		var err error
		bgColor, err := StringToColor(cfg.ColorScheme.TerminalBG)
		if err != nil {
			return err
		}

		err = typer.ClearWindow(bgColor)
		if err != nil {
			return err
		}
	}

	// get text color
	textColor, err := StringToColor(cfg.ColorScheme.PrimaryText)
	if err != nil {
		return err
	}

	// parse program from command string
	split := strings.Split(cmd.CmdString, " ")
	program := split[0]

	var command *exec.Cmd
	if strings.ContainsAny(cmd.CmdString, "|><") {
		command = exec.Command("bash", "-e", "-c", cmd.CmdString)
	} else {
		command = exec.Command(program, split[1:]...)
	}

	if cmd.Hidden {
		err := command.Run()
		if err != nil {
			return err
		}
	} else if cmd.Async {
		go cmd.runAsync(command, cfg, typer, textColor, exitChan)
	} else {
		// type the command into the console and wait for it to finish typing before further execution
		var err error
		err = writeCommand(cmd, *cfg, typer)
		if err != nil {
			panic(err)
		}

		// special handling for cd commands
		if program == "cd" {
			updateDirectory(cfg, split[1])
			err := os.Chdir(split[1])
			if err != nil {
				typer.Pos.X = 5
				err = typer.Print(err.Error(), textColor)
				if err != nil {
					return err
				}
			}
		} else {
			// set up a stdout and stderr pipe to capture the output
			output, err := command.StdoutPipe()
			if err != nil {
				return err
			}
			errOutput, err := command.StderrPipe()
			if err != nil {
				return err
			}

			// dont wait for command to finish
			err = command.Start()
			if err != nil {
				return err
			}

			// stream the output from the command in realtime
			// won't block so the sleep timer can run while printing
			go printOutputs(output, errOutput, typer, textColor)

			if cmd.Timeout != 0 {
				time.Sleep(time.Duration(cmd.Timeout) * time.Second)
				command.Process.Signal(os.Interrupt)
			} else {
				command.Wait()
			}

			// give any error messages some time to print to the window
			time.Sleep(time.Millisecond * 10)
		}

	}

	if !cmd.Async {
		// hide the window if cmd.HideWindow is set
		if cmd.HideWindow {
			typer.Window.Hide()
		}
	}

	return nil
}

// runAsync runs a command asynchronously
func (cmd Command) runAsync(command *exec.Cmd, cfg *Config, typer *Typer, textColor sdl.Color, exitChan chan bool) error {
	// type the command into the console and wait for it to finish typing before further execution
	var err error
	err = writeCommand(cmd, *cfg, typer)
	if err != nil {
		return err
	}

	// set up a stdout and stderr pipe to capture the output
	output, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	errOutput, err := command.StderrPipe()
	if err != nil {
		return err
	}

	// dont wait for command to finish
	err = command.Start()
	if err != nil {
		return err
	}

	// make sure process is not left running
	go func() {
		for {
			select {
			case <-exitChan:
				command.Process.Signal(os.Interrupt)
			}
		}
	}()

	// stream the output from the command in realtime
	// won't block so the sleep timer can run while printing
	go printOutputs(output, errOutput, typer, textColor)

	// set default timeout if not set
	if cmd.Timeout == 0 {
		cmd.Timeout = defaultTimeout
	}

	time.Sleep(time.Duration(cmd.Timeout) * time.Second)

	// hide the window if cmd.HideWindow is set
	if cmd.HideWindow {
		typer.Window.Hide()
	}

	command.Process.Signal(os.Interrupt)

	// give any error messages some time to print to the window
	time.Sleep(time.Millisecond * 10)
	return nil
}

func printOutputs(output io.ReadCloser, errOutput io.ReadCloser, typer *Typer, textColor sdl.Color) {
	go func() {
		scanner := bufio.NewScanner(output)

		for scanner.Scan() {
			typer.mutex.Lock()
			line := scanner.Text()
			typer.Pos.X = 5
			typer.Print(line, textColor)
			typer.mutex.Unlock()

		}
	}()
	go func() {
		errScanner := bufio.NewScanner(errOutput)

		for errScanner.Scan() {
			typer.mutex.Lock()
			line := errScanner.Text()
			typer.Pos.X = 5
			typer.Print(line, textColor)
			typer.mutex.Unlock()
		}
	}()
}

// ClearWindow removes all content on the window specified in typer
func (typer *Typer) ClearWindow(color sdl.Color) error {
	var surface *sdl.Surface
	var err error

	// get surface info
	surface, err = typer.Window.GetSurface()
	if err != nil {
		return err
	}

	err = surface.Blit(nil, surface, &sdl.Rect{X: 0, Y: -typer.Pos.Y, W: 0, H: 0})
	if err != nil {
		return err
	}

	// create a rectangle that fills the screen and make it black
	rect := sdl.Rect{
		X: 0,
		Y: 0,
		W: surface.W,
		H: surface.H,
	}

	// set color to correct for the weirdness in the Uint32 conversion
	colorFix := sdl.Color{
		R: color.A,
		G: color.R,
		B: color.G,
		A: color.B,
	}
	surface.FillRect(&rect, colorFix.Uint32())

	// draw the rect and update typer position
	typer.Window.UpdateSurface()
	typer.Pos = Position{
		X: 5,
		Y: 5,
		H: 0,
		W: 0,
	}

	return nil
}

// updateDirectory updates the config directory printed to the prompt for when a cd command is called
func updateDirectory(cfg *Config, path string) {
	if !strings.HasSuffix(cfg.Directory, "/") {
		cfg.Directory += "/"
	}

	tempPath := cfg.Directory + path

	// split path by slash to edit relative paths
	if strings.Contains(tempPath, ".") {

		// keep removing directories if a ../ is present
		for strings.Contains(tempPath, "..") {
			tempPath = removeOnePath(tempPath)
		}

		// ignore all ./ relative paths
		strings.ReplaceAll(tempPath, "./", "")

		if !strings.HasSuffix(tempPath, "/") {
			tempPath += "/"
		}
		cfg.Directory = tempPath

	} else {
		if !strings.HasSuffix(tempPath, "/") {
			path += "/"
		}
		cfg.Directory += path
	}
}

// removeOnePath removes a single pair of ../ and its preceding directory
func removeOnePath(path string) string {
	split := strings.Split(path, "/")
	for i, dir := range split {
		if dir == ".." {
			// overwrite a ../ path and its predecessor with the remaining path
			copy(split[i-1:], split[i+1:])

			// empty the old indexes
			split[len(split)-1] = ""
			split[len(split)-2] = ""

			// remove unused length
			split = split[:len(split)-2]
			break
		}
	}
	return strings.Join(split, "/")
}
