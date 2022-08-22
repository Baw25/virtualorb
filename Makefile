
compile:
	go build 

clean:
	(cd osx-cpu-temp && rm -rf osx-cpu-temp)

build_osx_cputemp:
	(cd osx-cpu-temp && make)

build: clean compile build_osx_cputemp

run: 
	PORT=5000 ./virtualorb

test: 
	go test -v ./...
