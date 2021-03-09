# Driver Binary Name
BIN_NAME := docker-machine-driver-ionoscloud

GOFILES= $(shell find . -type f -name '*.go')
ifeq ($(OS),Windows_NT)
	BIN_SUFFIX := ".exe"
else
	BIN_SUFFIX :=
endif

.PHONY: build
build: compile print-success

.PHONY : install
install: compile
ifeq ($(OS),Windows_NT)
	cp bin/${BIN_NAME} $(GOPATH)/bin/
else
	cp ./bin/${BIN_NAME}${BIN_SUFFIX} ${GOPATH}/bin/
endif

.PHONY: compile
compile:
	GOGC=off CGOENABLED=0 go build -o ./bin/${BIN_NAME}${BIN_SUFFIX} ./bin

.PHONY: print-success
print-success:
	@echo
	@echo "Plugin built."
	@echo
	@echo "To use it, either run 'make install' or set your PATH environment variable correctly."

.PHONY: test
test: test_unit

.PHONY: test_unit
test_unit:
	@echo "Run unit tests"
	@go test -cover .
	@echo "DONE"

.PHONY: gofmt_check
gofmt_check:
	@echo "Ensure code adheres to gofmt and list files whose formatting differs from gofmt's"
	@if [ "$(shell echo $$(gofmt -l ${GOFILES}))" != "" ]; then (echo "Format files: $(shell echo $$(gofmt -l ${GOFILES})) Hint: use \`make gofmt_update\`"; exit 1); fi
	@echo "DONE"

.PHONY: gofmt_update
gofmt_update:
	@echo "Ensure code adheres to gofmt and change files accordingly"
	@gofmt -w ${GOFILES}
	@echo "DONE"

.PHONY: clean
clean:
	rm -f ./bin/${BIN_NAME}*
	rm -f ${GOPATH}/bin/${BIN_NAME}*
