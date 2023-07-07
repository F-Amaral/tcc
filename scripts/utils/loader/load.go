package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
	generator "github.com/F-Amaral/tcc/scripts/utils/generator/impl"
	"github.com/F-Amaral/tcc/scripts/utils/loader/impl"
	vegeta "github.com/tsenart/vegeta/lib"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
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
	decode        bool
	file          string
}

func parseFlags() *Flags {
	flags := &Flags{}
	flag.IntVar(&flags.StartNumNodes, "n-start", 10, "Start number of nodes")
	flag.IntVar(&flags.MaxNumNodes, "n-max", 10, "Max number of nodes")
	flag.IntVar(&flags.StepNumNodes, "n-step", 10, "Step number of nodes")
	flag.IntVar(&flags.Depth, "depth-start", 3, "Depth side of ratio with width")
	flag.IntVar(&flags.Width, "width-start", 3, "width side of ratio with depth")
	flag.StringVar(&flags.OutputFolder, "output-folder", "./output", "Output folder")
	flag.StringVar(&flags.TargetPath, "target-path", "http://test.pi.hole:8080", "Target path")
	flag.StringVar(&flags.Format, "format", "csv", "Output format: path, csv")
	flag.BoolVar(&flags.DepthPriority, "depth-priority", true, "Depth priority")
	flag.BoolVar(&flags.decode, "decode", false, "Decode")
	flag.StringVar(&flags.file, "file", "/Users/famaral/go/src/github.com/f-amaral/tcc/scripts/utils/loader/test.json", "File to decode")
	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()
	Load(flags)
}

func Load(flags *Flags) {
	if flags.decode {
		file, err := os.Open(flags.file)
		if err != nil {
			log.Fatal(err)
		}
		results := &vegeta.Result{}
		decoder := vegeta.NewJSONDecoder(file)
		fmt.Println(decoder.Decode(results))

	}

	numIterations := (flags.MaxNumNodes - flags.StartNumNodes) / flags.StepNumNodes
	dataFiles := []string{}
	for i := 0; i <= numIterations; i++ {
		numNodes := flags.StartNumNodes + (i * flags.StepNumNodes)
		gen := generator.NewNodeGenerator(numNodes, flags.Depth, flags.Width)
		nodes := gen.GenerateRoot()

		//db, _ := json.Marshal(nodes)
		//fmt.Println(string(db))
		datasetName := fmt.Sprintf("dataset_%d_%s.csv", numNodes, flags.Format)
		dataFiles = append(dataFiles, path.Join(flags.OutputFolder, datasetName))
		generator.SaveNodesToFile(nodes, flags.OutputFolder, datasetName, flags.Format)
	}

	for _, dataFile := range dataFiles {
		dataset, err := csv.ReadFromCSV(dataFile)
		if err != nil {
			fmt.Sprintf("Error reading file %s", dataFile)
			continue
		}

		addPptStep := impl.NewStep("Add To PPT", flags.TargetPath, impl.Withmethod(http.MethodPost), impl.WithPathFormat("/ppt/%s/%s"))
		addNestedStep := impl.NewStep("Add To Nested", flags.TargetPath, impl.Withmethod(http.MethodPost), impl.WithPathFormat("/nested/%s/%s"))
		getPptRecursiveStep := impl.NewStep("Get recursive", flags.TargetPath, impl.Withmethod(http.MethodGet), impl.WithPathFormat("/ppt/%s?recursive=true"))
		getPptStep := impl.NewStep("Get Ppt", flags.TargetPath, impl.Withmethod(http.MethodGet), impl.WithPathFormat("/ppt/%s?recursive=false"))
		getNestedStep := impl.NewStep("GetNested", flags.TargetPath, impl.Withmethod(http.MethodGet), impl.WithPathFormat("/nested/%s"))

		loader := impl.NewLoader(dataset,
			impl.AddStep(addPptStep),
			impl.AddStep(addNestedStep),
			impl.AddStep(getPptRecursiveStep),
			impl.AddStep(getPptStep),
			impl.AddStep(getNestedStep),
		)
		loader.Run()
		res := loader.Results()
		resb, _ := json.Marshal(res)
		fname := strings.Split(dataFile, ".")
		outFile, err := os.Create(csv.ParsePath(fmt.Sprintf("%s.json", fname[0])))
		if err != nil {
			fmt.Println(fmt.Sprintf("Error creating file %s", dataFile))
			fmt.Println(string(resb))
		}

		defer outFile.Close()
		for _, result := range []vegeta.Result(*res) {
			vegeta.NewJSONEncoder(outFile).Encode(&result)
		}
	}
}
