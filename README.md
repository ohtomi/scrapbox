# scrapbox

This tool provides command line interface for scrapbox.io.

## Description

## Usage

```bash
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

$ scrapbox list --tags go-scrapbox english
title having paren ( ) mark --- #english #no-url #whitespace #no-slash #paren #no-plus #no-question
title having plus + mark --- #english #no-url #whitespace #no-slash #no-paren #plus #no-question
title having question ? mark --- #english #no-url #whitespace #no-slash #no-paren #no-plus #question
title having slash / mark --- #english #no-url #whitespace #slash #no-paren #no-plus #no-question
title having whitespaces --- #english #no-url #whitespace #no-slash #no-paren #no-plus #no-question

$ scrapbox read --no-cache go-scrapbox "title having paren ( ) mark"
title having paren ( ) mark
#english #no-url #whitespace #no-slash #paren #no-plus #no-question

$ scrapbox open go-scrapbox "title having paren ( ) mark"
https://scrapbox.io/go-scrapbox/title%20having%20paren%20(%20)%20mark

$ scrapbox link go-scrapbox "複数のリンクがあるページ"
https://www.google.co.jp
https://www.google.com
```

### Environment Variables

- `SCRAPBOX_TOKEN`: specify `token` instead of `--token` option.
- `SCRAPBOX_HOST`: specify `host` instead of `--host` option.

### Private Project

To access private project, use `--token` option:

```bash
$ scrapbox <sub command> --token s%3A... <arguments>
```

### Scrapbox Enterprise

To access Scrapbox Enterprise, use `--host` option:

```bash
$ scrapbox <sub command> --host http://host:port <arguments>
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
