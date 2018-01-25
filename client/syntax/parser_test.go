package syntax

import (
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
		expected []string
	}{
		{"   ", []string{"   "}},
		{"\t\t\t", []string{"\t\t\t"}},
	} {
		queryable, remaining := Parse([]byte(fixture.original), enablePrettyPrint)

		if len(remaining) != 0 {
			t.Fatalf("Got %q, but Want %q", string(remaining), "")
		}
		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected)+1 {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected)+1, queryable)
		}

		for i, expected := range fixture.expected {
			assertEqualTo(t, queryable.GetChildren()[i].GetName(), "indent")
			assertEqualTo(t, queryable.GetChildren()[i].GetValue(), expected)
		}
	}
}

func TestParse__quoted_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		expected []string
	}{
		{
			">https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox",
			[]string{"https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox"},
		},
		{
			"   >https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox",
			[]string{"https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox"},
		},
		{
			"\t\t\t>https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox",
			[]string{"https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox"},
		},
	} {
		queryable, remaining := Parse([]byte(fixture.original), enablePrettyPrint)

		if len(remaining) != 0 {
			t.Fatalf("Got %q, but Want %q", string(remaining), fixture.original)
		}
		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected)+1 {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected)+1, queryable)
		}

		for i, expected := range fixture.expected {
			assertEqualTo(t, queryable.GetChildren()[i+1].GetName(), "quoted")
			if len(queryable.GetChildren()[i+1].GetChildren()) != 2 {
				t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()[i+1].GetChildren()), 2, queryable.GetChildren()[i+1])
			}
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[0].GetName(), "q")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[0].GetValue(), ">")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[1].GetName(), "t")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[1].GetValue(), expected)
		}
	}
}

func TestParse__code_directive_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		expected []string
	}{
		{"code:sample.js", []string{"sample.js"}},
		{"   code:sample.js", []string{"sample.js"}},
		{"\t\t\tcode:sample.js", []string{"sample.js"}},
	} {
		queryable, remaining := Parse([]byte(fixture.original), enablePrettyPrint)

		if len(remaining) != 0 {
			t.Fatalf("Got %q, but Want %q", string(remaining), "")
		}
		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected)+1 {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected)+1, queryable)
		}

		for i, expected := range fixture.expected {
			assertEqualTo(t, queryable.GetChildren()[i+1].GetName(), "code")
			if len(queryable.GetChildren()[i+1].GetChildren()) != 2 {
				t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()[i+1].GetChildren()), 2, queryable.GetChildren()[i+1])
			}
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[0].GetName(), "c")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[0].GetValue(), "code:")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[1].GetName(), "n")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[1].GetValue(), expected)
		}
	}
}

func TestParse__table_directive_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		expected []string
	}{
		{"table:sample.js", []string{"sample.js"}},
		{"   table:sample.js", []string{"sample.js"}},
		{"\t\t\ttable:sample.js", []string{"sample.js"}},
	} {
		queryable, remaining := Parse([]byte(fixture.original), enablePrettyPrint)

		if len(remaining) != 0 {
			t.Fatalf("Got %q, but Want %q", string(remaining), "")
		}
		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected)+1 {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected)+1, queryable)
		}

		for i, expected := range fixture.expected {
			assertEqualTo(t, queryable.GetChildren()[i+1].GetName(), "table")
			if len(queryable.GetChildren()[i+1].GetChildren()) != 2 {
				t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()[i+1].GetChildren()), 2, queryable.GetChildren()[i+1])
			}
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[0].GetName(), "t")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[0].GetValue(), "table:")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[1].GetName(), "n")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetChildren()[1].GetValue(), expected)
		}
	}
}

func TestParse__image_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		expected []string
	}{
		{"https://avatars1.githubusercontent.com/u/1678258#.png", []string{"https://avatars1.githubusercontent.com/u/1678258#.png"}},
		{"   https://avatars1.githubusercontent.com/u/1678258#.png", []string{"https://avatars1.githubusercontent.com/u/1678258#.png"}},
		{"\t\t\thttps://avatars1.githubusercontent.com/u/1678258#.png", []string{"https://avatars1.githubusercontent.com/u/1678258#.png"}},
	} {
		queryable, remaining := Parse([]byte(fixture.original), enablePrettyPrint)

		if len(remaining) != 0 {
			t.Fatalf("Got %q, but Want %q", string(remaining), "")
		}
		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected)+1 {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected)+1, queryable)
		}

		for i, expected := range fixture.expected {
			assertEqualTo(t, queryable.GetChildren()[i+1].GetName(), "image")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetValue(), expected)
		}
	}
}

func TestParse__url_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		expected []string
	}{
		{"https://avatars1.githubusercontent.com/u/1678258", []string{"https://avatars1.githubusercontent.com/u/1678258"}},
		{"   https://avatars1.githubusercontent.com/u/1678258", []string{"https://avatars1.githubusercontent.com/u/1678258"}},
		{"\t\t\thttps://avatars1.githubusercontent.com/u/1678258", []string{"https://avatars1.githubusercontent.com/u/1678258"}},
	} {
		queryable, remaining := Parse([]byte(fixture.original), enablePrettyPrint)

		if len(remaining) != 0 {
			t.Fatalf("Got %q, but Want %q", string(remaining), "")
		}
		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected)+1 {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected)+1, queryable)
		}

		for i, expected := range fixture.expected {
			assertEqualTo(t, queryable.GetChildren()[i+1].GetName(), "url")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetValue(), expected)
		}
	}
}

func TestParse__text_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		expected []string
	}{
		{"github.com/ohtomi/scrapbox", []string{"github.com/ohtomi/scrapbox"}},
		{"   github.com/ohtomi/scrapbox", []string{"github.com/ohtomi/scrapbox"}},
		{"\t\t\tgithub.com/ohtomi/scrapbox", []string{"github.com/ohtomi/scrapbox"}},
	} {
		queryable, remaining := Parse([]byte(fixture.original), enablePrettyPrint)

		if len(remaining) != 0 {
			t.Fatalf("Got %q, but Want %q", string(remaining), fixture.original)
		}
		if queryable == nil {
			t.Fatalf("Failed to parse")
		}
		if len(queryable.GetChildren()) != len(fixture.expected)+1 {
			t.Fatalf("Found %d, but Want %d: %+v", len(queryable.GetChildren()), len(fixture.expected)+1, queryable)
		}

		for i, expected := range fixture.expected {
			assertEqualTo(t, queryable.GetChildren()[i+1].GetName(), "text")
			assertEqualTo(t, queryable.GetChildren()[i+1].GetValue(), expected)
		}
	}
}

func assertEqualTo(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Got %+v, but Want %+v", actual, expected)
	}
}
