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

	head := node.GetChildren()[0]
	if head.GetName() == "indent" {
		attributes["indent"] = []string{fmt.Sprintf("%d", len(head.GetValue()))}
	} else if head.GetName() == "missing" {
		attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
	}

	return &CodeBlock{parsec.NonTerminal{Name: "code_block", Children: children, Attributes: attributes}}
}
