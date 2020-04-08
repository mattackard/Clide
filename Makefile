# adds assets used by the program to their install location and adds binaries to /usr/bin
install: builds
	mkdir /usr/share/clide
	cp -r assets /usr/share/clide/
	cp -r cmd/clide-editor /usr/share/clide/editor/
	cp -r examples /usr/share/clide/
	cp build/bin/clide /usr/bin/
	cp build/bin/clide-editor /usr/bin/

# removes all clide configuration files from /usr filesystem and binaries from /usr/bin
uninstall:
	rm -r /usr/share/clide
	rm /usr/bin/clide
	rm /usr/bin/clide-editor

# builds clide and clide-editor into binaries
builds:
	go build -o clide cmd/clide/*.go
	mv clide build/bin/

	go build -o clide-editor cmd/clide-editor/main.go
	mv clide-editor build/bin/

# runs demo.json localated in the examples folder in project directory
run-demo:
	go run cmd/clide/*.go examples/demo.json

# creates a package ready to package for a debian distribution
create-package: builds
	mkdir -p 					./build/pkg/clide-pkg/usr/share
	mkdir  						./build/pkg/clide-pkg/usr/bin
	mkdir						./build/pkg/clide-pkg/DEBIAN
	cp -r assets 				./build/pkg/clide-pkg/usr/share/
	cp -r examples 				./build/pkg/clide-pkg/usr/share/
	cp -r cmd/clide-editor 		./build/pkg/clide-pkg/usr/share/
	cp CHANGELOG.md 			./build/pkg/clide-pkg/
	cp LICENSE 					./build/pkg/clide-pkg/
	rm							./build/pkg/clide-pkg/usr/share/clide-editor/main.go
	cp ./build/bin/clide		./build/pkg/clide-pkg/usr/bin/
	cp ./build/bin/clide-editor	./build/pkg/clide-pkg/usr/bin/

	echo "Package: clide\nVersion: 1.2.0\nSection: base\nPriority: optional\nArchitecture: amd64\nDepends: libsdl2-dev (>= 2.0.8), libsdl2-gfx-dev (>= 1.0.4), libsdl2-mixer-dev (>= 2.0.2), libsdl2-ttf-dev (>= 2.0.14), libsdl2-image-dev (>= 2.0.3)\nMaintainer: Matt Ackard <mattackard@gmail.com>\nDescription: clide\n Clide is a tool for creating automated demos for terminal-based applications." > build/pkg/clide-pkg/DEBIAN/control
