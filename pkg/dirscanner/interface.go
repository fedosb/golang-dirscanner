package dirscanner

import "context"

type IScanner interface {
	Scan(string, context.Context) *Node
}
