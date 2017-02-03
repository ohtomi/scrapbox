#!/bin/bash

function usage() {
  echo "
Usage: $0 [fmt|stringer|build|prep|test|install|package|release]
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
    go build -v \
      -ldflags "-X main.GitCommit=$(git describe --always)"
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
    env SCRAPBOX_HOME="`pwd`/testdata" SCRAPBOX_EXPIRATION=1 go test github.com/ohtomi/scrapbox/command -v
    ;;
  "install")
    go install \
      -ldflags "-X main.GitCommit=$(git describe --always)"
    ;;
  "package")
    $0 stringer

    rm -fr ./pkg
    gox \
      -ldflags "-X main.GitCommit=$(git describe --always)" \
      -os="darwin linux windows" \
      -arch="386 amd64" \
      -output "pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

    repo=$(grep "const Name " version.go | sed -E 's/.*"(.+)"$/\1/')
    version=$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/')

    rm -fr ./dist/${version}
    mkdir -p ./dist/${version}
    for platform in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
      platform_name=$(basename ${platform})
      archive_name=${repo}_${version}_${platform_name}
      pushd ${platform}
      zip ../../dist/${version}/${archive_name}.zip ./*
      popd
    done

    pushd ./dist/${version}
    shasum -a 256 * > ./${version}_SHASUMS
    popd
    ;;
  "release")
    version=$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/')
    ghr ${version} ./dist/${version}
    ;;
  *)
    usage
    exit 1
    ;;
esac
