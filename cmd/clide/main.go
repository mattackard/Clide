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
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	goroutineMax = 100
	fontPath     = "assets/UbuntuMono-B.ttf"
	fontSize     = 18
	helpText     = `Clide CLI Usage:
		clide example.json		runs the clide demo stored in example.json
		clide-sh script.sh		converts script.sh into script.json formatted as a clide demo
		clide-build			opens the clide demo builder interface				
		clide				shows this help message`
)

var goroutineCount int

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

	//create a channel to exit out of goroutines when program exits
	exitChan := make(chan bool, goroutineMax)

	//check if os.Args[1] exists
	if len(os.Args) < 2 {
		//open a window for error message
		window, err := newWindow("Clide", clide.Position{
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
		typer := clide.Typer{
			Window: window,
			Pos: clide.Position{
				X: 5,
				Y: 5,
				H: 0,
				W: 0,
			},
			Font: clide.Font{
				Path: fontPath,
				Size: fontSize,
			},
			Speed:    100,
			Humanize: 0.9,
		}

		fmt.Println("You must provide a clide configured json file to run a demo.")
		typer.Pos, err = clide.Print(typer, "You must provide a clide configured json file to run a demo.", sdl.Color{R: 255, G: 0, B: 0, A: 255})
		if err != nil {
			panic(err)
		}
		fmt.Println("\n" + helpText)
		typer.Pos.X = 20
		typer.Pos.Y += 20
		typer.Pos, err = clide.Print(typer, helpText, sdl.Color{R: 255, G: 0, B: 0, A: 255})
		if err != nil {
			panic(err)
		}
		typer.Pos.X = 5
		typer.Pos.Y += 20
		typer.Pos, err = clide.Print(typer, "Exiting in 10 seconds", sdl.Color{R: 255, G: 0, B: 0, A: 255})
		if err != nil {
			panic(err)
		}
		exit(10, exitChan)
	}

	//support missing file extension
	if !strings.HasSuffix(os.Args[1], ".json") {
		os.Args[1] += ".json"
	}

	//check if os.Args[1] file exists
	var resp *http.Response
	file, err := os.Open(os.Args[1])
	if err != nil {
		//open a window for error message
		window, err := newWindow("Clide", clide.Position{
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
		typer := clide.Typer{
			Window: window,
			Pos: clide.Position{
				X: 5,
				Y: 5,
				H: 0,
				W: 0,
			},
			Font: clide.Font{
				Path: fontPath,
				Size: fontSize,
			},
			Speed:    100,
			Humanize: 0.9,
		}

		errorText := fmt.Sprintf("File %s does not exists in current directory, checking /usr/share/clide", os.Args[1])
		typer.Pos, err = clide.Print(typer, errorText, sdl.Color{R: 255, G: 0, B: 0, A: 255})
		if err != nil {
			panic(err)
		}

		//if not check usr/share/clide
		file, err = os.Open("/usr/share/clide/" + os.Args[1])
		if err != nil {
			errorText = fmt.Sprintf("File %s does not exists /usr/share/clide. Checking for clide examples on clide.sh with name %s ...", os.Args[1], os.Args[1])
			typer.Pos, err = clide.Print(typer, errorText, sdl.Color{R: 255, G: 0, B: 0, A: 255})
			if err != nil {
				panic(err)
			}

			//if not finally check clide demo fileserver
			resp, err = http.Get("https://clide.sh/demos/" + os.Args[1])
			if err != nil || resp.StatusCode != 200 {
				typer.Pos, err = clide.Print(typer, "Could not find file at clide.sh/demos", sdl.Color{R: 255, G: 0, B: 0, A: 255})
				if err != nil {
					panic(err)
				}
				exit(10, exitChan)
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

	cfg, err = cfg.Validate()
	if err != nil {
		panic(err)
	}

	//open a window for each defined in json
	typerList := []clide.Typer{}
	for _, w := range cfg.Windows {
		window, err := newWindow(w.Name, clide.Position{
			X: w.X,
			Y: w.Y,
			H: w.Height,
			W: w.Width,
		})
		if err != nil {
			panic(err)
		}
		defer window.Destroy()

		//initialize typer values
		typer := clide.Typer{
			Window: window,
			Pos: clide.Position{
				X: 5,
				Y: 5,
				H: 0,
				W: 0,
			},
			Font: clide.Font{
				Path: fontPath,
				Size: fontSize,
			},
			Speed:    cfg.TypeSpeed,
			Humanize: cfg.Humanize,
		}

		//reset window background color
		bgColor, err := clide.StringToColor(cfg.ColorScheme.TerminalBG)
		if err != nil {
			panic(err)
		}
		typer, err = clide.ClearWindow(typer, bgColor)
		if err != nil {
			panic(err)
		}

		typerList = append(typerList, typer)
	}

	//listen for quit events to close program
	go listenForQuit(exitChan)

	//run each command in the commands slice
	for _, cmd := range cfg.Commands {
		//add async commands to goroutine count
		if cmd.Async {
			goroutineCount++
		}

		//find the typer for the window specified in cmd
		index := 0
		if cmd.Window != "" {
			for i, v := range typerList {
				if v.Window.GetTitle() == cmd.Window {
					index = i
				}
			}
		}

		err := cmd.Validate()
		if err != nil {
			clide.Print(typerList[index], err.Error(), sdl.Color{R: 255, G: 0, B: 0, A: 255})
		}

		if cmd.IsInstalled() {
			typerList[index], err = cmd.Run(cfg, typerList[index], exitChan)
			if err != nil {
				panic(err)
			}
		} else {
			if !cfg.HideWarnings {
				warning := fmt.Sprintf("WARNING: %s is not installed! Skipping command: '%s'.\n", strings.Split(cmd.CmdString, " ")[0], cmd.CmdString)
				typerList[index].Pos, err = clide.Print(typerList[index], warning, sdl.Color{R: 255, G: 0, B: 0, A: 255})
				if err != nil {
					panic(err)
				}
			}
		}
	}
	exit(1, exitChan)
}

// newWindow creates a new sdl2 window
func newWindow(title string, pos clide.Position) (*sdl.Window, error) {
	var window *sdl.Window
	var err error

	// Create a window for us to draw the text on
	if window, err = sdl.CreateWindow(title, pos.X, pos.Y, pos.W, pos.H, sdl.WINDOW_SHOWN); err != nil {
		return nil, err
	}

	iconSurface, err := sdl.LoadBMP("./assets/clide_icon.bmp")
	if err != nil {
		return nil, err
	}
	window.SetIcon(iconSurface)

	return window, nil
}

// listenForQuit watches for a quit event on any window and exits clide with status 1 when found
func listenForQuit(exitChan chan bool) {
	for {
		//keep checking keyboard events until a trigger key is pressed
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch target := event.(type) {

			//if quit event, close program
			case *sdl.QuitEvent:
				exit(1, exitChan)
			//if any window is closed, close program
			case *sdl.WindowEvent:
				if target.Event == sdl.WINDOWEVENT_CLOSE {
					exit(1, exitChan)
				}
				//keyboard keys to quit
				// case *sdl.KeyboardEvent:
				// 	if target.Keysym.Sym == sdl.K_ESCAPE {
				// 		exit(1, exitChan)
				// 	}
			}
		}
	}
}

// exit stops the program after delaying for x seconds
func exit(delay int, exit chan bool) {
	//send exit signal to all async commands
	for goroutineCount > 0 {
		exit <- true
		goroutineCount--
	}

	//used to give time for goroutines to exit
	time.Sleep(time.Second * time.Duration(delay))

	os.Exit(1)
}
