all: test

clean:
	if [ -f automation ]; then rm automation; fi
	if [ -f output.txt ]; then rm output.txt; fi

cleanall: clean
	go clean -r

build:
	go build -o automation

test:
	$(MAKE) cleanall
	$(MAKE) build
	./automation
