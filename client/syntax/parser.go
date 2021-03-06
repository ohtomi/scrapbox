package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type AST struct {
	ast    *parsec.AST
	parser parsec.Parser
}

func NewAST() AST {
	ast := parsec.NewAST("ast", 1000)

	lf := parsec.Token("\n", "lf")
	end := parsec.End()

	ws := parsec.Token("[ \t]+", "ws")

	indent := ast.Maybe("indent", nil, ws)

	quoted := parsec.Atom(">", "quoted")
	code := parsec.Atom("code:", "code")
	table := parsec.Atom("table:", "table")

	mark := ast.OrdChoice("mark", nil, quoted, code, table)
	head := ast.Maybe("head", nil, mark)

	// [$ text]
	math := parsec.Token("\\[\\$[ \t]+[^\n]*?\\]", "math")
	// [[*/-_]+ url]
	styled_url := parsec.Token("\\[[*/\\-_]+[ \t]+https?://[^ \t\n]*?\\]", "styled_url")
	// [[*/-_]+ text]
	styled_text := parsec.Token("\\[[*/\\-_]+[ \t]+[^\n]*?\\]", "styled_text")
	// [/text(/text)*]
	project_link := parsec.Token("\\[(/[^/ \n]+)+\\]", "project_link")
	// [image url]
	image_link1 := parsec.Token("\\[(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))[ \t]+https?://[^ \t\n]+\\]", "image_link1")
	// [url image]
	image_link2 := parsec.Token("\\[https?://[^ \t\n]+[ \t]+(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))\\]", "image_link2")
	// [text url]
	labeled_link1 := parsec.Token("\\[[^[\n]+[ \t]+https?://[^ \t\n]+\\]", "labeled_link1")
	// [url text]
	labeled_link2 := parsec.Token("\\[https?://[^ \t\n]+[ \t]+[^[\n]+\\]", "labeled_link2")
	// [url]
	external_link := parsec.Token("\\[https?://[^[ \t\n]+\\]", "external_link")
	// [/text(/text)*.icon]
	page_icon := parsec.Token("\\[(/[^\n]+)+\\.icon\\]", "page_icon")
	// [text.icon]
	icon := parsec.Token("\\[[^\n]+\\.icon\\]", "icon")
	// [text+]
	internal_link := parsec.Token("\\[[^[\n]*?\\]", "internal_link")
	// image
	image := parsec.Token("(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))", "image")
	// url
	url := parsec.Token("https?://[^ \t\n]+", "url")
	// [[image]]
	bold_image := parsec.Token("\\[\\[(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))?\\]\\]", "bold_image")
	// [[text]]
	bold_text := parsec.Token("\\[\\[[^\n]*?\\]\\]", "bold_text")
	// `text+`
	snippet := parsec.Token("`[^`]*?`", "snippet")
	// #[text( text)*] | #text
	tag := parsec.Token("#(\\[[^[\n]+\\]|[^ \t\n]+)", "tag")
	// text
	text := parsec.Token("[^\n]+", "text")

	token := ast.OrdChoice("token", nil,
		math,
		styled_url,
		styled_text,
		project_link,
		image_link1,
		image_link2,
		labeled_link1,
		labeled_link2,
		external_link,
		page_icon,
		icon,
		internal_link,
		image,
		url,
		bold_image,
		bold_text,
		snippet,
		tag,
		text)
	rest := ast.Kleene("rest", nil, token)

	callback := func(name string, s parsec.Scanner, node parsec.Queryable) parsec.Queryable {
		indent := node.GetChildren()[0]
		attributes := map[string][]string{}
		if indent.GetName() == "ws" {
			attributes["indent"] = []string{fmt.Sprintf("%d", len(indent.GetValue()))}
		} else if indent.GetName() == "missing" {
			attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
		}

		head := node.GetChildren()[1]
		var newName string
		switch head.GetName() {
		case "quoted":
			newName = "quoted_text"
		case "code":
			newName = "code_block"
		case "table":
			newName = "table_block"
		default:
			newName = "simple_text"
		}

		rest := node.GetChildren()[2]
		children := rest.GetChildren()

		return &parsec.NonTerminal{Name: newName, Children: children, Attributes: attributes}
	}

	line := ast.And("line", callback, indent, head, rest)
	root := ast.ManyUntil("root", nil, line, lf, end)

	return AST{ast: ast, parser: root}
}

func Parse(contents []byte, debug bool) parsec.Queryable {
	ast := NewAST()
	scanner := parsec.NewScanner(contents).SetWSPattern("\r\n")
	queryable, _ := ast.ast.Parsewith(ast.parser, scanner)

	if debug {
		if queryable != nil {
			ast.ast.Prettyprint()
		}
	}

	return queryable
}
