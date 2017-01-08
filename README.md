# scrapbox

This tool provides command line interface for scrapbox.io.

## Description

## Usage

```bash
$ scrapbox import ohtomi bookmark
Imported keyword data from https://scrapbox.io/ohtomi to ~/.go-scrapbox/scrapbox.io/ohtomi/bookmark/db/

$ scrapbox list ohtomi bookmark
Go Advent Calendar 2016 - Qiita #Go #adventcalendar #Qiita #Bookmark
Go (その2) Advent Calendar 2016 - Qiita #Go #adventcalendar #Qiita #Bookmark
Go (その3) Advent Calendar 2016 - Qiita #Go #adventcalendar #Qiita #Bookmark
高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA #gcli #Go #generator #Bookmark
...

$ scrapbox show ohtomi "高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA"
高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA
http://deeeet.com/writing/2014/06/22/cli-init/
https://github.com/tcnksm/gcli

#gcli #Go #generator #Bookmark

$ scrapbox open ohtomi "高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA"
-> open http://deeeet.com/writing/2014/06/22/cli-init/

$ scrapbox download ohtomi "高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA" ./
$ ls .
高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA

$ scrapbox upload ohtomi "./高速にGo言語のCLIツールをつくるcli-initというツールをつくった | SOTA"
```

### Local Cache

```bash
$ scrapbox show --no-cache project-name page-name
$ scrapbox open --no-cache project-name page-name
```

### Private Project

```bash
$ scrapbox import   --token "your token" project-name tag-name
$ scrapbox show     --token "your token" project-name page-name
$ scrapbox download --token "your token" project-name page-name /path/to/
$ scrapbox upload   --token "your token" project-name /path/to/page-name
```

### Scrapbox Enterprise

```bash
$ scrapbox import   --url http://host:port/ project-name tag-name
$ scrapbox list     --url http://host:port/ project-name tag-name
$ scrapbox show     --url http://host:port/ project-name page-name
$ scrapbox open     --url http://host:port/ project-name page-name
$ scrapbox download --url http://host:port/ project-name page-name /path/to/
$ scrapbox upload   --url http://host:port/ project-name /path/to/page-name
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/ohtomi/scrapbox
```

## Contribution

1. Fork ([https://github.com/ohtomi/scrapbox/fork](https://github.com/ohtomi/scrapbox/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[Kenichi Ohtomi](https://github.com/ohtomi)
