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

	if len(node.GetChildren()) > 0 {
		child := node.GetChildren()[0]
		if child.GetName() == "indent" {
			attributes["indent"] = []string{fmt.Sprintf("%d", len(child.GetValue()))}
		} else if child.GetName() == "missing" {
			attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
		}
	}

	return &QuotedText{parsec.NonTerminal{Name: "quoted_text", Children: children, Attributes: attributes}}
}
