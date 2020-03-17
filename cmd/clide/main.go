package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/mattackard/Clide/pkg/clide"
)

func main() {
	//open default terminal instance
	// terminal1 := exec.Command("konsole")
	// terminal1.Stdin = os.Stdin
	// terminal1.Stdout = os.Stdout
	// err := terminal1.Start()
	// defer terminal1.Process.Signal(os.Interrupt)
	// if err != nil {
	// 	panic(err)
	// }

	//check if os.Args[1] exists
	if len(os.Args) < 2 {
		color.Println("<comment>You must provide a clide configured json file to run a demo.</>")
		showHelp()
		os.Exit(1)
	}

	//support missing file extension
	if !strings.HasSuffix(os.Args[1], ".json") {
		os.Args[1] += ".json"
	}

	//check if os.Args[1] file exists
	var resp *http.Response
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("File %s does not exists in current directory, checking /usr/share/clide \n", os.Args[1])

		//if not check usr/share/clide
		file, err = os.Open("/usr/share/clide/" + os.Args[1])
		if err != nil {
			fmt.Printf("File %s does not exists /usr/share/clide. Checking for clide examples on clide.sh with name %s ... \n", os.Args[1], os.Args[1])

			//if not finally check clide demo fileserver
			resp, err = http.Get("https://clide.sh/demos/" + os.Args[1])
			if err != nil || resp.StatusCode != 200 {
				panic("Could not find file " + os.Args[1])
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

	//run each command in the commands slice
	for _, cmd := range cfg.Commands {
		cmd.Validate()
		if cmd.IsInstalled() {
			err = cmd.Run(cfg)
			if err != nil {
				panic(err)
			}
		} else {
			if !cfg.HideWarnings {
				color.Printf("<comment>WARNING</>: %s is not installed! Skipping command: '%s'.\n", strings.Split(cmd.CmdString, " ")[0], cmd.CmdString)
			}
		}
	}
}

func showHelp() {
	helpText := `
	Clide CLI Usage:
		clide example.json		runs the clide demo stored in example.json
		clide-sh script.sh		converts script.sh into script.json formatted as a clide demo
		clide-build			opens the clide demo builder interface				
		clide				shows this help message
	`
	fmt.Println(helpText)
}
