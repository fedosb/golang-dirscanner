package dirscanner

type dirScanner struct {
}

func NewScanner() IScanner {
	return &dirScanner{}
}

func (s *dirScanner) Scan(path string) *Node {
	return &Node{}
}
