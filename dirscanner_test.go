package golang_dirscanner

import (
	"fmt"
	"golang-dirscanner/pkg/dirscanner"
	"testing"
)

func TestDirScanner_Scan(t *testing.T) {

	scanner := dirscanner.NewScanner()
	tree := scanner.Scan("./")
	fmt.Println(tree)
}
