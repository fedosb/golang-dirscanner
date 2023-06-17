package golang_dirscanner

import (
	"fmt"
	"golang-dirscanner/pkg/dirscanner"
	"strings"
	"testing"
)

func printTree(node *dirscanner.Node, indent, prefix, path string) {
	fmt.Println(indent+prefix, node.Type, strings.ReplaceAll(node.Name, path, ""),
		fmt.Sprintf("(%d B; %.4f %%)", node.Size, node.Percentage),
	)
	for _, child := range node.Children {
		printTree(child, indent+" │  ", " ├", node.Name+"/")
	}
}

func TestDirScanner_Scan(t *testing.T) {

	s := dirscanner.NewScanner(1)
	tree := s.Scan("./")

	printTree(tree, "", "", "")
}
