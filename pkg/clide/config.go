package clide

import (
	"errors"
	"os"
	"strconv"
	"strings"

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
	HideWindows    bool      `json:"hideWindows"`
	FontPath       string    `json:"fontPath"`
	FontSize       int       `json:"fontSize"`
	Windows        []Window  `json:"windows"`
	TiggerKeys     []string  `json:"triggerKeys"`
	ColorScheme    Colors    `json:"colorScheme"`
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

// Colors holds information for the color scheme of the terminal window and text
type Colors struct {
	UserText      string `json:"userText"`
	DirectoryText string `json:"directoryText"`
	PrimaryText   string `json:"primaryText"`
	TerminalBG    string `json:"terminalBG"`
}

//Validate checks for potential issues in a Config and
//adds some default values if they are not present
func (cfg Config) Validate() (Config, error) {
	var window *sdl.Window
	var err error

	// Create a window to draw the text on
	if window, err = sdl.CreateWindow("Clide", 0, 0, 600, 600, sdl.WINDOW_SHOWN); err != nil {
		return Config{}, err
	}
	defer window.Destroy()

	//font defaults
	if cfg.FontPath == "" {
		cfg.FontPath = "/usr/share/clide/assets/UbuntuMono-B.ttf"
	}
	if cfg.FontSize == 0 {
		cfg.FontSize = 18
	}

	//initialize typer values
	typer := Typer{
		Window: window,
		Pos: Position{
			X: 50,
			Y: 5,
			H: 0,
			W: 0,
		},
		Font: Font{
			Path: cfg.FontPath,
			Size: cfg.FontSize,
		},
		Speed:    cfg.TypeSpeed,
		Humanize: cfg.Humanize,
	}

	//if no windows were provided in json, create one
	if len(cfg.Windows) == 0 {
		cfg.Windows = []Window{Window{
			Name:   "Clide",
			X:      0,
			Y:      0,
			Height: 600,
			Width:  1000,
		}}
	}

	//throw error when no commands are present
	if len(cfg.Commands) == 0 {
		typer.Print("No commands found in provided json file", sdl.Color{R: 255, G: 0, B: 0, A: 255})
		return cfg, errors.New("No commands found in provided json file")
	}

	//default directory
	if cfg.Directory == "" {
		var err error
		cfg.Directory, err = os.Getwd()
		if err != nil {
			return cfg, err
		}
	}

	//default user
	if cfg.User == "" {
		cfg.User = "demo@clide"
	}

	//default color scheme
	if cfg.ColorScheme.UserText == "" || cfg.ColorScheme.DirectoryText == "" || cfg.ColorScheme.PrimaryText == "" || cfg.ColorScheme.TerminalBG == "" {
		cfg.ColorScheme = Colors{
			UserText:      "0,150,255,255",
			DirectoryText: "150,255,150,255",
			PrimaryText:   "220,220,220,255",
			TerminalBG:    "30,30,30,255",
		}
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
			typer.Print("WARNING: At least one command is not installed on the system! The following commands will be skipped:", sdl.Color{R: 255, G: 0, B: 0, A: 255})
			for _, badCmd := range notInstalled {
				typer.Print(badCmd, sdl.Color{R: 255, G: 0, B: 0, A: 255})
			}
		}
	}
	return cfg, nil
}

//getWindow return the sindow object with the specified name
func (cfg Config) getWindow(name string) *sdl.Window {
	var targetWindow *sdl.Window
	for _, win := range cfg.Windows {
		if win.Name == name {
			targetWindow = win.Window
		}
	}
	return targetWindow
}

//StringToColor converts a rgb or rgba formatted string to an sdl.Color struct
func StringToColor(color string) (sdl.Color, error) {
	var sdlColor sdl.Color
	split := strings.Split(color, ",")

	//set either rgb or rgba
	switch len(split) {
	case 3:
		r, err := strconv.Atoi(split[0])
		g, err := strconv.Atoi(split[1])
		b, err := strconv.Atoi(split[2])
		if err != nil {
			return sdl.Color{}, err
		}
		sdlColor = sdl.Color{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: 255,
		}
	case 4:
		r, err := strconv.Atoi(split[0])
		g, err := strconv.Atoi(split[1])
		b, err := strconv.Atoi(split[2])
		a, err := strconv.Atoi(split[3])
		if err != nil {
			return sdl.Color{}, err
		}
		sdlColor = sdl.Color{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(a),
		}
	default:
		err := errors.New("Invalid olor value")
		return sdl.Color{}, err
	}

	return sdlColor, nil
}
