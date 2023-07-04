package main

import (
	"flag"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
	"os"
)

var (
	pathPrefixFlag = flag.String("p", "", "Base path")
	filenameFlag   = flag.String("o", "tree.csv", "Output filename")
	numNodesFlag   = flag.Int("n", 100, "Number of nodes")
	avgDepthFlag   = flag.Int("avg-depth", 3, "Average depth")
	maxDepthFlag   = flag.Int("max-depth", 4, "Maximum depth")
	probFlag       = flag.Float64("prob", 0.6, "Probability")
)

func main() {
	flag.Parse()

	basePath, _ := os.Getwd()
	if *pathPrefixFlag != "" {
		basePath = *pathPrefixFlag
	}
	numNodes := *numNodesFlag
	avgDepth := *avgDepthFlag
	maxDepth := *maxDepthFlag
	prob := *probFlag
	filename := *filenameFlag
	nodes := generator.Generate(numNodes, avgDepth, maxDepth, prob)
	generator.SaveNodesToFile(nodes, basePath, filename)
}
