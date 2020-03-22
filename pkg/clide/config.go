package clide

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// Config holds the global configuration for a clide demo
type Config struct {
	User           string    `json:"user"`
	Directory      string    `json:"directory"`
	TypeSpeed      int       `json:"typeSpeed"`
	Humanize       float64   `json:"humanize"`
	HideWarnings   bool      `json:"hideWarnings"`
	ClearBeforeAll bool      `json:"clearBeforeAll"`
	KeyTriggerAll  bool      `json:"keyTriggerAll"`
	Windows        []Window  `json:"windows"`
	TiggerKeys     []string  `json:"triggerKeys"`
	Commands       []Command `json:"commands"`
}

// Window holds data for a window created in sdl
type Window struct {
	Window *sdl.Window
	Name   string `json:"name"`
	X      int32  `json:"x"`
	Y      int32  `json:"y"`
	Height int32  `json:"height"`
	Width  int32  `json:"width"`
}

//Validate checks for potential issues in a Config and
//adds some default values if they are not present
func (cfg Config) Validate() {
	var window *sdl.Window
	var err error

	// Create a window for us to draw the text on
	if window, err = sdl.CreateWindow("Clide", 0, 0, 600, 600, sdl.WINDOW_SHOWN); err != nil {
		panic(err)
	}
	defer window.Destroy()

	//initialize typer values
	typer := Typer{
		Window: window,
		Pos: Position{
			X: 5,
			Y: 5,
			H: 0,
			W: 0,
		},
		Font: Font{
			Path: "assets/UbuntuMono-B.ttf",
			Size: 18,
		},
		Speed:    cfg.TypeSpeed,
		Humanize: cfg.Humanize,
	}

	//throw error when no commands are present
	if len(cfg.Commands) == 0 {
		Print(typer, "No commands found in provided json file")
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
			if !v.IsInstalled() {
				notInstalled = append(notInstalled, v.CmdString)
			}
		}

		//comfirm user wants to run program even though uninstalled commands will be skipped
		if len(notInstalled) != 0 {
			Print(typer, "WARNING: At least one command is not installed on the system! The following commands will be skipped:")
			for _, badCmd := range notInstalled {
				Print(typer, badCmd)
			}
		}
	}
}
