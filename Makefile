builds:
	go build -o clide cmd/clide/*.go
	mv clide build/bin/

run-demo:
	go run cmd/clide/*.go examples/demo.json