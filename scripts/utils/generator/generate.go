package main

import (
	"flag"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
)

type Flags struct {
	PathPrefix    string
	Filename      string
	NumNodes      int
	Depth         int
	Width         int
	Format        string
	DepthPriority bool
}

func parseFlags() *Flags {
	flags := &Flags{}
	flag.StringVar(&flags.PathPrefix, "output-path", "", "Base path")
	flag.StringVar(&flags.Filename, "output-file", "tree.csv", "Output filename")
	flag.IntVar(&flags.NumNodes, "nodes", 20, "Number of nodes")
	flag.IntVar(&flags.Depth, "depth", 3, "Tree depth")
	flag.IntVar(&flags.Width, "width", 3, "Tree width")
	flag.StringVar(&flags.Format, "format", "csv", "Output format: csv, vegeta")
	flag.BoolVar(&flags.DepthPriority, "depth-priority", true, "Depth priority")
	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()

	gen := generator.NewNodeGenerator(flags.NumNodes, flags.Depth, flags.Width, false)
	nodes := gen.GenerateRoot()
	generator.SaveNodesToFile(nodes, flags.PathPrefix, flags.Filename, flags.Format)
}
