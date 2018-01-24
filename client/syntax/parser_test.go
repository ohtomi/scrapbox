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
	} {
		queryable, remaining := Parse([]byte("   "), enablePrettyPrint)

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
			[]string{">https://avatars1.githubusercontent.com/u/1678258#.png https://avatars1.githubusercontent.com/u/1678258 github.com/ohtomi/scrapbox"},
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
			assertEqualTo(t, queryable.GetChildren()[i+1].GetValue(), expected)
		}
	}
}

func TestParse__image_node(t *testing.T) {
	for _, fixture := range []struct {
		original string
		expected []string
	}{
		{"https://avatars1.githubusercontent.com/u/1678258#.png", []string{"https://avatars1.githubusercontent.com/u/1678258#.png"}},
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
