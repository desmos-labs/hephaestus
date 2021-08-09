BUILDDIR ?= $(CURDIR)/build

export GO111MODULE = on

###############################################################################
###                                   All                                   ###
###############################################################################

all: lint build test-unit

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

build-linux: go.sum
	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

build-arm64: go.sum
	GOOS=linux GOARCH=arm64 $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

tools:
	@go get -u github.com/client9/misspell/cmd/misspell
	@go get golang.org/x/tools/cmd/goimports

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

clean:
	rm -rf $(BUILDDIR)/

.PHONY: clean

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test-unit:
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -w -local github.com/desmos-labs/hephaestus
.PHONY: format