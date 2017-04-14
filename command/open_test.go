package command

import (
	"bytes"
	"strings"
	"testing"

	_ "github.com/mitchellh/cli"
)

func TestOpenCommand__print_url_having_paren(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having paren ( ) mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20paren%20(%20)%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_plus(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having plus + mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20plus%20%2B%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_question(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having question ? mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20question%20%3F%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_slash(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having slash / mark", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20slash%20%2F%20mark"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_whitespace(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  title having whitespaces", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/title%20having%20whitespaces"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}

func TestOpenCommand__print_url_having_japanese(t *testing.T) {

	outStream, errStream, inStream := new(bytes.Buffer), new(bytes.Buffer), strings.NewReader("")
	meta := NewTestMeta(outStream, errStream, inStream)
	command := &OpenCommand{
		Meta: *meta,
	}

	args := strings.Split("go-scrapbox  日本語タイトルのページ", "  ")
	exitStatus := command.Run(args)
	if ExitCode(exitStatus) != ExitCodeOK {
		t.Fatalf("ExitStatus is %s, but want %s", ExitCode(exitStatus), ExitCodeOK)
	}

	expected := "https://scrapbox.io/go-scrapbox/%E6%97%A5%E6%9C%AC%E8%AA%9E%E3%82%BF%E3%82%A4%E3%83%88%E3%83%AB%E3%81%AE%E3%83%9A%E3%83%BC%E3%82%B8"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("Output is %q, but want %q", outStream.String(), expected)
	}
}
