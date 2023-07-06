package main

import (
	"encoding/json"
	"flag"
	"fmt"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
	vegeta "github.com/F-Amaral/tcc/scripts/utils/vegeta/impl"
)

var (
	startNumNodesFlag = flag.Int("n-start", 10, "Start number of nodes")
	maxNumNodesFlag   = flag.Int("n-max", 10, "Max number of nodes")
	stepNumNodesFlag  = flag.Int("n-step", 10, "Step number of nodes")
	startDepthFlag    = flag.Int("depth-start", 10, "Maximum depth")
	stepDepthFlag     = flag.Int("step-depth", 1, "Step depth")
	outputFolderFlag  = flag.String("output-folder", "./output", "Output folder")
	formatFlag        = flag.String("format", "vegeta", "Output format: path, csv")
)

func main() {
	flag.Parse()
	startNumNodes := *startNumNodesFlag
	maxNumNodes := *maxNumNodesFlag
	stepNumNodes := *stepNumNodesFlag
	startDepth := *startDepthFlag
	stepDetph := *stepDepthFlag
	outputFolder := *outputFolderFlag
	format := *formatFlag
	Load(startNumNodes, maxNumNodes, stepNumNodes, startDepth, stepDetph, outputFolder, format)
}

func Load(startNumNodes int, maxNumNodes int, stepNumNodes int, startDepth int, stepDepth int, outputFolder string, format string) {
	numIterations := (maxNumNodes - startNumNodes) / stepNumNodes
	for i := 0; i <= numIterations; i++ {
		numNodes := startNumNodes + (i * stepNumNodes)
		depth := startDepth + (i * stepDepth)
		gen := generator.NewNodeGenerator(numNodes, depth)
		nodes := gen.GenerateRoot()
		db, _ := json.Marshal(nodes)
		fmt.Println(string(db))
		datasetName := fmt.Sprintf("dataset_%d_%s.csv", numNodes, format)
		generator.SaveNodesToFile(nodes, outputFolder, datasetName, format)
	}

	switch format {
	case "vegeta":
		vegeta.Run(outputFolder, "", "input", "all", "http://test.pi.hole:8080", true)
		break
	default:
		print("Creating dataset without execution")
		break
	}
}
