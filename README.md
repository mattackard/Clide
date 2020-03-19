# Clide
[![Release](https://img.shields.io/github/v/release/mattackard/Clide)](https://github.com/mattackard/Clide/releases)
[![Build](https://img.shields.io/github/workflow/status/mattackard/Clide/Clide)](https://github.com/mattackard/Clide/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattackard/Clide)](https://goreportcard.com/report/github.com/mattackard/Clide)

Clide is a tool for creating automated demos for terminal-based applications. Using a json file, you can script out a demo in the terminal to free your hands up from having to copy and paste commands, or handle typos when you're in the middle of giving a presentation.

## Getting Started

Clide is not ready for a full release at this time, so the best way to get started using it is to clone the repository. From there you can run the go file in `cmd/clide` directly, or build it and move the binary into your bin path.

## Prerequisites

Clide uses the [go-sdl2](https://github.com/veandco/go-sdl2) project to create it's windows and track program events. Go-sdl2 requires SDL to be installed to you pkg-config. Check the [installation requirements](https://github.com/veandco/go-sdl2#requirements) in sdl2's readme for the most up-to-date guide for installation.

The Clide source contains an `examples` folder that contains sample json demos. It is recommended to try out running clide with one of the included demo json files to confirm all dependencies are installed.

**Running a demo from a built binary:**

`cmd/clide/*.go examples/demo.json`

**Running a demo from a built binary:**

`clide examples/demo.json`

## Development

To get started working with clide yourself, clone or fork the repository. This project follows the common [Go Standard Project Layout](https://github.com/golang-standards/project-layout) for organization of program files. For this reason, you can find the main application inside the cmd folder and any supporting packages inside the pkg folder. 

The included `Makefile` has small helper scripts to run and build the project quickly. You can run these by navigating to the root project directory and running make along with the script you want to execute.

```
make builds
make run-demo
```

The demo files in the `examples` folder is a good place to start tweaking the setting and getting familiar with everything clide can do.

## Built With

* [Go](http://golang.org) - The Go programming language
* [go-sdl2](github.com/veandco/go-sdl2) - Clide terminal and window management

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/mattackard/Clide/tags). 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details