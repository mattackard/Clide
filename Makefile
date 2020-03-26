# adds assets used by the program to their install location
set-env:
	cp -r assets /usr/share/clide/
	cp -r examples /usr/share/clide/

# removes all clide configuration files from /usr filesystem
rm-env:
	rm -r /usr/share/clide

# builds the go program into a binary
builds:
	go build -o clide cmd/clide/*.go
	mv clide build/bin/

# runs demo.json localated in the examples folder in project directory
run-demo:
	go run cmd/clide/*.go examples/demo.json