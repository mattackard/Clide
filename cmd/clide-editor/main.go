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
	w.Navigate("file:///home/kubuntu/go/src/github.com/mattackard/Clide/cmd/clide-editor/convert.html")

	http.HandleFunc("/getFiles", getFiles)
	http.HandleFunc("/save", saveFiles)
	http.HandleFunc("/convert", convertToClide)
	http.HandleFunc("/run", runDemo)

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
