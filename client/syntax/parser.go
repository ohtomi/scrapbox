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

	image := parsec.Token("(https://gyazo.com/[^ \t\n]+)|https?://[^ \t\n]+(\\.png|\\.gif|\\.jpg|\\.jpeg)", "image")
	url := parsec.Token("https?://[^ \t\n]+", "url")
	tag := parsec.Token("#[^ \t\n]+", "tag")
	text := parsec.Token("[^\n]+", "text")

	token := ast.OrdChoice("token", nil, image, url, tag, text)
	rest := ast.ManyUntil("rest", nil, token, lf)

	// [text+]
	// [url]
	// [text url]
	// [url text]
	// [image url]
	// [url image]
	// [/text(/text)*]
	// [text.icon]
	// [/text(/text)*.icon]
	// [[text]]
	// [[image]]
	// [[*/-_]+ text]
	// [$ text]
	// #text -> tag
	// #[text+]
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
