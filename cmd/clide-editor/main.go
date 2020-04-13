package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/mattackard/Clide/pkg/clide"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/zserge/webview"
)

// Request holds data received from the GUI
type Request struct {
	Filename     string `json:"filename"`
	FileContents string `json:"fileContents"`
}

//Files stores the information for the files being worked on in the editor
type Files struct {
	ScriptText string `json:"scriptText"`
	JSONText   string `json:"jsonText"`
}

// Initialize a files struct for persistent storage
var files = Files{}

var (
	fontPath = "/usr/share/clide/assets/UbuntuMono-B.ttf"
	fontSize = 18
)

// open a webview connection using the client.html file
// css and js is loaded in through the html
func main() {
	//get file contents passed in as an argumemnt
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}

		contents, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		if strings.HasSuffix(os.Args[1], ".json") {
			files.JSONText = string(contents)
		} else if strings.HasSuffix(os.Args[1], ".sh") {
			files.ScriptText = string(contents)

			bytes, err := buildClide(files.ScriptText)
			if err != nil {
				log.Println("Could not convert contents of given file to clide configuration")
			}

			files.JSONText = string(bytes)
		}
	}

	//set up webview gui
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Clide Editor")
	w.SetSize(1920, 1080, webview.HintNone)
	// w.Navigate("file:///usr/share/clide/editor/edit.html")
	w.Navigate("file:///home/xubuntu/go/src/github.com/mattackard/Clide/cmd/clide-editor/edit.html")

	http.HandleFunc("/getFiles", getFiles)
	http.HandleFunc("/save", saveFiles)
	http.HandleFunc("/convert", convertToClide)
	http.HandleFunc("/run", runDemo)
	http.HandleFunc("/arrangeWindows", arrangeWindows)

	// start http server for communiaction with gui js
	go http.ListenAndServe(":8080", nil)

	// launch webview gui
	w.Run()
}

// runDemo test runs the json file provided in the request with clide
func runDemo(w http.ResponseWriter, r *http.Request) {
	//parse json
	body, err := getReqBody(r)
	httpError(w, err, http.StatusInternalServerError)

	//create a temp file and write the json into it
	file, err := os.Create("temp.json")
	httpError(w, err, http.StatusInternalServerError)

	_, err = file.WriteString(body.FileContents)
	httpError(w, err, http.StatusInternalServerError)

	//update stored json contents
	files.JSONText = body.FileContents

	file.Close()
	defer os.Remove("temp.json")

	//execute clide with the temp json as an argument
	cmd := exec.Command("clide", "temp.json")
	err = cmd.Start()
	httpError(w, err, http.StatusInternalServerError)

	w = setHeaders(w)
	w.Write([]byte("OK"))

	//wait for clide execution to finish before deleting temp file
	cmd.Wait()
}

// getFiles sends the contents of the currently stored file contents
func getFiles(w http.ResponseWriter, r *http.Request) {
	w = setHeaders(w)

	bytes, err := json.Marshal(files)
	httpError(w, err, http.StatusInternalServerError)

	w.Write(bytes)
}

// saveFiles saves the contents of the files from the client editor to the files struct
func saveFiles(w http.ResponseWriter, r *http.Request) {
	//read request body
	body, err := ioutil.ReadAll(r.Body)
	httpError(w, err, http.StatusInternalServerError)

	//store contents of request body into global files struct
	err = json.Unmarshal(body, &files)
	httpError(w, err, http.StatusInternalServerError)

	w = setHeaders(w)
	w.Write([]byte("OK"))
}

// convertToClide takes the file sent and converts it into a clide demo
func convertToClide(w http.ResponseWriter, r *http.Request) {
	body, err := getReqBody(r)
	httpError(w, err, http.StatusInternalServerError)

	contentType := getFileType(body)

	switch contentType {
	case "json":
		//store json into files
		files.JSONText = body.FileContents

		//write the json back as response
		w = setHeaders(w)
		w.Write([]byte(body.FileContents))
	case "script":
		//store script into files
		files.ScriptText = body.FileContents

		clide, err := buildClide(body.FileContents)
		httpError(w, err, http.StatusInternalServerError)
		files.JSONText = string(clide)

		w = setHeaders(w)
		w.Write(clide)
	default:
		fmt.Println("unknown file type")
		err := errors.New("the file provided could not be recognized as a json or script file")
		httpError(w, err, http.StatusInternalServerError)
	}
}

func arrangeWindows(w http.ResponseWriter, r *http.Request) {
	//get request body
	body, err := getReqBody(r)
	httpError(w, err, http.StatusInternalServerError)

	//create a new config struct and unmarshal filecontents into it
	cfg := clide.Config{}

	err = json.Unmarshal([]byte(body.FileContents), &cfg)

	httpError(w, err, http.StatusInternalServerError)

	//initialize sdl2
	err = ttf.Init()
	httpError(w, err, http.StatusInternalServerError)
	defer ttf.Quit()

	err = sdl.Init(sdl.INIT_VIDEO)
	httpError(w, err, http.StatusInternalServerError)
	defer sdl.Quit()

	//open a window for each defined in json
	typerList := []*clide.Typer{}
	for i, win := range cfg.Windows {
		window, err := newWindow(win.Name, clide.Position{
			X: win.X,
			Y: win.Y,
			H: win.Height,
			W: win.Width,
		})
		httpError(w, err, http.StatusInternalServerError)

		// go listenForResize(window)

		//set the window object in the cfg window
		cfg.Windows[i].Window = window

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
		}

		typerList = append(typerList, &typer)
	}

	cfg.TyperList = typerList

	for _, typer := range typerList {
		err := typer.Print("Press enter to store window positions", sdl.Color{R: 255, G: 255, B: 255, A: 255})
		httpError(w, err, http.StatusInternalServerError)
	}

	//listen for quit events to close program
	listenForQuit()

	for i, newPos := range cfg.Windows {
		//store the new position
		newX, newY := newPos.Window.GetPosition()
		cfg.Windows[i].X = newX
		cfg.Windows[i].Y = newY

		//store the new size
		newWidth, newHeight := newPos.Window.GetSize()
		cfg.Windows[i].Width = newWidth
		cfg.Windows[i].Height = newHeight
	}

	bytes, err := json.Marshal(cfg)
	httpError(w, err, http.StatusInternalServerError)

	//write the json back as response
	w = setHeaders(w)
	w.Write(bytes)
}

func resizeHandler(window clide.Window, event sdl.Event) {

}

func buildClide(text string) ([]byte, error) {
	//initialize a command slice
	commands := []clide.Command{}

	//create a command struct for each line in the script
	split := strings.Split(text, "\n")
	for _, line := range split {

		//filter out comments and empty lines
		if !strings.HasPrefix(line, "#") && len(strings.Trim(line, " ")) > 0 {
			commands = append(commands, clide.Command{
				CmdString: strings.Trim(line, " "),
				PreDelay:  500,
				PostDelay: 500,
			})
		}
	}

	//create a config and put all commands in it
	cfg := clide.Config{
		User:      "demo@clide",
		Directory: "/",
		Commands:  commands,
	}

	bytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// determines the file's content type to see if conversion to json is needed
func getFileType(req Request) string {
	if (strings.HasSuffix(req.Filename, ".json") || strings.HasPrefix(req.FileContents, "{")) && json.Valid([]byte(req.FileContents)) {
		return "json"
	} else if strings.HasSuffix(req.Filename, ".sh") || strings.HasPrefix(req.FileContents, "#!") {
		return "script"
	}
	return "unknown"
}

// returns the http request body as a string
func getReqBody(r *http.Request) (Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Request{}, err
	}

	file := Request{}
	err = json.Unmarshal(body, &file)
	if err != nil {
		return Request{}, err
	}

	return file, nil
}

//set header to expect json and allow cors
func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	return w
}

// sends an error status code as a response
func httpError(w http.ResponseWriter, err error, statusCode int) {
	if err != nil {
		log.Println(err)
		w = setHeaders(w)
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
	}
}

// newWindow creates a new sdl2 window
func newWindow(title string, pos clide.Position) (*sdl.Window, error) {
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

// listenForQuit watches for a quit event on any window and exits clide with status 1 when found
func listenForQuit() {
	for {
		//keep checking keyboard events until a trigger key is pressed
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch target := event.(type) {

			//if quit event, close program
			case *sdl.QuitEvent:
				return
			//if any window is closed, close program
			case *sdl.WindowEvent:
				if target.Event == sdl.WINDOWEVENT_CLOSE {
					return
				}
			//keyboard keys to quit
			case *sdl.KeyboardEvent:
				if target.Keysym.Sym == sdl.K_KP_ENTER {
					return
				}
				if target.Keysym.Sym == sdl.K_RETURN {
					return
				}
			}
		}
	}
}

func listenForResize(window *sdl.Window) {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	for {
		//keep checking keyboard events until a trigger key is pressed
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch target := event.(type) {

			//if any window is closed, close program
			case *sdl.WindowEvent:
				if target.Event == sdl.WINDOWEVENT_RESIZED {
					fmt.Println("window resized")
					err := surface.FillRect(nil, sdl.Color{R: 255, G: 255, B: 255, A: 255}.Uint32())
					if err != nil {
						panic(err)
					}

					window.UpdateSurface()
				}
			}
		}
	}
}
