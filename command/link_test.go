package command

import (
	"bytes"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func TestLinkCommand__print_http_link(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &LinkCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox HTTPなリンクのあるページ", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "http://www.sphinx-doc.org/en/stable/"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestLinkCommand__print_https_link(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &LinkCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox HTTPSなリンクのあるページ", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://www.google.co.jp"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestLinkCommand__print_link_with_name_1(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &LinkCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox 文章のなかにリンクがあるページ1", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://www.google.co.jp"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestLinkCommand__print_link_with_name_2(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &LinkCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox 文章のなかにリンクがあるページ2", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://www.google.com"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestLinkCommand__print_multiple_links(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &LinkCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox 複数のリンクがあるページ", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://www.google.co.jp\nhttps://www.google.com"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestLinkCommand__print_no_links(t *testing.T) {

	InitializeMeta()

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &LinkCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox 日本語タイトルのページ", " ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := ""
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}
