package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type QuotedText struct {
	parsec.NonTerminal
}

func NewQuotedText(node parsec.Queryable) *QuotedText {
	children := node.GetChildren()
	attributes := map[string][]string{}

	for i, child := range children {
		if child.GetName() == "indent" {
			children = append(children[:i], children[i+1:]...)
			attributes["indent"] = []string{fmt.Sprintf("%d", len(child.GetValue()))}
		} else if child.GetName() == "missing" {
			children = append(children[:i], children[i+1:]...)
			attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
		}
	}

	return &QuotedText{parsec.NonTerminal{Name: "quoted_text", Children: children, Attributes: attributes}}
}
