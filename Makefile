default: build

version := "v1.3.4"

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
name := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

ifeq ($(OS),Windows_NT)
	bin_suffix := ".exe"
else
	bin_suffix :=
endif

clean:
	rm -f ./bin/$(name)*
	rm -f /usr/local/bin/$(name)*

compile:
	GOGC=off CGOENABLED=0 go build -i -o ./bin/$(name)$(bin_suffix) ./bin

print-success:
	@echo
	@echo "Plugin built."
	@echo
	@echo "To use it, either run 'make install' or set your PATH environment variable correctly."

build: compile print-success

release:
	GOOS=linux GOARCH=amd64 GOGC=off CGOENABLED=0 go build  -i -o bin/$(name) ./bin
	tar  -cvzf bin/$(name)-$(version)-linux-amd64.tar.gz -C bin $(name)
	GOOS=darwin GOARCH=amd64 GOGC=off CGOENABLED=0 go build  -i -o bin/$(name) ./bin
	tar -cvzf bin/$(name)-$(version)-darwin-amd64.tar.gz -C bin $(name)
	GOOS=windows GOARCH=amd64 GOGC=off CGOENABLED=0 go build  -i -o bin/$(name).exe ./bin
	tar -cvzf bin/$(name)-$(version)-windows-amd64.tar.gz -C bin $(name).exe

install: compile
ifeq ($(OS),Windows_NT)
	cp bin/$(name) $(GOPATH)/bin/
else
	cp ./bin/$(name)$(bin_suffix) ${GOPATH}/bin/
endif


.PHONY : build release install
