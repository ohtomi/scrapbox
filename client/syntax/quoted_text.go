package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type QuotedText struct {
	parsec.NonTerminal
}

func NewQuotedText(node parsec.Queryable) *QuotedText {
	children := node.GetChildren()[1:]
	attributes := map[string][]string{}

	head := node.GetChildren()[0]
	if head.GetName() == "indent" {
		attributes["indent"] = []string{fmt.Sprintf("%d", len(head.GetValue()))}
	} else if head.GetName() == "missing" {
		attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
	}

	return &QuotedText{parsec.NonTerminal{Name: "quoted_text", Children: children, Attributes: attributes}}
}
