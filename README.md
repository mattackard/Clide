# Clide

Clide is a tool for creating automated demos for terminal-based applications. Using a json file, you can script out a demo in the terminal to free your hands up from having to copy and paste commands, or handle typos when you're in the middle of giving a presentation.

## Getting Started

To get started using Clide, you can clone the repository and run the program using Go directly, or you can download a binary from the release tab and get started without having to install go at all.

### Prerequisites

In order to use this application you must first have a json file configured for your demo. You can check out the clide demo by running clide with the demo.json in the `examples` folder.

```
clide ./examples/demo.json
```

### Development

To get started working with clide yourself, just clone or fork the repository. This project follows the common [Go Standard Project Layout](https://github.com/golang-standards/project-layout) for organization of program files. For this reason, you can find the main application inside the cmd folder and any supporting packages inside the pkg folder.

The included `Makefile` has small helper scripts to run and build the project quickly. You can run these by navigating to the root project directory and running make along with the script you want to execute.

```
make builds
make demo-run
```

The demo file in the `examples` folder is a good place to start tweaking the setting and getting familiar with everything clide can do.

## Built With

* [Go](http://golang.org) - The Go programming language
* [gookit/color](https://github.com/gookit/color) - Terminal text coloring

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/mattackard/Clide/tags). 

## Authors

* **Matt Ackard** - *Initial work* - [Github](https://github.com/mattackard)

See also the list of [contributors](https://github.com/mattackard/Clide/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details