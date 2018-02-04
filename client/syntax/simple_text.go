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

	head := node.GetChildren()[0]
	if head.GetName() == "indent" {
		attributes["indent"] = []string{fmt.Sprintf("%d", len(head.GetValue()))}
	} else if head.GetName() == "missing" {
		attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
	}

	return &SimpleText{parsec.NonTerminal{Name: "simple_text", Children: children, Attributes: attributes}}
}
