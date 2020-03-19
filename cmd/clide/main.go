package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattackard/Clide/pkg/clide"
	"github.com/mattackard/Clide/pkg/sdltyper"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	fontPath = "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
	fontSize = 16
	helpText = `Clide CLI Usage:
		clide example.json		runs the clide demo stored in example.json
		clide-sh script.sh		converts script.sh into script.json formatted as a clide demo
		clide-build			opens the clide demo builder interface				
		clide				shows this help message`
)

func main() {
	//initialize sdl2
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer ttf.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	//open a window for each defined in json
	window, err := newWindow("Clide", sdltyper.Position{
		X: 0,
		Y: 0,
		H: 800,
		W: 1000,
	})
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	//initialize typer values
	typer := sdltyper.Typer{
		Window: window,
		Pos: sdltyper.Position{
			X: 5,
			Y: 5,
			H: 0,
			W: 0,
		},
		Font: sdltyper.Font{
			Path: fontPath,
			Size: fontSize,
		},
		Speed:    100,
		Humanize: 0.9,
	}

	//check if os.Args[1] exists
	if len(os.Args) < 2 {
		fmt.Println("You must provide a clide configured json file to run a demo.")
		fmt.Println(typer.Pos)
		typer.Pos, err = sdltyper.Print(typer, "You must provide a clide configured json file to run a demo.")
		if err != nil {
			panic(err)
		}
		fmt.Println("\n" + helpText)
		fmt.Println(typer.Pos)
		typer.Pos, err = sdltyper.Print(typer, helpText)
		if err != nil {
			panic(err)
		}
		fmt.Println(typer.Pos)
		typer.Pos, err = sdltyper.Print(typer, "Exiting in 10 seconds")
		if err != nil {
			panic(err)
		}
		exit(10)
	}

	//support missing file extension
	if !strings.HasSuffix(os.Args[1], ".json") {
		os.Args[1] += ".json"
	}

	//check if os.Args[1] file exists
	var resp *http.Response
	file, err := os.Open(os.Args[1])
	if err != nil {
		errorText := fmt.Sprintf("File %s does not exists in current directory, checking /usr/share/clide", os.Args[1])
		typer.Pos, err = sdltyper.Print(typer, errorText)
		if err != nil {
			panic(err)
		}

		//if not check usr/share/clide
		file, err = os.Open("/usr/share/clide/" + os.Args[1])
		if err != nil {
			errorText = fmt.Sprintf("File %s does not exists /usr/share/clide. Checking for clide examples on clide.sh with name %s ...", os.Args[1], os.Args[1])
			typer.Pos, err = sdltyper.Print(typer, errorText)
			if err != nil {
				panic(err)
			}

			//if not finally check clide demo fileserver
			resp, err = http.Get("https://clide.sh/demos/" + os.Args[1])
			if err != nil || resp.StatusCode != 200 {
				typer.Pos, err = sdltyper.Print(typer, "Could not find file at clide.sh/demos")
				if err != nil {
					panic(err)
				}
				exit(10)
			}
		}
	}

	//unmarshal clide json into config struct
	var bytes []byte
	cfg := clide.Config{}
	if file != nil {
		bytes, err = ioutil.ReadAll(file)
	} else if resp != nil {
		bytes, err = ioutil.ReadAll(resp.Body)
	} else {
		panic("File and HTTP Response are both nil : " + err.Error())
	}
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	cfg.Validate()

	//adjust typer values to match cfg
	typer.Speed = cfg.TypeSpeed
	typer.Humanize = cfg.Humanize

	//run each command in the commands slice
	for _, cmd := range cfg.Commands {
		cmd.Validate()
		if cmd.IsInstalled() {
			typer, err = cmd.Run(cfg, typer)
			if err != nil {
				panic(err)
			}
		} else {
			if !cfg.HideWarnings {
				warning := fmt.Sprintf("WARNING: %s is not installed! Skipping command: '%s'.\n", strings.Split(cmd.CmdString, " ")[0], cmd.CmdString)
				typer.Pos, err = sdltyper.Print(typer, warning)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// newWindow creates a new sdl2 window
func newWindow(title string, pos sdltyper.Position) (*sdl.Window, error) {
	var window *sdl.Window
	var err error

	// Create a window for us to draw the text on
	if window, err = sdl.CreateWindow(title, pos.X, pos.Y, pos.W, pos.H, sdl.WINDOW_SHOWN); err != nil {
		return nil, err
	}

	return window, nil
}

// exit stops the program after delaying for x seconds
func exit(delay int) {
	time.Sleep(time.Second * time.Duration(delay))
	os.Exit(1)
}
