set-env:
	cp -r assets /usr/share/clide/
	cp -r examples /usr/share/clide/

builds:
	go build -o clide cmd/clide/*.go
	mv clide build/bin/

run-demo:
	go run cmd/clide/*.go examples/demo.json