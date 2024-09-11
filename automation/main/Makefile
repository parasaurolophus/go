all: test

clean:
	if [ -f automation ]; then rm automation; fi

cleanall: clean
	go clean -r

build:
	go build -o automation

test: cleanall build
	./automation
