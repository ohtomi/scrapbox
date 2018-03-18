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

	project_link := parsec.Token("\\[(/[^/ \n]+)+\\]", "link")
	image_link1 := parsec.Token("\\[(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))[ \t]+https?://[^ \t\n]+\\]", "link")
	image_link2 := parsec.Token("\\[https?://[^ \t\n]+[ \t]+(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))\\]", "link")
	labeled_link1 := parsec.Token("\\[[^[\n]+[ \t]+https?://[^ \t\n]+\\]", "link")
	labeled_link2 := parsec.Token("\\[https?://[^ \t\n]+[ \t]+[^[\n]+\\]", "link")
	external_link := parsec.Token("\\[https?://[^[ \t\n]+\\]", "link")
	internal_link := parsec.Token("\\[[^[\n]*?\\]", "link")
	link := ast.OrdChoice("link", nil, project_link, image_link1, image_link2, labeled_link1, labeled_link2, external_link, internal_link)

	image := parsec.Token("(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))", "image")
	url := parsec.Token("https?://[^ \t\n]+", "url")
	tag := parsec.Token("#(\\[[^[\n]+\\]|[^ \t\n]+)", "tag")
	text := parsec.Token("[^\n]+", "text")

	bold_styled_url := parsec.Token("\\[\\[[*/\\-_]+[ \t]+https?://[^ \t\n]*?\\]\\]", "bold")
	bold_styled_text := parsec.Token("\\[\\[[*/\\-_]+[ \t]+[^\n]*?\\]\\]", "bold")
	bold_image := parsec.Token("\\[\\[(https://gyazo.com/[^ \t\n]+|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg))?\\]\\]", "bold")
	bold_text := parsec.Token("\\[\\[[^\n]*?\\]\\]", "bold")
	bold := ast.OrdChoice("bold", nil, bold_styled_url, bold_styled_text, bold_image, bold_text)

	token := ast.OrdChoice("token", nil, link, image, url, bold, tag, text)
	rest := ast.Kleene("rest", nil, token)

	// [text+]					-> internal_link
	// [url]					-> external_link
	// [text url]				-> labeled_link1
	// [url text]				-> labeled_link2
	// [image url]				-> image_link1
	// [url image]				-> image_link2
	// [/text(/text)*]			-> project_link
	// [text.icon]
	// [/text(/text)*.icon]
	// [[text]]					-> bold_text
	// [[image]]				-> bold_image
	// [[*/-_]+ url]
	// [[*/-_]+ text]			-> bold_styled_text
	// [$ text]
	// #text 					-> tag
	// #[text+] 				-> tag
	// `text+`

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
