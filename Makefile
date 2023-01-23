.DEFAULT_GOAL := build

# Driver Binary Name
BIN_NAME := docker-machine-driver-ionoscloud
GOFILES_NOVENDOR=$(shell find . -type f -name '*.go' | grep -v vendor)
ifeq ($(OS),Windows_NT)
	BIN_SUFFIX := .exe
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
	@echo "NOTE: Please copy the binary somewhere in your PATH"
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
	@go test -cover ./...
	@echo "DONE"

.PHONY: gofmt_check
gofmt_check:
	@echo "Ensure code adheres to gofmt and list files whose formatting differs from gofmt's (vendor directory excluded)"
	@if [ "$(shell echo $$(gofmt -l ${GOFILES_NOVENDOR}))" != "" ]; then (echo "Format files: $(shell echo $$(gofmt -l ${GOFILES_NOVENDOR})) Hint: use \`make gofmt_update\`"; exit 1); fi
	@echo "DONE"

.PHONY: gofmt_update
gofmt_update:
	@echo "Ensure code adheres to gofmt and change files accordingly (vendor directory excluded)"
	@gofmt -w ${GOFILES_NOVENDOR}
	@echo "DONE"

.PHONY: mock_update
mock_update:
	@echo "Update mock for tests"
	@mockgen -source=internal/utils/client_service.go > internal/utils/mocks/ClientService.go
	@echo "DONE"

.PHONY: vendor_status
vendor_status:
	@govendor status

.PHONY: vendor_update
vendor_update:
	@echo "Update vendor dependencies"
	@go mod vendor
	@go mod tidy
	@echo "DONE"

.PHONY: clean
clean:
	rm -f ./bin/${BIN_NAME}*
	rm -f ${GOPATH}/bin/${BIN_NAME}*

.PHONY: upload
upload: clean compile
ifdef to_clip
	./scripts/publish_image.py -c True
else
	./scripts/publish_image.py
endif
