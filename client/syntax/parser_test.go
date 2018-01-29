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

func TestParse__indent_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   int
	}{
		{"   ", 3},
		{"\t\t\t", 3},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}

		for _, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "simple_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})
		}
	}
}

func TestParse__quoted_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   int
		expected [][]string
	}{
		{
			">https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox",
			0,
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox"},
			},
		},
		{
			"   >https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox",
			3,
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox"},
			},
		},
		{
			"\t\t\t>https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox",
			3,
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox"},
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
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "quoted")
				if len(node.GetChildren()[j].GetChildren()) != 2 {
					t.Fatalf("Found %d, but Want %d: %+v", len(node.GetChildren()[j].GetChildren()), 2, node.GetChildren()[j])
				}
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[0].GetName(), "q")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[0].GetValue(), ">")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[1].GetName(), "t")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[1].GetValue(), expected)
			}
		}
	}
}

func TestParse__code_directive_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   int
		expected [][]string
	}{
		{"code:sample.js", 0, [][]string{{"sample.js"}}},
		{"   code:sample.js", 3, [][]string{{"sample.js"}}},
		{"\t\t\tcode:sample.js", 3, [][]string{{"sample.js"}}},
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
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "code")
				if len(node.GetChildren()[j].GetChildren()) != 2 {
					t.Fatalf("Found %d, but Want %d: %+v", len(node.GetChildren()[j].GetChildren()), 2, node.GetChildren()[j])
				}
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[0].GetName(), "c")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[0].GetValue(), "code:")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[1].GetName(), "n")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[1].GetValue(), expected)
			}
		}
	}
}

func TestParse__table_directive_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   int
		expected [][]string
	}{
		{"table:sample.js", 0, [][]string{{"sample.js"}}},
		{"   table:sample.js", 3, [][]string{{"sample.js"}}},
		{"\t\t\ttable:sample.js", 3, [][]string{{"sample.js"}}},
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
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "table")
				if len(node.GetChildren()[j].GetChildren()) != 2 {
					t.Fatalf("Found %d, but Want %d: %+v", len(node.GetChildren()[j].GetChildren()), 2, node.GetChildren()[j])
				}
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[0].GetName(), "t")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[0].GetValue(), "table:")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[1].GetName(), "n")
				assertEqualTo(t, node.GetChildren()[j].GetChildren()[1].GetValue(), expected)
			}
		}
	}
}

func TestParse__image_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   int
		expected [][]string
	}{
		{"https://avatars1.githubusercontent.com/u/1678258#.png", 0, [][]string{{"https://avatars1.githubusercontent.com/u/1678258#.png"}}},
		{"   https://avatars1.githubusercontent.com/u/1678258#.png", 3, [][]string{{"https://avatars1.githubusercontent.com/u/1678258#.png"}}},
		{"\t\t\thttps://avatars1.githubusercontent.com/u/1678258#.png", 3, [][]string{{"https://avatars1.githubusercontent.com/u/1678258#.png"}}},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected) {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected), queryable)
		}

		for i, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "simple_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "image")
				assertEqualTo(t, node.GetChildren()[j].GetValue(), expected)
			}
		}
	}
}

func TestParse__url_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   int
		expected [][]string
	}{
		{"https://avatars1.githubusercontent.com/u/1678258", 0, [][]string{{"https://avatars1.githubusercontent.com/u/1678258"}}},
		{"   https://avatars1.githubusercontent.com/u/1678258", 3, [][]string{{"https://avatars1.githubusercontent.com/u/1678258"}}},
		{"\t\t\thttps://avatars1.githubusercontent.com/u/1678258", 3, [][]string{{"https://avatars1.githubusercontent.com/u/1678258"}}},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected) {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected), queryable)
		}

		for i, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "simple_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})

			for j, expected := range fixture.expected[i] {
				assertEqualTo(t, node.GetChildren()[j].GetName(), "url")
				assertEqualTo(t, node.GetChildren()[j].GetValue(), expected)
			}
		}
	}
}

func TestParse__text_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   int
		expected [][]string
	}{
		{"github.com/ohtomi/scrapbox", 0, [][]string{{"github.com/ohtomi/scrapbox"}}},
		{"   github.com/ohtomi/scrapbox", 3, [][]string{{"github.com/ohtomi/scrapbox"}}},
		{"\t\t\tgithub.com/ohtomi/scrapbox", 3, [][]string{{"github.com/ohtomi/scrapbox"}}},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected) {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected), queryable)
		}

		for i, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "simple_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent)})

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
