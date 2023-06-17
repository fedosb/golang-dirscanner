# golang-dirscanner
### Multithreaded utility for analyzing the ratio of file sizes and directories within a selected directory

## Analysis of File and Directory Sizes

This project provides a tool for analyzing the sizes of files and directories. The analysis is performed in a multithreaded mode using the breadth-first search (BFS) algorithm.

## Features

- Analysis of file and directory sizes is performed in a multithreaded mode using a breadth-first search (BFS) algorithm.
- Each directory is processed in a separate thread. Processing includes building a tree of directories and files for further processing.
- Maximum thread limit: The number of threads involved in the analysis is limited.
- The analysis can be canceled via context.

### NB:

The thread processing a directory does not wait for the processing of all nested directories; it simply queues them for processing. Otherwise, in cases with high levels of nesting, the threads would idle while waiting for the completion of threads launched for the nested directories.

Additionally, when setting the limit for the number of active threads to a value lower than the nesting level of directories in the scanned directory, in order to avoid program "hang-ups" due to mutual blocking (when all threads are occupied waiting for the completion of traversal in nested directories, and no threads are left to process them), the recalculation of directory sizes considering nested directories is separately performed using depth-first search (DFS), after the sizes of all files have been calculated.

## Usage Example (dirscanner_test.go):

```go
// printTree is a utility function used to print the directory tree in a visually appealing way.
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

// TestDirScanner_Scan is a test function that demonstrates the usage of the dirscanner library.
func TestDirScanner_Scan(t *testing.T) {
    // Set up a context and defer its cancellation to ensure clean-up.
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Create a new scanner with a concurrency level of 1.
    s := dirscanner.NewScanner(1)

    // Scan the specified directory and get the resulting tree.
    tree := s.Scan("./", ctx)

    // Print the directory tree by calling printTree with appropriate arguments.
    printTree(tree, "", "", "")
}
```

### Output:

```text
DIR /Users/fedosb/golang-dirscanner (63250 B; 100.0000 %)
│   ├ DIR .git (51104 B; 80.7968 %)
 ... (the .git folder is omitted in this example due to lengthy output)
│   ├ FILE .gitignore (485 B; 0.7668 %)
│   ├ DIR .idea (6587 B; 10.4142 %)
 ... (the .idea folder is also omitted)
│   ├ FILE README.md (131 B; 0.2071 %)
│   ├ FILE dirscanner_test.go (671 B; 1.0609 %)
│   ├ FILE go.mod (68 B; 0.1075 %)
│   ├ FILE go.sum (153 B; 0.2419 %)
│   ├ DIR pkg (4051 B; 6.4047 %)
│   │   ├ DIR dirscanner (3475 B; 5.4941 %)
│   │   │   ├ FILE dirscanner.go (2934 B; 4.6387 %)
│   │   │   ├ FILE interface.go (103 B; 0.1628 %)
│   │   │   ├ FILE node.go (278 B; 0.4395 %)
│   │   ├ DIR queue (448 B; 0.7083 %)
│   │   │   ├ FILE queue.go (352 B; 0.5565 %)
```