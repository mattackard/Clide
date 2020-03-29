package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mattackard/Clide/pkg/clide"

	"github.com/zserge/webview"
)

// Request holds data received from the GUI
type Request struct {
	Filename     string `json:"filename"`
	FileContents string `json:"fileContents"`
}

// open a webview connection using the client.html file
// css and js is loaded in through the html
func main() {
	//set up webview gui
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Clide Editor")
	w.SetSize(1920, 1080, webview.HintNone)
	w.Navigate("file:///home/ubuntu/go/src/github.com/mattackard/Clide/cmd/clide-editor/convert.html")

	http.HandleFunc("/convert", convertToClide)

	// start http server for communiaction with gui js
	go http.ListenAndServe(":8080", nil)

	// launch webview gui
	w.Run()
}

// convertToClide takes the file sent and converts it into a clide demo
func convertToClide(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got a request from ", r.RemoteAddr)
	body, err := getReqBody(r)
	httpError(w, err, http.StatusInternalServerError)

	contentType := getFileType(body)

	switch contentType {
	case "json":
		w = setHeaders(w)
		w.Write([]byte(body.FileContents))
	case "script":
		clide, err := buildClide(body.FileContents)
		httpError(w, err, http.StatusInternalServerError)

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
		w = setHeaders(w)
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
	}
}
