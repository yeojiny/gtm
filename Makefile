BINARY=gtm
VERSION=v1.2.8-beta

LDFLAGS=-ldflags "-X main.Version=${VERSION}"
BLD_TAGS=--tags static

build:
	go build ${BLD_TAGS} ${LDFLAGS} -o ${BINARY}

test:
	go test ${BLD_TAGS} $$(go list ./... | grep -v vendor)

vet:
	go vet $$(go list ./... | grep -v vendor)

fmt:
	go fmt $$(go list ./... | grep -v vendor)

install-git2go:
	cd ${GOPATH}/src/github.com/git-time-metric/git2go; git checkout v25; git submodule update --init; make install-static

install:
	go install ${BLD_TAGS} ${LDFLAGS}

clean:
	go clean

.PHONY: test vet install clean fmt todo note
