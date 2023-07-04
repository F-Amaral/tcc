package main

import (
	"flag"
	vegeta "github.com/F-Amaral/tcc/scripts/utils/vegeta/impl"
)

var (
	csvInputPathFlag = flag.String("file", "./data.csv", "Path to csv file")
	csvFolderFlag    = flag.String("folder", "", "Path to csv file")
	outFileFlag      = flag.String("output", "input", "Path to output file")
	modeFlag         = flag.String("mode", "nested", "Algorithm to test: ppt, recursive (ppt) , nested, all")
	targetsFlag      = flag.String("targets", "http://localhost:8080", "Target URL")
	executionFlag    = flag.Bool("exec", false, "Execution mode: vegeta, use with -input and -mode all ")
)

func main() {
	flag.Parse()
	csvInputPath := *csvInputPathFlag
	csvFolder := *csvFolderFlag
	outFile := *outFileFlag
	mode := *modeFlag
	targets := *targetsFlag
	execution := *executionFlag

	vegeta.Run(csvFolder, csvInputPath, outFile, mode, targets, execution)
}
