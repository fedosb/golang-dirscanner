package dirscanner

type Node struct {
	Name       string
	Type       NodeType
	Size       uint64
	Percentage float64
	Children   []*Node
}

type NodeType string

const (
	FileTypeDir     = "DIR"
	FileTypeDevice  = "DEV"
	FileTypeSymlink = "SYMLINK"
	FileTypeFile    = "FILE"
)
