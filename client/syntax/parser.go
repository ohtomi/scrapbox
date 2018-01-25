package syntax

import (
	"github.com/prataprc/goparsec"
)

var (
	ast = parsec.NewAST("ast", 1000)

	indent = parsec.Token("[ \t]+", "indent")

	quoted = ast.And("quoted", nil, parsec.Atom(">", "q"), parsec.Token(".+", "t"))

	image  = parsec.Token("(https://gyazo.com/[^ \t]+)|https?://[^ \t]+(\\.png|\\.gif|\\.jpg|\\.jpeg)", "image")
	url    = parsec.Token("https?://[^ \t]+", "url")
	text   = parsec.Token(".+", "text")

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
	// code:text
	// table:text

	Y = ast.And("Y", nil,
		ast.Maybe("maybe", nil, indent),
		ast.Maybe("maybe", nil, ast.OrdChoice("or", nil, quoted, image, url, text)),
	)
)

func Parse(line []byte, debug bool) (parsec.Queryable, []byte) {
	ast.Reset()
	scanner := parsec.NewScanner(line).SetWSPattern("^[\r\n]+")
	queryable, scanner := ast.Parsewith(Y, scanner)
	remaining, _ := scanner.Match(".+")

	if debug {
		ast.Prettyprint()
	}
	return queryable, remaining
}
