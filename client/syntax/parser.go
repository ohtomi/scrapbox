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

	indent := parsec.Token("[ \t]+", "indent")

	quoted := ast.And("quoted", nil, parsec.Atom(">", "q"), parsec.Token(".+", "t"))
	code := ast.And("code", nil, parsec.Atom("code:", "c"), parsec.Token(".+", "n"))
	table := ast.And("table", nil, parsec.Atom("table:", "t"), parsec.Token(".+", "n"))

	image := parsec.Token("(https://gyazo.com/[^ \t]+)|https?://[^ \t]+(\\.png|\\.gif|\\.jpg|\\.jpeg)", "image")
	url := parsec.Token("https?://[^ \t]+", "url")
	text := parsec.Token(".+", "text")

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

	Y := ast.And("Y", nil,
		ast.Maybe("maybe", nil, indent),
		ast.Maybe("maybe", nil, ast.OrdChoice("or", nil, quoted, code, table, image, url, text)),
	)

	return AST{ast: ast, parser: Y}
}

func Parse(line []byte, debug bool) parsec.Queryable {
	ast := NewAST()
	scanner := parsec.NewScanner(line).SetWSPattern("^[\r\n]+")
	queryable, _ := ast.ast.Parsewith(ast.parser, scanner)

	if debug {
		ast.ast.Prettyprint()
	}

	return queryable
}
