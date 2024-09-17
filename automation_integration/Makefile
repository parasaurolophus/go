all: cleanall build

clean:
	go clean

cleanall:
	go clean -r
	if [ -f output.txt ]; then rm output.txt; fi

build:
	go build

test: all
	./automation_integration
