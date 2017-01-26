#!/bin/bash

if [ $# -ne 1 ]; then
  echo "
Usage: $0 [fmt|stringer|build|prep|test]
"
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
    env SCRAPBOX_DEBUG=1 ./scrapbox import go-scrapbox english
    env SCRAPBOX_DEBUG=1 ./scrapbox import go-scrapbox japanese
    echo
    ls -lR ./testdata/scrapbox.io/go-scrapbox
    ;;
  "test")
    echo cleaning up ~/.scrapbox ...
    rm -fr ~/.scrapbox
    echo
    echo testing ...
    go test github.com/ohtomi/scrapbox/command -v
    ;;
esac
