package main

import (
	"flag"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
	"os"
)

var (
	pathPrefixFlag = flag.String("output-path", "", "Base path")
	filenameFlag   = flag.String("output-file", "tree.csv", "Output filename")
	numNodesFlag   = flag.Int("nodes", 10, "Number of nodes")
	depthFlag      = flag.Int("depth", 3, "Tree depth")
	formatFlag     = flag.String("format", "csv", "Output format: csv, vegeta")
)

func main() {
	flag.Parse()

	basePath, _ := os.Getwd()
	if *pathPrefixFlag != "" {
		basePath = *pathPrefixFlag
	}
	numNodes := *numNodesFlag
	depth := *depthFlag
	filename := *filenameFlag
	format := *formatFlag

	gen := generator.NewNodeGenerator(numNodes, depth)
	nodes := gen.GenerateRoot()
	generator.SaveNodesToFile(nodes, basePath, filename, format)
}
