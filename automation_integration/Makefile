all: test

clean:
	go clean

cleanall:
	go clean -r
	if [ -f output.txt ]; then rm output.txt; fi

build:
	go build

test:
	$(MAKE) cleanall
	$(MAKE) build
	./automation_integration
