
compile:
	go build 

clean:
	(cd osx-cpu && rm -rf osx-cpu)

build_osx_cpu:
	(cd osx-cpu && make)

build: clean compile build_osx_cpu

run: 
	PORT=5000 ./virtualorb

test: 
	go test -v ./...
