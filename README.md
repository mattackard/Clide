# Clide

[![Release](https://img.shields.io/github/v/release/mattackard/Clide)](https://github.com/mattackard/Clide/releases)
[![Build](https://img.shields.io/github/workflow/status/mattackard/Clide/Clide)](https://github.com/mattackard/Clide/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattackard/Clide)](https://goreportcard.com/report/github.com/mattackard/Clide)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://godoc.org/github.com/mattackard/Clide/pkg/clide)

Clide is a tool for creating automated demos for terminal-based applications. Using a json file, you can script out a demo in the terminal to free your hands up from having to copy and paste commands, or handle typos when you're in the middle of giving a presentation. The scripted commands run in real-time and display the output in a window to emulate a terminal.

## Getting Started

To try out clide, go to the [github release page](https://github.com/mattackard/Clide/releases) and download the package for your OS (currently only debian-based linux is supported).

An alternative to using the binary is to clone the github project and run clide from the source. More information about how to run the program from the source is in the sections below.

## Installation

The github release contains .deb and .rpm packages for installation on debain and rpm based linux distributions. To install clide, go to the [clide homepage](https://mattackard.github.io/Clide) and click the download for your operating system. Alternatively you can clone the repository or download the source code and use the makefile command `make install` to install clide from the source code. See below for dependency information if you want to install clide from the source code.

**Running a demo from a built binary:**

`clide demo`
`clide windows`
`clide network`

The demos included with the binary are `demo`, `windows`, and `network`. These json files can be found in the `/usr/share/clide/examples` folder for reference. Looking at the included demo json files is a good starting point to see all the options available when creating a demo with clide. Note that some included demos may use linux packages not installed on your system. If you see an error about a command not running because a program is not installed, you must install the package for the demo to execute properly.

**Opening a demo or script in the editor**

`clide-editor`
`clide-editor script.sh`
`clide-editor demo.json`

The editor provides a user-friendly method to edit and create clide demos. You can convert a script file into a clide demo using the `script convert` page. From there you can edit the details of the demo in the `modular editor` or edit the json directly in the script conversion page. You can test run the demos in each page to check on the accuracy of your demo at any time.

## Development

Clide uses the [go-sdl2](https://github.com/veandco/go-sdl2) project to create it's windows and track program events, and [webview](https://github.com/zserge/webview) for creating the editor GUI. Go-sdl2 requires SDL to be installed to your pkg-config. Check the [installation requirements](https://github.com/veandco/go-sdl2#requirements) in sdl2's readme for the most up-to-date guide for installation. Webview requires some graphical libraries that can also be found in the [github install notes](https://github.com/zserge/webview#notes) at the bottom of the readme.

To get started working with clide yourself, clone or fork the repository. This project follows the common [Go Standard Project Layout](https://github.com/golang-standards/project-layout) for organization of program files. For this reason, you can find the main application inside the cmd folder and any supporting packages inside the pkg folder.

The Clide source contains an `examples` folder that contains sample json demos. It is recommended to try out running clide with one of the included demo json files to confirm all dependencies are installed.

The included `Makefile` has small helper scripts to run and build the project quickly. You can run these by navigating to the root project directory and running make along with the script you want to execute.

`make install` is required for running the program from source and sets up the environment with the asset files clide needs to run properly. Look into the `Makefile` for more scripts to help in development.

The demo files in the `examples` folder is a good place to start looking at the settings to get familiar with everything clide can do.

## Built With

- [Go](http://golang.org) - The Go programming language
- [go-sdl2](https://github.com/veandco/go-sdl2) - Clide terminal and window management
- [webview](https://github.com/zserge/webview) - Creating the clide-editor GUI

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/mattackard/Clide/tags).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
