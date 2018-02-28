package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type TableBlock struct {
	parsec.NonTerminal
}

func NewTableBlock(node parsec.Queryable) *TableBlock {
	indent := node.GetChildren()[0]
	attributes := map[string][]string{}
	if indent.GetName() == "ws" {
		attributes["indent"] = []string{fmt.Sprintf("%d", len(indent.GetValue()))}
	} else if indent.GetName() == "missing" {
		attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
	}

	rest := node.GetChildren()[2]
	children := rest.GetChildren()

	return &TableBlock{parsec.NonTerminal{Name: "table_block", Children: children, Attributes: attributes}}
}
