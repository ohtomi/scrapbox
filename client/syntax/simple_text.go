package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type SimpleText struct {
	parsec.NonTerminal
}

func NewSimpleText(node parsec.Queryable) *SimpleText {
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

	return &SimpleText{parsec.NonTerminal{Name: "simple_text", Children: children, Attributes: attributes}}
}
