package dirscanner

import (
	"context"
	"golang-dirscanner/pkg/queue"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type dirScanner struct {
	queue                  queue.Queue[*Node]
	sem                    semaphore.Weighted
	totalDirSize           uint64
	maxGoroutinesCount     int64
	currentGoroutinesCount atomic.Int32
}

func NewScanner(cnt int64) IScanner {

	if cnt < 1 {
		panic("INVALID GOROUTINES MAX COUNT")
	}

	return &dirScanner{
		maxGoroutinesCount: cnt,
		sem:                *semaphore.NewWeighted(cnt),
	}
}

func (s *dirScanner) Scan(path string, ctx context.Context) *Node {
	absPath, _ := filepath.Abs(path)
	tree := s.buildTree(absPath, ctx)
	s.traverseTreeDFS(tree)
	return tree
}

func (s *dirScanner) buildTree(rootPath string, ctx context.Context) *Node {

	wg := sync.WaitGroup{}

	root := &Node{
		Name:     rootPath,
		Type:     FileTypeDir,
		Children: make([]*Node, 0),
	}

	s.queue.Push(root)

bfsLoop:
	for s.queue.Size() > 0 || s.currentGoroutinesCount.Load() > 0 {
		select {
		case <-ctx.Done():
			break bfsLoop
		case <-time.After(time.Microsecond):
			if s.queue.Size() == 0 {
				continue bfsLoop
			}
		}

		currentNode := s.queue.Pop()

		wg.Add(1)
		_ = s.sem.Acquire(ctx, 1)

		s.currentGoroutinesCount.Add(1)
		go s.processDirNode(&wg, currentNode)
	}

	wg.Wait()

	return root
}

func (s *dirScanner) processDirNode(wg *sync.WaitGroup, currentNode *Node) {
	defer wg.Done()
	defer s.sem.Release(1)
	defer s.currentGoroutinesCount.Add(-1)

	files, err := ioutil.ReadDir(currentNode.Name)
	if err != nil {
		log.Print(err)
	}

	for _, file := range files {
		filePath := filepath.Join(currentNode.Name, file.Name())

		child := &Node{
			Name:     filePath,
			Size:     uint64(file.Size()),
			Children: make([]*Node, 0),
		}

		switch true {
		case file.Mode()&os.ModeSymlink == os.ModeSymlink:
			child.Type = FileTypeSymlink
			continue
		case file.Mode()&os.ModeDevice == os.ModeDevice:
			child.Type = FileTypeDevice
			continue
		case file.Mode()&os.ModeDir == os.ModeDir:
			child.Type = FileTypeDir
		default:
			child.Type = FileTypeFile
		}

		s.totalDirSize += child.Size

		currentNode.Children = append(currentNode.Children, child)

		if file.IsDir() {
			s.queue.Push(child)
		}
	}
}

func (s *dirScanner) traverseTreeDFS(cur *Node) {
	var wg = sync.WaitGroup{}

	for i := range cur.Children {

		acquired := s.sem.TryAcquire(1)

		if acquired {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				defer s.sem.Release(1)

				s.traverseTreeDFS(cur.Children[i])
				cur.Size += cur.Children[i].Size
			}(i)
		} else {
			s.traverseTreeDFS(cur.Children[i])
			cur.Size += cur.Children[i].Size
		}
	}

	wg.Wait()

	cur.Percentage = float64(cur.Size) / float64(s.totalDirSize) * 100
}
