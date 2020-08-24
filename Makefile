# go installation
GOPATH=$(HOME)/go
GOBINDIR=$(GOPATH)/bin

# commands
GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
MOCKGEN=$(GOBINDIR)/mockgen

BINARY_NAME=trac2gitea
PACKAGES=$(shell go list ./...)
ROOTPACKAGE=github.com/stevejefferson/trac2gitea

MOCKFILES=\
	accessor/mock_gitea/accessor.go \
	accessor/mock_trac/accessor.go

.PHONY: all install build test
all: build test

install: build
	$(GOINSTALL)

test: build mocks
	$(GOTEST) ./...

build: mocks
	$(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: allclean clean
allclean: mockclean clean 

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: alldeps deps
alldeps: deps mockdeps lintdeps

deps:
	$(GOGET) github.com/go-ini/ini 
	$(GOGET) github.com/mattn/go-sqlite3
	$(GOGET) github.com/spf13/pflag
	$(GOGET) gopkg.in/src-d/go-git.v4

.PHONY: mocks mockdeps mockclean
mocks: $(MOCKFILES)

# mock generation:
accessor/mock_gitea/accessor.go: accessor/gitea/accessor.go
	$(MOCKGEN) -destination=$@ $(ROOTPACKAGE)/$(<D) Accessor

accessor/mock_trac/accessor.go: accessor/trac/accessor.go
	$(MOCKGEN) -destination=$@ $(ROOTPACKAGE)/$(<D) Accessor

mockdeps:
	GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.3

mockclean:
	rm -rf accessor/mock_gitea accessor/mock_giteawiki accessor/mock_trac

.PHONY: lint lintdeps
lint:
	@for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

lintdeps:
	$(GOGET) -u github.com/golang/lint/golint

.PHONY: modtidy
modtidy:
	$(GOMOD) tidy