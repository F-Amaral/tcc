package main

import (
	"encoding/json"
	"flag"
	"fmt"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
	vegeta "github.com/F-Amaral/tcc/scripts/utils/vegeta/impl"
)

type Flags struct {
	StartNumNodes int
	MaxNumNodes   int
	StepNumNodes  int
	StartDepth    int
	StepDepth     int
	StartWidth    int
	StepWidth     int
	OutputFolder  string
	Format        string
	DepthPriority bool
}

func parseFlags() *Flags {
	flags := &Flags{}
	flag.IntVar(&flags.StartNumNodes, "n-start", 10, "Start number of nodes")
	flag.IntVar(&flags.MaxNumNodes, "n-max", 10, "Max number of nodes")
	flag.IntVar(&flags.StepNumNodes, "n-step", 10, "Step number of nodes")
	flag.IntVar(&flags.StartDepth, "depth-start", 3, "Maximum depth")
	flag.IntVar(&flags.StepDepth, "step-depth", 1, "Step depth")
	flag.IntVar(&flags.StartWidth, "width-start", 3, "Maximum Width")
	flag.IntVar(&flags.StepWidth, "step-width", 1, "Step Width")
	flag.StringVar(&flags.OutputFolder, "output-folder", "./output", "Output folder")
	flag.StringVar(&flags.Format, "format", "vegeta", "Output format: path, csv")
	flag.BoolVar(&flags.DepthPriority, "depth-priority", true, "Depth priority")
	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()
	Load(flags)
}

func Load(flags *Flags) {
	numIterations := (flags.MaxNumNodes - flags.StartNumNodes) / flags.StepNumNodes
	for i := 0; i <= numIterations; i++ {
		numNodes := flags.StartNumNodes + (i * flags.StepNumNodes)
		depth := flags.StartDepth + (i * flags.StepDepth)
		gen := generator.NewNodeGenerator(numNodes, depth, flags.StartWidth, flags.DepthPriority)
		nodes := gen.GenerateRoot()
		db, _ := json.Marshal(nodes)
		fmt.Println(string(db))
		datasetName := fmt.Sprintf("dataset_%d_%s.csv", numNodes, flags.Format)
		generator.SaveNodesToFile(nodes, flags.OutputFolder, datasetName, flags.Format)
	}

	switch flags.Format {
	case "vegeta":
		vegeta.Run(flags.OutputFolder, "", "input", "all", "http://test.pi.hole:8080", true)
	default:
		print("Creating dataset without execution")
	}
}
