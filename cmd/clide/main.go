package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	helpText     = `Clide CLI Usage:
		clide example.json		runs the clide demo stored in example.json
		clide-editor			opens the clide demo editor GUI interface				
		clide				shows this help message`
)

var goroutineCount int

func main() {
	// initialize sdl2
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer ttf.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// create a channel to exit out of goroutines when program exits
	exitChan := make(chan bool, goroutineMax)

	// check if os.Args[1] exists
	if len(os.Args) < 2 {
		cfg, err := clide.NewDefaultConfig()
		if err != nil {
			panic(err)
		}
		noFileError(cfg)
		exit(10, exitChan)
	}

	// support missing file extension
	if !strings.HasSuffix(os.Args[1], ".json") {
		os.Args[1] += ".json"
	}

	// check if os.Args[1] file exists
	var resp *http.Response
	file, err := os.Open(os.Args[1])
	if err != nil {
		cfg, err := clide.NewDefaultConfig()
		if err != nil {
			panic(err)
		}
		file, err = checkAlternateFileLocations(cfg, exitChan)
		if err != nil {
			panic(err)
		}
	}

	// unmarshal clide json into config struct
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

	// open a window for each defined in json
	typerList, err := cfg.BuildTyperList()
	if err != nil {
		panic(err)
	}

	// reset window background color for each window
	for _, typer := range typerList {
		bgColor, err := clide.StringToColor(cfg.ColorScheme.TerminalBG)
		if err != nil {
			panic(err)
		}
		err = typer.ClearWindow(bgColor)
		if err != nil {
			panic(err)
		}
	}

	cfg.TyperList = typerList

	// listen for quit events to close program
	go listenForQuit(exitChan)

	// run each command in the commands slice
	for _, cmd := range cfg.Commands {
		// add async commands to goroutine count
		if cmd.Async {
			goroutineCount++
		}

		// find the typer for the window specified in cmd
		index := 0
		if cmd.Window != "" {
			for i, v := range cfg.TyperList {
				if v.Window.GetTitle() == cmd.Window {
					index = i
				}
			}
		}

		err = cmd.Run(&cfg, typerList[index], exitChan)
		if err != nil {
			panic(err)
		}
	}
	exit(1, exitChan)
}

// noFileError createse a window containing an error message when no file is passed to clide
func noFileError(cfg clide.Config) {
	// open a window for error message
	window, err := clide.NewWindow("Clide", clide.Position{
		X: 0,
		Y: 0,
		H: 800,
		W: 1000,
	})
	if err != nil {
		panic(err)
	}

	// initialize typer values
	typer := cfg.NewTyper(window)

	fmt.Println("You must provide a clide configured json file to run a demo.")
	err = typer.Print("You must provide a clide configured json file to run a demo.", sdl.Color{R: 255, G: 100, B: 100, A: 255})
	if err != nil {
		panic(err)
	}
	fmt.Println("\n" + helpText)
	typer.Pos.X = 20
	typer.Pos.Y += 20
	err = typer.Print(helpText, sdl.Color{R: 255, G: 100, B: 100, A: 255})
	if err != nil {
		panic(err)
	}
	typer.Pos.X = 5
	typer.Pos.Y += 20
	err = typer.Print("Exiting in 10 seconds", sdl.Color{R: 255, G: 100, B: 100, A: 255})
	if err != nil {
		panic(err)
	}
}

// checkAlernateFileLocations checks for the given clide demo file in built-in demo locations and returns an new
// path string if the file is found
func checkAlternateFileLocations(cfg clide.Config, exitChan chan bool) (*os.File, error) {
	var file *os.File

	// open a window for error message
	window, err := clide.NewWindow("Clide", clide.Position{
		X: 0,
		Y: 0,
		H: 800,
		W: 1000,
	})
	if err != nil {
		return nil, err
	}

	// initialize typer values
	typer := cfg.NewTyper(window)

	errorText := fmt.Sprintf("File %s does not exists in current directory, checking /usr/share/clide/examples/", os.Args[1])
	log.Println(errorText)
	err = typer.Print(errorText, sdl.Color{R: 255, G: 100, B: 100, A: 255})
	if err != nil {
		return nil, err
	}

	// if not check usr/share/clide
	file, err = os.Open("/usr/share/clide/examples/" + os.Args[1])
	if err != nil {
		errorText = fmt.Sprintf("File %s does not exists /usr/share/clide/examples/. Checking for clide examples on the web with name %s ...", os.Args[1], os.Args[1])
		log.Println(errorText)
		typer.Pos.X = 5
		err = typer.Print(errorText, sdl.Color{R: 255, G: 100, B: 100, A: 255})
		if err != nil {
			return nil, err
		}
	}

	// if demo is found in an alternate location, destroy error window
	window.Destroy()
	return file, nil
}

// listenForQuit watches for a quit event on any window and exits clide with status 1 when found
func listenForQuit(exitChan chan bool) {
	for {
		// keep checking keyboard events until a trigger key is pressed
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch target := event.(type) {

			// if quit event, close program
			case *sdl.QuitEvent:
				exit(1, exitChan)
			// if any window is closed, close program
			case *sdl.WindowEvent:
				if target.Event == sdl.WINDOWEVENT_CLOSE {
					exit(1, exitChan)
				}
			}
		}
	}
}

// exit stops the program after delaying for x seconds
func exit(delay int, exit chan bool) {
	// send exit signal to all async commands
	for goroutineCount > 0 {
		exit <- true
		goroutineCount--
	}

	// used to give time for goroutines to exit
	time.Sleep(time.Second * time.Duration(delay))

	os.Exit(0)
}
