package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type SimpleText struct {
	parsec.NonTerminal
}

func NewSimpleText(node parsec.Queryable, name string) *SimpleText {
	indent := node.GetChildren()[0]
	attributes := map[string][]string{}
	if indent.GetName() == "ws" {
		attributes["indent"] = []string{fmt.Sprintf("%d", len(indent.GetValue()))}
	} else if indent.GetName() == "missing" {
		attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
	}

	rest := node.GetChildren()[2]
	children := rest.GetChildren()

	return &SimpleText{parsec.NonTerminal{Name: name, Children: children, Attributes: attributes}}
}
