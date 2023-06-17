package dirscanner

type IScanner interface {
	Scan(path string) *Node
}
