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
	queryable, remaining := Parse([]byte("   "), enablePrettyPrint)

	if len(remaining) != 0 {
		t.Fatalf("Got %q, but Want %q", string(remaining), "")
	}
	if queryable == nil {
		t.Fatalf("Failed to parse")
	}
	if len(queryable.GetChildren()) == 0 {
		t.Fatalf("Not found children: %+v", queryable)
	}

	assertEqualTo(t, queryable.GetChildren()[0].GetName(), "indent")
	assertEqualTo(t, queryable.GetChildren()[0].GetValue(), "   ")
}

func TestParse__image_node(t *testing.T) {
	queryable, remaining := Parse([]byte("https://avatars1.githubusercontent.com/u/1678258#.png"), enablePrettyPrint)

	if len(remaining) != 0 {
		t.Fatalf("Got %q, but Want %q", string(remaining), "")
	}
	if queryable == nil {
		t.Fatalf("Failed to parse")
	}
	if len(queryable.GetChildren()) == 0 {
		t.Fatalf("Not found children: %+v", queryable)
	}

	assertEqualTo(t, queryable.GetChildren()[1].GetName(), "image")
	assertEqualTo(t, queryable.GetChildren()[1].GetValue(), "https://avatars1.githubusercontent.com/u/1678258#.png")
}

func TestParse__url_node(t *testing.T) {
	queryable, remaining := Parse([]byte("https://avatars1.githubusercontent.com/u/1678258"), enablePrettyPrint)

	if len(remaining) != 0 {
		t.Fatalf("Got %q, but Want %q", string(remaining), "")
	}
	if queryable == nil {
		t.Fatalf("Failed to parse")
	}
	if len(queryable.GetChildren()) == 0 {
		t.Fatalf("Not found children: %+v", queryable)
	}

	assertEqualTo(t, queryable.GetChildren()[1].GetName(), "url")
	assertEqualTo(t, queryable.GetChildren()[1].GetValue(), "https://avatars1.githubusercontent.com/u/1678258")
}

func TestParse__text_node(t *testing.T) {
	queryable, remaining := Parse([]byte("github.com/ohtomi/scrapbox"), enablePrettyPrint)

	if len(remaining) != 0 {
		t.Fatalf("Got %q, but Want %q", string(remaining), "")
	}
	if queryable == nil {
		t.Fatalf("Failed to parse")
	}
	if len(queryable.GetChildren()) == 0 {
		t.Fatalf("Not found children: %+v", queryable)
	}

	assertEqualTo(t, queryable.GetChildren()[1].GetName(), "text")
	assertEqualTo(t, queryable.GetChildren()[1].GetValue(), "github.com/ohtomi/scrapbox")
}

func assertEqualTo(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Got %+v, but Want %+v", actual, expected)
	}
}
