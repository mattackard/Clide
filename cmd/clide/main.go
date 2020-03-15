package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/mattackard/Clide/pkg/clide"
)

func main() {
	//support missing file extension
	if !strings.HasSuffix(os.Args[1], ".json") {
		os.Args[1] += ".json"
	}

	//unmarshal clide json into config struct
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(file)
	cfg := clide.Config{}
	err = json.Unmarshal(bytes, &cfg)

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
				color.Printf("<yellow>WARNING</>: %s is not installed! Skipping command: '%s'.\n", strings.Split(cmd.CmdString, " ")[0], cmd.CmdString)
			}
		}
	}
}
