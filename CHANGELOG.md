# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2020-03-26

### Added

- Released as .deb package to auto-install dependencies
- Support for custom colors for text, background, and terminal prompt
- More error handling associated with defaults/missing json fields
- Icon added for clide terminal windows
- JSON documentation hosted on github-pages
- Clide website hosted on github-pages
- Clide package testing and integration with github actions

## [1.0.0] - 2020-03-22

### Added

- SDL2 implementation for creating and managing windows
- Ability to clear window before running a command
- Run commands triggered by a key press
- Commands can be run asynchronously
- Better error handling and validation of JSON demo files
- Clide checks all commands are installed before running
- Support for a bunch more JSON configuration fields:
    - hideWarnings - ignore warnings about uninstalled commands
    - clearBeforeAll - clears terminal before every command
    - keyTriggerAll - waits for a keypress to trigger all timings
    - windows - specifies window title, size, and position
    - triggerKeys - specifies the keys used to trigger command execution
    - window - specifies which window to run the given command in
    - clearBeforeRun - clear window before running the given command
    - waitForKey - requires key to trigger timings for given command
    - timeout - sets a timeout for the given command
    - hidden - runs the command hiding all output from the user
    - async - runs the command asynchronously, immediately continuing to the next


### Changed

- Clide now uses sdl2 windows instead of the os terminal window

## [0.0.2] - 2020-03-15

### Added

- Support for emulated typing in terminal
- Support for humanizing typed output

### Changed

- Update to example demo type speed and humanize ratio
- Pre and post delay location adjusted to better emulate terminal usage

## [0.0.1] - 2020-03-14

### Added

- Clide demos can be run using a json file
- Examples folder to hold example json demos
- Makefile to build and run program quickly from project directory
- Github actions build program when pushees are sent to master
- Issue templates for bugs and features
- Preliminary README.md 

[1.1.0]: https://github.com/mattackard/clide/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/mattackard/clide/compare/v0.0.2...v1.0.0
[0.0.2]: https://github.com/mattackard/clide/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/mattackard/clide/releases/tag/v0.0.1
