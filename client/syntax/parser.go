package syntax

import (
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
	text := parsec.Token("[^\n]+", "text")

	token := ast.OrdChoice("xx", nil, image, url, text)
	rest := ast.ManyUntil("rest", nil, token, end)

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
	// #text
	// #[text+]
	// `text+`

	callback := func(name string, s parsec.Scanner, node parsec.Queryable) parsec.Queryable {
		tokens := node.GetChildren()[0]
		children := tokens.GetChildren()
		body := children[1]

		switch body.GetName() {
		case "quoted":
			return NewQuotedText(tokens)
		case "code":
			return NewCodeBlock(tokens)
		case "table":
			return NewTableBlock(tokens)
		default:
			return NewSimpleText(tokens)
		}
	}

	tokens := ast.And("tokens", nil, indent, head, rest)
	line := ast.ManyUntil("line", callback, tokens, end)
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
