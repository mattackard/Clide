package clide

import (
	"math/rand"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Typer holds all information neeeded to type or print text to a window
type Typer struct {
	Window   *sdl.Window
	Pos      Position
	Font     Font
	Speed    int
	Humanize float64
}

// Position holds the positional data of a sdl2 surface
type Position struct {
	X int32
	Y int32
	H int32
	W int32
}

// Font holds font information for printing text
type Font struct {
	Path string
	Size int
}

var textColor = sdl.Color{R: 220, G: 220, B: 220, A: 255}

// Type prints the text to a sdl2 window and simulates a user typing the string
// it returns a position struct describing the completed surface
func Type(typer Typer, text string) (Position, error) {
	//split string into array of cingle characters
	split := strings.Split(text, "")

	var surface *sdl.Surface
	var font *ttf.Font
	var textsurface *sdl.Surface

	var err error
	var lastX int32
	var lastY int32
	lastX = typer.Pos.X
	lastY = typer.Pos.Y

	//get surface info
	if surface, err = typer.Window.GetSurface(); err != nil {
		return Position{}, err
	}

	// Load the font for our text
	if font, err = ttf.OpenFont(typer.Font.Path, typer.Font.Size); err != nil {
		return Position{}, err
	}
	defer font.Close()

	for _, char := range split {
		// Wait to simulate typing speed
		time.Sleep(getKeyDelay(typer.Speed, typer.Humanize))

		// Load the font for our text
		if font, err = ttf.OpenFont(typer.Font.Path, typer.Font.Size); err != nil {
			return Position{}, err
		}
		defer font.Close()

		// Create text using font
		if textsurface, err = font.RenderUTF8Blended(char, textColor); err != nil {
			return Position{}, err
		}
		defer textsurface.Free()

		// Wrap text if it's too long
		if lastX > surface.W-15 || []byte(char)[0] == []byte("\n")[0] {
			lastX = typer.Pos.X
			lastY += textsurface.H + 2
		}

		// Dont print newline characters
		if []byte(char)[0] != []byte("\n")[0] {
			if err := textsurface.Blit(nil, surface, &sdl.Rect{X: lastX, Y: lastY, W: 0, H: 0}); err != nil {
				return Position{}, err
			}

			lastX += textsurface.W

			// Update the window surface with what we have drawn
			typer.Window.UpdateSurface()
		}
	}

	lastY += textsurface.H + 2

	//return position for next line to be typed from
	return Position{
		X: lastX,
		Y: lastY,
		W: surface.W,
		H: surface.H,
	}, nil
}

// Print prints text to the sdl2 window all at once
func Print(typer Typer, text string) (Position, error) {
	split := strings.Split(text, "\n")

	var surface *sdl.Surface
	var font *ttf.Font
	var textsurface *sdl.Surface

	var err error
	var lastY int32
	lastY = typer.Pos.Y

	if surface, err = typer.Window.GetSurface(); err != nil {
		return Position{}, err
	}

	// Load the font for our text
	if font, err = ttf.OpenFont(typer.Font.Path, typer.Font.Size); err != nil {
		return Position{}, err
	}
	defer font.Close()

	//print each line individually
	for _, line := range split {
		if len(line) > 0 {
			// Create text using font
			if textsurface, err = font.RenderUTF8Blended(line, textColor); err != nil {
				return Position{}, err
			}
			defer textsurface.Free()

			if err := textsurface.Blit(nil, surface, &sdl.Rect{X: typer.Pos.X, Y: lastY, W: 0, H: 0}); err != nil {
				return Position{}, err
			}
			lastY += textsurface.H + 2
		}
	}

	// Update the window surface with what we have drawn
	err = typer.Window.UpdateSurface()
	if err != nil {
		panic(err)
	}

	//in case textSurface was never defined
	if textsurface == nil {
		//return position for next line to be typed from
		return Position{
			X: 5,
			Y: lastY,
			W: surface.W,
			H: surface.H,
		}, nil
	}

	//return position for next line to be typed from
	return Position{
		X: textsurface.W + 7,
		Y: lastY,
		W: surface.W,
		H: surface.H,
	}, nil
}

// ListenForKey blocks execution until a key is pressed.
// Use in a goroutine to watch in the background
func ListenForKey(cfg Config) {
	pressed := false
	for !pressed {
		//keep checking keyboard events until a trigger key is pressed
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				for _, key := range cfg.TiggerKeys {
					if t.Keysym.Sym == sdl.GetKeyFromName(key) {
						pressed = true
					}
				}
			}
		}
	}
}

//getKeyDelay calculates and returns a time to wait based on type speed and humanization ratio
func getKeyDelay(typeSpeed int, humanize float64) time.Duration {
	if humanize > 0 {
		//set up a seeded random
		rand.Seed(time.Now().UnixNano())

		//calculate speed variance based on humanize field
		variance := (1 - humanize - rand.Float64()) * float64(typeSpeed)

		return time.Duration(float64(typeSpeed)+variance) * time.Millisecond
	}
	return time.Duration(typeSpeed) * time.Millisecond
}
