package syntax

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var (
	enablePrettyPrint = os.Getenv("SCRAPBOX_DEBUG") != ""
)

func TestParse__indent_level(t *testing.T) {
	for _, fixture := range []struct {
		source string
		indent int
	}{
		{" ", 1},
		{"\t", 1},
		{" \t ", 3},
		{"\t \t", 3},
	} {
		queryable := Parse([]byte(fixture.source), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}

		if len(queryable.GetChildren()) > 1 {
			t.Fatalf("%d root children found", len(queryable.GetChildren()))
		}
		node := queryable.GetChildren()[0]

		assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})
	}
}

func TestParse__single_node(t *testing.T) {
	for _, fixture := range []struct {
		source string
		name   string
		value  string
	}{
		// link
		{"[$ 1+2 = 3]", "link", "[$ 1+2 = 3]"},
		{"[_-/*/-_ https://avatars1.githubusercontent.com/u/1678258#.png]", "link", "[_-/*/-_ https://avatars1.githubusercontent.com/u/1678258#.png]"},
		{"[_-/*/-_ github.com/ohtomi/scrapbox]", "link", "[_-/*/-_ github.com/ohtomi/scrapbox]"},
		{"[/foo/bar/baz]", "link", "[/foo/bar/baz]"},
		{"[https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258]", "link", "[https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258]"},
		{"[https://avatars1.githubusercontent.com/u/1678258 https://avatars1.githubusercontent.com/u/1678258#.png]", "link", "[https://avatars1.githubusercontent.com/u/1678258 https://avatars1.githubusercontent.com/u/1678258#.png]"},
		{"[avatar https://avatars1.githubusercontent.com/u/1678258]", "link", "[avatar https://avatars1.githubusercontent.com/u/1678258]"},
		{"[https://avatars1.githubusercontent.com/u/1678258 avatar]", "link", "[https://avatars1.githubusercontent.com/u/1678258 avatar]"},
		{"[https://avatars1.githubusercontent.com/u/1678258]", "link", "[https://avatars1.githubusercontent.com/u/1678258]"},
		{"[ user.icon]", "link", "[ user.icon]"},
		{"[github.com/ohtomi/scrapbox]", "link", "[github.com/ohtomi/scrapbox]"},
		// image
		{"http://avatars1.githubusercontent.com/u/1678258#.png", "image", "http://avatars1.githubusercontent.com/u/1678258#.png"},
		{"http://avatars1.githubusercontent.com/u/1678258#.gif", "image", "http://avatars1.githubusercontent.com/u/1678258#.gif"},
		{"https://avatars1.githubusercontent.com/u/1678258#.jpg", "image", "https://avatars1.githubusercontent.com/u/1678258#.jpg"},
		{"https://avatars1.githubusercontent.com/u/1678258#.jpeg", "image", "https://avatars1.githubusercontent.com/u/1678258#.jpeg"},
		{"https://gyazo.com/1678258/avatar", "image", "https://gyazo.com/1678258/avatar"},
		// url
		{"http://avatars1.githubusercontent.com/u/1678258", "url", "http://avatars1.githubusercontent.com/u/1678258"},
		{"https://avatars1.githubusercontent.com/u/1678258", "url", "https://avatars1.githubusercontent.com/u/1678258"},
		// bold
		{"[[http://avatars1.githubusercontent.com/u/1678258#.png]]", "bold", "[[http://avatars1.githubusercontent.com/u/1678258#.png]]"},
		{"[[github.com/ohtomi/scrapbox]]", "bold", "[[github.com/ohtomi/scrapbox]]"},
		{"[[ github.com\t/ohtomi/\tscrapbox ]]", "bold", "[[ github.com\t/ohtomi/\tscrapbox ]]"},
		// tag
		{"#[github.com/ohtomi/scrapbox/]", "tag", "#[github.com/ohtomi/scrapbox/]"},
		{"#[ github.com\t/ohtomi/\tscrapbox/ ]", "tag", "#[ github.com\t/ohtomi/\tscrapbox/ ]"},
		{"#github.com/ohtomi/scrapbox", "tag", "#github.com/ohtomi/scrapbox"},
		// text
		{"x github.com\t/ohtomi/\tscrapbox/ x", "text", "x github.com\t/ohtomi/\tscrapbox/ x"},
	} {
		queryable := Parse([]byte(fixture.source), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}

		if len(queryable.GetChildren()) > 1 {
			t.Fatalf("%d root children found", len(queryable.GetChildren()))
		}
		node := queryable.GetChildren()[0]

		if len(node.GetChildren()) > 1 {
			t.Fatalf("%d children found", len(node.GetChildren()))
		}
		item := node.GetChildren()[0]

		assertEqualTo(t, item.GetName(), fixture.name)
		assertEqualTo(t, item.GetValue(), fixture.value)
	}
}

func TestParse__quoted_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   []int
		expected [][]string
	}{
		{
			">https://avatars1.githubusercontent.com/u/1678258",
			[]int{0},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258"},
			},
		},
		{
			"   >https://avatars1.githubusercontent.com/u/1678258",
			[]int{3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258"},
			},
		},
		{
			"\t\t\t>https://avatars1.githubusercontent.com/u/1678258",
			[]int{3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258"},
			},
		},
		{
			">https://avatars1.githubusercontent.com/u/1678258#1\n" +
				"   >https://avatars1.githubusercontent.com/u/1678258#2\n" +
				"\t\t\t>https://avatars1.githubusercontent.com/u/1678258#3",
			[]int{0, 3, 3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#1"},
				{"https://avatars1.githubusercontent.com/u/1678258#2"},
				{"https://avatars1.githubusercontent.com/u/1678258#3"},
			},
		},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected) {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected), queryable)
		}

		for i, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "quoted_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent[i])})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "url")
				assertEqualTo(t, node.GetChildren()[j].GetValue(), expected)
			}
		}
	}
}

func TestParse__code_directive_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   []int
		expected [][]string
	}{
		{"code:sample.js",
			[]int{0},
			[][]string{
				{"sample.js"},
			},
		},
		{"   code:sample.js",
			[]int{3},
			[][]string{
				{"sample.js"},
			},
		},
		{"\t\t\tcode:sample.js",
			[]int{3},
			[][]string{
				{"sample.js"},
			},
		},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected) {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected), queryable)
		}

		for i, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "code_block")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent[i])})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "text")
				assertEqualTo(t, node.GetChildren()[j].GetValue(), expected)
			}
		}
	}
}

func TestParse__table_directive_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   []int
		expected [][]string
	}{
		{"table:sample",
			[]int{0},
			[][]string{
				{"sample"},
			},
		},
		{"   table:sample",
			[]int{3},
			[][]string{
				{"sample"},
			},
		},
		{"\t\t\ttable:sample",
			[]int{3},
			[][]string{
				{"sample"},
			},
		},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected) {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected), queryable)
		}

		for i, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "table_block")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent[i])})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "text")
				assertEqualTo(t, node.GetChildren()[j].GetValue(), expected)
			}
		}
	}
}

func assertEqualTo(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Got %+v, but Want %+v", actual, expected)
	}
}
