package clide

import (
	"math/rand"
	"strings"
	"sync"
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
	mutex    sync.Mutex
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

// Type prints the text to a sdl2 window and simulates a user typing the string
// it returns a position struct describing the completed surface
func (typer *Typer) Type(text string, color sdl.Color) error {
	// split string into array of cingle characters
	split := strings.Split(text, "")

	var surface *sdl.Surface
	var font *ttf.Font
	var textSurface *sdl.Surface

	var err error
	lastX := typer.Pos.X
	lastY := typer.Pos.Y

	// get surface info
	if surface, err = typer.Window.GetSurface(); err != nil {
		return err
	}

	// Load the font for our text
	if font, err = ttf.OpenFont(typer.Font.Path, typer.Font.Size); err != nil {
		return err
	}
	defer font.Close()

	for _, char := range split {
		// Wait to simulate typing speed
		time.Sleep(getKeyDelay(typer.Speed, typer.Humanize))

		// Load the font for our text
		if font, err = ttf.OpenFont(typer.Font.Path, typer.Font.Size); err != nil {
			return err
		}
		defer font.Close()

		// Create text using font
		if textSurface, err = font.RenderUTF8Blended(char, color); err != nil {
			return err
		}
		defer textSurface.Free()

		// Wrap text if it's too long
		if lastX > surface.W-15 || []byte(char)[0] == []byte("\n")[0] {
			lastX = 5
			if lastY+textSurface.H+2 > surface.H {
				scrollBottom(surface, textSurface)
				lastY -= textSurface.H + 2
			} else {
				lastY += textSurface.H + 2
			}
		}

		// Dont print newline characters
		if []byte(char)[0] != []byte("\n")[0] {
			if lastY+textSurface.H+2 > surface.H {
				scrollBottom(surface, textSurface)
				lastY -= textSurface.H + 2
			}

			err := textSurface.Blit(nil, surface, &sdl.Rect{X: lastX, Y: lastY, W: 0, H: 0})
			if err != nil {
				return err
			}

			lastX += textSurface.W

			// Update the window surface with what we have drawn
			typer.Window.UpdateSurface()
		}
	}
	lastY += textSurface.H + 2

	// return position for next line to be typed from
	typer.Pos = Position{
		X: lastX,
		Y: lastY,
		W: surface.W,
		H: surface.H,
	}
	return nil
}

// Print prints text to the sdl2 window all at once
func (typer *Typer) Print(text string, color sdl.Color) error {
	split := strings.Split(text, "\n")

	var surface *sdl.Surface
	var font *ttf.Font
	var textSurface *sdl.Surface

	var err error
	var lastY int32
	lastY = typer.Pos.Y

	if surface, err = typer.Window.GetSurface(); err != nil {
		return err
	}

	// Load the font for our text
	if font, err = ttf.OpenFont(typer.Font.Path, typer.Font.Size); err != nil {
		return err
	}
	defer font.Close()

	// print each line individually
	for _, line := range split {
		if len(line) > 0 {
			// Create text using font
			if textSurface, err = font.RenderUTF8Blended(line, color); err != nil {
				return err
			}
			defer textSurface.Free()

			// check if its planning to print to a newline here?
			if lastY+textSurface.H+2 > surface.H {
				scrollBottom(surface, textSurface)
				lastY -= textSurface.H + 2
			}

			err := textSurface.Blit(nil, surface, &sdl.Rect{X: typer.Pos.X, Y: lastY, W: 0, H: 0})
			if err != nil {
				return err
			}

			lastY += textSurface.H + 2

		}
	}

	// Update the window surface with what we have drawn
	err = typer.Window.UpdateSurface()
	if err != nil {
		return err
	}

	// in case textSurface was never defined
	if textSurface == nil {
		// return position for next line to be typed from
		typer.Pos = Position{
			X: 5,
			Y: lastY,
			W: surface.W,
			H: surface.H,
		}
		return nil
	}

	// return position for next line to be typed from
	typer.Pos = Position{
		X: typer.Pos.X + textSurface.W + 1,
		Y: lastY,
		W: surface.W,
		H: surface.H,
	}
	return nil
}

// BuildTyperList creates a window for each window defined in config, and then
// attaches a typer to each window that is created
func (cfg Config) BuildTyperList() ([]*Typer, error) {
	typerList := []*Typer{}
	for i, win := range cfg.Windows {
		window, err := NewWindow(win.Name, Position{
			X: win.X,
			Y: win.Y,
			H: win.Height,
			W: win.Width,
		})
		if err != nil {
			return nil, err
		}

		// set the window object in the cfg window
		cfg.Windows[i].Window = window

		// initialize typer values
		typer := cfg.NewTyper(window)

		typerList = append(typerList, typer)

		if cfg.HideWindows {
			window.Hide()
		}
	}
	return typerList, nil
}

// ListenForKey blocks execution until a key is pressed.
// Use in a goroutine to watch in the background
func ListenForKey(cfg Config) {
	pressed := false
	for !pressed {
		// keep checking keyboard events until a trigger key is pressed
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

// NewWindow creates a new clide window using sdl2
func NewWindow(title string, pos Position) (*sdl.Window, error) {
	var window *sdl.Window
	var err error

	// Create a window for us to draw the text on
	if window, err = sdl.CreateWindow(title, pos.X, pos.Y, pos.W, pos.H, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE); err != nil {
		return nil, err
	}

	iconSurface, err := sdl.LoadBMP("/usr/share/clide/assets/clide_icon.bmp")
	if err != nil {
		return nil, err
	}
	window.SetIcon(iconSurface)

	return window, nil
}

// NewTyper creates a new typer referencing the given window
func (cfg Config) NewTyper(window *sdl.Window) *Typer {
	typer := Typer{
		Window: window,
		Pos: Position{
			X: 5,
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
	return &typer
}

// getKeyDelay calculates and returns a time to wait based on type speed and humanization ratio
func getKeyDelay(typeSpeed int, humanize float64) time.Duration {
	if humanize > 0 {
		// set up a seeded random
		rand.Seed(time.Now().UnixNano())

		// calculate speed variance based on humanize field
		variance := (1 - humanize - rand.Float64()) * float64(typeSpeed)

		return time.Duration(float64(typeSpeed)+variance) * time.Millisecond
	}
	return time.Duration(typeSpeed) * time.Millisecond
}

// scrollBottom moves the surface to make room for another line if there is no more room at the bottom of the window
func scrollBottom(surface *sdl.Surface, lastText *sdl.Surface) error {
	// store the data in the current surface
	surfaceStore := *surface

	r, g, b, a := surfaceStore.At(0, 0).RGBA()

	colorFix := sdl.Color{
		R: uint8(a),
		G: uint8(b),
		B: uint8(g),
		A: uint8(r),
	}

	// move the surface up by one line
	err := surfaceStore.Blit(nil, surface, &sdl.Rect{X: 0, Y: -lastText.H - 2, W: 0, H: 0})
	if err != nil {
		return err
	}

	// erase the old line that was moved up one
	err = surface.FillRect(&sdl.Rect{X: 0, Y: surface.H - lastText.H - 2, W: surface.W, H: lastText.H}, colorFix.Uint32())
	if err != nil {
		return err
	}
	return nil
}
