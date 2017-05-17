# scrapbox

This tool provides command line interface for scrapbox.io.

## Description

This is a tool to search pages by keywords, to print a content of a page, to print an encoded URL of a page, to print URLs linked by a page.

## Usage

### List page titles containing specified tags

```console
$ scrapbox list -h
usage: scrapbox list [options...] PROJECT [TAGs...]

Options:
  --token, -t  Scrapbox connect.sid used to access private project.
  --host, -h   Scrapbox Host. By default, "https://scrapbox.io".
  --expire     Local Cache Expiration. By default, 3600 seconds.


$ scrapbox list go-scrapbox
title having paren ( ) mark
title having plus + mark
title having question ? mark
title having slash / mark
title having whitespaces
日本語タイトルのページ
地のリンクがあるページ
複数のリンクがあるページ
文章のなかにリンクがあるページ1
文章のなかにリンクがあるページ2

$ scrapbox list go-scrapbox english
title having paren ( ) mark
title having plus + mark
title having question ? mark
title having slash / mark
title having whitespaces

$ scrapbox list go-scrapbox english paren
title having paren ( ) mark
```

### Print the content of the scrapbox page

```console
$ scrapbox read -h
usage: scrapbox read [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid used to access private project.
	--host, -h   Scrapbox Host. By default, "https://scrapbox.io".
  --expire     Local Cache Expiration. By default, 3600 seconds.


$ scrapbox read go-scrapbox "title having paren ( ) mark"
title having paren ( ) mark
#english #no-url #whitespace #no-slash #paren #no-plus #no-question
```

### Print the URL of the scrapbox page

```console
$ scrapbox open -h
usage: scrapbox open [options...] PROJECT PAGE

Options:
	--host, -h   Scrapbox Host. By default, "https://scrapbox.io".


$ scrapbox open go-scrapbox "title having paren ( ) mark"
https://scrapbox.io/go-scrapbox/title%20having%20paren%20(%20)%20mark
```

### Print all URLs in the scrapbox page

```console
$ scrapbox link -h
usage: scrapbox link [options...] PROJECT PAGE

Options:
  --token, -t  Scrapbox connect.sid used to access private project.
	--host, -h   Scrapbox Host. By default, "https://scrapbox.io".
  --expire     Local Cache Expiration. By default, 3600 seconds.


$ scrapbox link go-scrapbox "複数のリンクがあるページ"
https://www.google.co.jp
https://www.google.com
```

### Environment Variables

- `SCRAPBOX_TOKEN`: specify `token` instead of `--token` option.
- `SCRAPBOX_HOST`: specify `host` instead of `--host` option.
- `SCRAPBOX_EXPIRATION`: specify `expire` instead of `--expire` option.
- `SCRAPBOX_HOME`: specify `scrapbox` home directory. By default `~/.scrapbox/`
- `SCRAPBOX_DEBUG`: whether or not print stack trace at error.
- `SCRAPBOX_LONG_RUN_TEST`: execute long-run test.

### Private Project

To access private project, use `--token` option:

```console
$ scrapbox <sub command> --token s%3A... <arguments>
```

### Scrapbox Enterprise

To access Scrapbox Enterprise, use `--host` option:

```console
$ scrapbox <sub command> --host http://host:port <arguments>
```

### Local Cache Control

To ignore local caches, set `expire` to zero:

```console
$ scrapbox <sub command> --expire <expiration> <arguments>
```

## Install

To install, use `go get`:

```console
$ go get -d github.com/ohtomi/scrapbox
```

Or get binary from [release page](../../releases/latest).

## Contribution

1. Fork ([https://github.com/ohtomi/scrapbox/fork](https://github.com/ohtomi/scrapbox/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## License

MIT

## Author

[Kenichi Ohtomi](https://github.com/ohtomi)
