package clide

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
)

//Config holds the global configuration for a clide demo
type Config struct {
	User           string    `json:"user"`
	Directory      string    `json:"directory"`
	TypeSpeed      int       `json:"typeSpeed"`
	Humanize       float64   `json:"humanize"`
	HideWarnings   bool      `json:"hideWarnings"`
	ClearBeforeAll bool      `json:"clearBeforeAll"`
	TiggerKeys     []string  `json:"triggerKeys"`
	Commands       []Command `json:"commands"`
}

//Validate checks for potential issues in a Config and
//adds some default values if they are not present
func (cfg Config) Validate() {
	//throw error when no commands are present
	if len(cfg.Commands) == 0 {
		panic("No commands found in provided json file")
	}

	//default directory
	if cfg.Directory == "" {
		var err error
		cfg.Directory, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	//default user
	if cfg.User == "" {
		cfg.User = "demo@clide"
	}

	if !cfg.HideWarnings {
		//check if all commands are installed
		notInstalled := []string{}
		for _, v := range cfg.Commands {
			if !v.IsInstalled() {
				notInstalled = append(notInstalled, v.CmdString)
			}
		}

		//comfirm user wants to run program even though uninstalled commands will be skipped
		if len(notInstalled) != 0 {
			color.Println("<yellow>WARNING</>: At least one command is not installed on the system! The following commands will be skipped:")
			for _, badCmd := range notInstalled {
				split := strings.Split(badCmd, " ")
				joined := strings.Join(split[1:], " ")
				color.Printf("\t<green>%s</> %s \n", split[0], joined)
			}
			fmt.Print("Would you still like to run the demo? [y/n]: ")

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answer := scanner.Text()
			fmt.Scan(answer)
			if !strings.HasPrefix(answer, "y") && !strings.HasPrefix(answer, "Y") {
				fmt.Println("Exiting program")
				os.Exit(1)
			}
		}
	}
}
