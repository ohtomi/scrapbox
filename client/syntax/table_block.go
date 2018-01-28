package syntax

import (
	"fmt"

	"github.com/prataprc/goparsec"
)

type TableBlock struct {
	parsec.NonTerminal
}

func NewTableBlock(node parsec.Queryable) *TableBlock {
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

	return &TableBlock{parsec.NonTerminal{Name: "table_block", Children: children, Attributes: attributes}}
}
