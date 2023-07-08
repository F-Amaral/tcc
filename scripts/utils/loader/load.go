package main

import (
	"flag"
	"fmt"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
	vegeta "github.com/F-Amaral/tcc/scripts/utils/vegeta/impl"
	"path"
)

type Flags struct {
	TargetPath    string
	StartNumNodes int
	MaxNumNodes   int
	StepNumNodes  int
	Depth         int
	Width         int
	OutputFolder  string
	Format        string
	DepthPriority bool
	exec          bool
}

func parseFlags() *Flags {
	flags := &Flags{}
	flag.IntVar(&flags.StartNumNodes, "n-start", 10, "Start number of nodes")
	flag.IntVar(&flags.MaxNumNodes, "n-max", 10, "Max number of nodes")
	flag.IntVar(&flags.StepNumNodes, "n-step", 10, "Step number of nodes")
	flag.IntVar(&flags.Depth, "depth", 3, "Depth side of ratio with width")
	flag.IntVar(&flags.Width, "width", 3, "width side of ratio with depth")
	flag.StringVar(&flags.OutputFolder, "output-folder", "./output-local", "Output folder")
	flag.StringVar(&flags.TargetPath, "target-path", "http://test.pi.hole:8080", "Target path")
	flag.StringVar(&flags.Format, "format", "csv", "Output format: path, csv")
	flag.BoolVar(&flags.DepthPriority, "depth-priority", true, "Depth priority")
	flag.BoolVar(&flags.exec, "exec", false, "Execute load")
	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()
	Load(flags)
}

func Load(flags *Flags) {

	numIterations := (flags.MaxNumNodes - flags.StartNumNodes) / flags.StepNumNodes
	dataFiles := []string{}
	for i := 0; i <= numIterations; i++ {
		numNodes := flags.StartNumNodes + (i * flags.StepNumNodes)
		gen := generator.NewNodeGenerator(numNodes, flags.Depth, flags.Width)
		nodes := gen.GenerateRoot()

		//db, _ := json.Marshal(nodes)
		//fmt.Println(string(db))
		outFolder := path.Join(flags.OutputFolder, fmt.Sprintf("%d-%d-%d", numNodes, flags.Depth, flags.Width))
		datasetName := fmt.Sprintf("dataset_%d_%s.csv", numNodes, flags.Format)
		dataFiles = append(dataFiles, path.Join(outFolder, datasetName))
		generator.SaveNodesToFile(nodes, outFolder, datasetName, flags.Format)
		vegeta.Run(outFolder, "", "vegeta", "all", "http://localhost:8080", false)
	}

	if flags.exec {
		vegeta.Run(flags.OutputFolder, "", "output", "all", "http://test.pi.hole:8080", true)
	} else {
	}
}
