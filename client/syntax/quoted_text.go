package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type QuotedText struct {
	parsec.NonTerminal
}

func NewQuotedText(node parsec.Queryable) *QuotedText {
	indent := node.GetChildren()[0]
	attributes := map[string][]string{}
	if indent.GetName() == "ws" {
		attributes["indent"] = []string{fmt.Sprintf("%d", len(indent.GetValue()))}
	} else if indent.GetName() == "missing" {
		attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
	}

	rest := node.GetChildren()[2]
	children := rest.GetChildren()

	return &QuotedText{parsec.NonTerminal{Name: "quoted_text", Children: children, Attributes: attributes}}
}
