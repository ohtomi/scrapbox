#!/bin/bash

function usage() {
  echo "
Usage: $0 [fmt|stringer|build|prep|test|install]
"
}

if [ $# -ne 1 ]; then
  usage
  exit 1
fi

case "$1" in
  "fmt")
    gofmt -w .
    ;;
  "stringer")
    cd command
    stringer -type ExitCode -output meta_exitcode_string.go meta.go
    ;;
  "build")
    go build -v
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
    echo cleaning up ~/.scrapbox ...
    rm -fr ~/.scrapbox
    echo
    echo testing ...
    env SCRAPBOX_HOME="`pwd`/testdata" go test github.com/ohtomi/scrapbox/command -v
    ;;
  "install")
    go install
    ;;
  *)
    usage
    exit 1
    ;;
esac
