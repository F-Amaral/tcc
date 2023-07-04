package main

import (
	"flag"
	"fmt"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
	vegeta "github.com/F-Amaral/tcc/scripts/utils/vegeta/impl"
)

var (
	startNumNodesFlag = flag.Int("n-start", 10, "Start number of nodes")
	maxNumNodesFlag   = flag.Int("n-max", 100, "Max number of nodes")
	stepNumNodesFlag  = flag.Int("n-step", 10, "Step number of nodes")
	avgDepthFlag      = flag.Int("avg-depth", 3, "Average depth")
	avgDepthIncFlag   = flag.Int("avg-depth-inc", 1, "Average depth increment")
	maxDepthFlag      = flag.Int("max-depth", 10, "Maximum depth")
	probFlag          = flag.Float64("prob", 1, "Probability of child node creation")
	outputFolderFlag  = flag.String("output-folder", "./output", "Output folder")
)

func main() {
	flag.Parse()
	startNumNodes := *startNumNodesFlag
	maxNumNodes := *maxNumNodesFlag
	stepNumNodes := *stepNumNodesFlag
	avgDepth := *avgDepthFlag
	avgDepthInc := *avgDepthIncFlag
	maxDepth := *maxDepthFlag
	prob := *probFlag
	outputFolder := *outputFolderFlag
	Load(startNumNodes, maxNumNodes, stepNumNodes, avgDepth, avgDepthInc, maxDepth, prob, outputFolder)
}

func Load(startNumNodes, maxNumNodes, stepNumNodes, avgDepth, avgDepthInc, maxDepth int, prob float64, outputFolder string) {
	numIterations := (maxNumNodes - startNumNodes) / stepNumNodes
	for i := 0; i < numIterations; i++ {
		numNodes := startNumNodes + (i * stepNumNodes)
		avgDepth += i * avgDepthInc
		nodes := generator.Generate(numNodes, avgDepth, maxDepth, prob)
		datasetName := fmt.Sprintf("dataset_%d_%d_%d.csv", numNodes, avgDepth, maxDepth)
		generator.SaveNodesToFile(nodes, outputFolder, datasetName)
	}

	vegeta.Run(outputFolder, "", "input", "all", "http://localhost:8080", true)
}
