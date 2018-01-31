package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type SimpleText struct {
	parsec.NonTerminal
}

func NewSimpleText(node parsec.Queryable) *SimpleText {
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

	return &SimpleText{parsec.NonTerminal{Name: "simple_text", Children: children, Attributes: attributes}}
}
