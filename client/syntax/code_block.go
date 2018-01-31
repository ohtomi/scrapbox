package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type CodeBlock struct {
	parsec.NonTerminal
}

func NewCodeBlock(node parsec.Queryable) *CodeBlock {
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

	return &CodeBlock{parsec.NonTerminal{Name: "code_block", Children: children, Attributes: attributes}}
}
