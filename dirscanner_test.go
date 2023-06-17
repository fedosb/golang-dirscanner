package golang_dirscanner

import (
	"context"
	"fmt"
	"golang-dirscanner/pkg/dirscanner"
	"strings"
	"testing"
)

func printTree(node *dirscanner.Node, indent, prefix, path string) {
	if node == nil {
		return
	}

	fmt.Println(indent+prefix, node.Type, strings.ReplaceAll(node.Name, path, ""),
		fmt.Sprintf("(%d B; %.4f %%)", node.Size, node.Percentage),
	)
	for _, child := range node.Children {
		printTree(child, indent+" │  ", " ├", node.Name+"/")
	}
}

func TestDirScanner_Scan(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := dirscanner.NewScanner(1)
	tree := s.Scan("./", ctx)

	printTree(tree, "", "", "")

}
