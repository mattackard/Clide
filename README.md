# Clide
[![Release](https://img.shields.io/github/v/release/mattackard/Clide)](https://github.com/mattackard/Clide/releases)
[![Build](https://img.shields.io/github/workflow/status/mattackard/Clide/Clide)](https://github.com/mattackard/Clide/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattackard/Clide)](https://goreportcard.com/report/github.com/mattackard/Clide)

Clide is a tool for creating automated demos for terminal-based applications. Using a json file, you can script out a demo in the terminal to free your hands up from having to copy and paste commands, or handle typos when you're in the middle of giving a presentation. The scripted commands run in real-time and display the output in a window to emulate a terminal.

## Getting Started

To try out clide, go to the [github release page](https://github.com/mattackard/Clide/releases) and download the package for your OS (currently only debian-based linux is supported).

An alternative to using the binary is to clone the github project and run clide from the source. More information about how to run the program from the source is in the sections below.

## Installation

The github release file is a debian package that can be installed with the built-in software manager available in debian operating systems. To install the .deb file from within the terminal run:

`sudo dpkg -i clide_1.1.0.deb`

**Running a demo from a built binary:**

`clide demo`
`clide windows`
`clide network`

The demos included with the binary are `demo`, `windows`, and `network`. These json files can be found in the `/usr/share/clide/examples` folder for reference. Looking at the included demo json files is a good starting point to see all the options available when creating a demo with clide. Note that some included demos may use linux packages not installed on your system. If you see an error about a command not running because a program is not installed, you must install the package for the demo to execute properly.

## Development

Clide uses the [go-sdl2](https://github.com/veandco/go-sdl2) project to create it's windows and track program events. Go-sdl2 requires SDL to be installed to your pkg-config. Check the [installation requirements](https://github.com/veandco/go-sdl2#requirements) in sdl2's readme for the most up-to-date guide for installation.

To get started working with clide yourself, clone or fork the repository. This project follows the common [Go Standard Project Layout](https://github.com/golang-standards/project-layout) for organization of program files. For this reason, you can find the main application inside the cmd folder and any supporting packages inside the pkg folder. 

The Clide source contains an `examples` folder that contains sample json demos. It is recommended to try out running clide with one of the included demo json files to confirm all dependencies are installed.

The included `Makefile` has small helper scripts to run and build the project quickly. You can run these by navigating to the root project directory and running make along with the script you want to execute.

`make set-env` is required for running the program from source and sets up the environment with the asset files clide needs to run properly. 

```
# removes all clide configuration files from /usr filesystem
make rm-env

# builds the go program into a binary
make builds

# runs demo.json localated in the examples folder in project directory
make run-demo
```

The demo files in the `examples` folder is a good place to start tweaking the settings to get familiar with everything clide can do.

## Built With

* [Go](http://golang.org) - The Go programming language
* [go-sdl2](github.com/veandco/go-sdl2) - Clide terminal and window management

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/mattackard/Clide/tags). 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
