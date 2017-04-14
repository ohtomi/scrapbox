#!/bin/bash

MAIN_PACKAGE=.
REL_TO_ROOT=.

TEST_ENVIRONMENT="SCRAPBOX_DEBUG=1 SCRAPBOX_LONG_RUN_TEST=$3 SCRAPBOX_HOME=`pwd`/testdata SCRAPBOX_EXPIRATION=1"

GOX_ALL_OS="darwin linux windows"
GOX_ALL_ARCH="386 amd64"
GOX_MAIN_OS="darwin"
GOX_MAIN_ARCH="amd64"


function usage() {
  echo "
Usage: $0 [fmt|stringer|compile|prep|test|package|release]
"
}


case "$1" in
  "fmt")
    gofmt -w .
    ;;
  "stringer")
    cd command
    stringer -type ExitCode -output meta_exitcode_string.go meta.go
    ;;
  "compile")
    $0 stringer

    cd "${MAIN_PACKAGE}"
    gox \
      -ldflags "-X main.GitCommit=$(git describe --always)" \
      -os="${GOX_MAIN_OS}" \
      -arch="${GOX_MAIN_ARCH}" \
      -output "${REL_TO_ROOT}/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"
    ;;
  "prep")
    echo cleaning up ./testdata ...
    rm -fr ./testdata
    echo
    echo running debug mode to dump api reponse ...
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
    echo
    ls -l ./testdata/query/scrapbox.io/go-scrapbox/english
    ls -l ./testdata/page/scrapbox.io/go-scrapbox
    ;;
  "test")
    env ${TEST_ENVIRONMENT} go test ./... $2
    ;;
  "package")
    $0 stringer

    cd "${MAIN_PACKAGE}"
    rm -fr "${REL_TO_ROOT}/pkg"
    gox \
      -ldflags "-X main.GitCommit=$(git describe --always)" \
      -os="${GOX_ALL_OS}" \
      -arch="${GOX_ALL_ARCH}" \
      -output "${REL_TO_ROOT}/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

    repo=$(grep "const Name " version.go | sed -E 's/.*"(.+)"$/\1/')
    version=$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/')
    cd "${REL_TO_ROOT}"

    rm -fr "./dist/${version}"
    mkdir -p "./dist/${version}"
    for platform in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
      platform_name=$(basename ${platform})
      archive_name=${repo}_${version}_${platform_name}
      pushd "${platform}"
      zip "../../dist/${version}/${archive_name}.zip" ./*
      popd
    done

    pushd "./dist/${version}"
    shasum -a 256 * > "./${version}_SHASUMS"
    popd
    ;;
  "release")
    version=$(grep "const Version " ${MAIN_PACKAGE}/version.go | sed -E 's/.*"(.+)"$/\1/')
    ghr "${version}" "./dist/${version}"
    ;;
  *)
    usage
    exit 1
    ;;
esac
