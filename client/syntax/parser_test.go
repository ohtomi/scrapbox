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
		indent   []int
	}{
		{"   ",
			[]int{3},
		},
		{"\t\t\t",
			[]int{3},
		},
		{"   \n" +
			"\t\t\t",
			[]int{3, 3},
		},
	} {
		queryable := Parse([]byte(fixture.original), enablePrettyPrint)

		if queryable == nil {
			t.Fatalf("Failed to parse")
		}

		for i, node := range queryable.GetChildren() {
			assertEqualTo(t, node.GetName(), "simple_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent[i])})
		}
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

func TestParse__image_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		indent   []int
		expected [][]string
	}{
		{"https://avatars1.githubusercontent.com/u/1678258#.png",
			[]int{0},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#.png"},
			},
		},
		{"   https://avatars1.githubusercontent.com/u/1678258#.png",
			[]int{3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#.png"},
			},
		},
		{"\t\t\thttps://avatars1.githubusercontent.com/u/1678258#.png",
			[]int{3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#.png"},
			},
		},
		{"https://avatars1.githubusercontent.com/u/1678258#.png\n" +
			"   https://avatars1.githubusercontent.com/u/1678258#.jpg\n" +
			"\t\t\thttps://avatars1.githubusercontent.com/u/1678258#.gif",
			[]int{0, 3, 3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258#.png"},
				{"https://avatars1.githubusercontent.com/u/1678258#.jpg"},
				{"https://avatars1.githubusercontent.com/u/1678258#.gif"},
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
			assertEqualTo(t, node.GetName(), "simple_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent[i])})

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
		indent   []int
		expected [][]string
	}{
		{"https://avatars1.githubusercontent.com/u/1678258",
			[]int{0},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258"},
			},
		},
		{"   https://avatars1.githubusercontent.com/u/1678258",
			[]int{3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258"},
			},
		},
		{"\t\t\thttps://avatars1.githubusercontent.com/u/1678258",
			[]int{3},
			[][]string{
				{"https://avatars1.githubusercontent.com/u/1678258"},
			},
		},
		{"https://avatars1.githubusercontent.com/u/1678258#1\n" +
			"   https://avatars1.githubusercontent.com/u/1678258#2\n" +
			"\t\t\thttps://avatars1.githubusercontent.com/u/1678258#3",
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
			assertEqualTo(t, node.GetName(), "simple_text")
			assertEqualTo(t, node.GetAttribute("indent"), []string{fmt.Sprintf("%d", fixture.indent[i])})

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
		indent   []int
		expected [][]string
	}{
		{"github.com/ohtomi/scrapbox",
			[]int{0},
			[][]string{
				{"github.com/ohtomi/scrapbox"},
			},
		},
		{"   github.com/ohtomi/scrapbox",
			[]int{3},
			[][]string{
				{"github.com/ohtomi/scrapbox"},
			},
		},
		{"\t\t\tgithub.com/ohtomi/scrapbox",
			[]int{3},
			[][]string{
				{"github.com/ohtomi/scrapbox"},
			},
		},
		{"github.com/ohtomi/scrapbox/1\n" +
			"   github.com/ohtomi/scrapbox/2\n" +
			"\t\t\tgithub.com/ohtomi/scrapbox/3",
			[]int{0, 3, 3},
			[][]string{
				{"github.com/ohtomi/scrapbox/1"},
				{"github.com/ohtomi/scrapbox/2"},
				{"github.com/ohtomi/scrapbox/3"},
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
			assertEqualTo(t, node.GetName(), "simple_text")
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
