MAIN_PACKAGE = $(dir $(shell grep -ir -l --exclude-dir vendor --exclude Makefile "func main()" ./*))
REPO = $(notdir $(CURDIR))
VERSION = $(shell grep 'Version string' $(MAIN_PACKAGE)/version.go | sed -E 's/.*"(.+)"$$/\1/')
COMMIT = $(shell git describe --always)
PACKAGES = $(shell go list ./... | grep -v '/vendor/')

GOX_OS = darwin linux windows
GOX_ARCH = amd64 386

default: test

build: stringer
	@cd $(MAIN_PACKAGE) ; \
	gox \
	  -ldflags "-X main.GitCommit=$(COMMIT)" \
	  -os="$(firstword $(GOX_OS))" \
	  -arch="$(firstword $(GOX_ARCH))" \
	  -output="$(CURDIR)/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

prep:
	@rm -fr ./testdata

	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go list go-scrapbox
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go list go-scrapbox english
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go list go-scrapbox english paren
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go list go-scrapbox english whitespaces
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "title having paren ( ) mark"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "title having plus + mark"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "title having question ? mark"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "title having slash / mark"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "title having whitespaces"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "日本語タイトルのページ"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "HTTPなリンクのあるページ"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "HTTPSなリンクのあるページ"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "地のリンクがあるページ"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "複数のリンクがあるページ"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "文章のなかにリンクがあるページ1"
	env SCRAPBOX_HOME="`pwd`/testdata" go run ./*.go read go-scrapbox "文章のなかにリンクがあるページ2"

test:
	rm -fr ./testdata/query/127.0.0.1
	rm -fr ./testdata/page/127.0.0.1
	env SCRAPBOX_DEBUG=1 SCRAPBOX_LONG_RUN_TEST=1 SCRAPBOX_HOME=`pwd`/testdata SCRAPBOX_EXPIRATION=1 go test -v -parallel=4 ${PACKAGES}

test-race:
	go test -v -race ${PACKAGES}

vet:
	go vet ${PACKAGES}

clean:
	@rm -fr ./pkg
	@rm -fr ./dist/$(VERSION)

install: clean build
	cp "$(CURDIR)/pkg/$(firstword $(GOX_OS))_$(firstword $(GOX_ARCH))/$(REPO)" "${GOPATH}/bin"

package: clean stringer
	@cd $(MAIN_PACKAGE) ; \
	gox \
	  -ldflags "-X main.GitCommit=$(COMMIT)" \
	  -parallel=3 \
	  -os="$(GOX_OS)" \
	  -arch="$(GOX_ARCH)" \
	  -output="$(CURDIR)/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

	@mkdir -p ./dist/$(VERSION)

	@for platform in $(foreach os,$(GOX_OS),$(foreach arch,$(GOX_ARCH),$(os)_$(arch))) ; do \
	  echo "zip ../../dist/$(VERSION)/$(REPO)_$(VERSION)_$$platform.zip ./*" ; \
	  (cd ./pkg/$$platform && zip ../../dist/$(VERSION)/$(REPO)_$(VERSION)_$$platform.zip ./*) ; \
	done

	@cd ./dist/$(VERSION) ; \
	echo "shasum -a 256 * > ./$(VERSION)_SHASUMS" ; \
	shasum -a 256 * > ./$(VERSION)_SHASUMS

release: package
	ghr $(VERSION) ./dist/$(VERSION)

fmt:
	gofmt -w .

stringer:
	@cd command ; \
	stringer -type ExitCode -output meta_exitcode_string.go meta.go

.PHONY: build test test-race vet clean install package release fmt stringer
