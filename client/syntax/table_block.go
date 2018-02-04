package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type TableBlock struct {
	parsec.NonTerminal
}

func NewTableBlock(node parsec.Queryable) *TableBlock {
	children := node.GetChildren()[1:]
	attributes := map[string][]string{}

	head := node.GetChildren()[0]
	if head.GetName() == "indent" {
		attributes["indent"] = []string{fmt.Sprintf("%d", len(head.GetValue()))}
	} else if head.GetName() == "missing" {
		attributes["indent"] = []string{fmt.Sprintf("%d", 0)}
	}

	return &TableBlock{parsec.NonTerminal{Name: "table_block", Children: children, Attributes: attributes}}
}
