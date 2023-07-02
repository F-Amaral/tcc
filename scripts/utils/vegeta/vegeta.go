package main

import (
	"flag"
	"fmt"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
	"log"
	"os"
)

const (
	pptRecursiveFormat = "POST %s/ppt/%s\n"
	pptDefaultFormat   = "POST %s/ppt/%s=?recursive=false\n"
	nestedFormat       = "POST %s/nested/%s\n"
)

var (
	csvInputPathFlag = flag.String("f", "./data.csv", "Path to csv file")
	outFileFlag      = flag.String("o", "input", "Path to output file")
	modeFlag         = flag.String("mode", "ppt", "Algorithm to test: ppt (can be used with -r), nested")
	recursiveFlag    = flag.Bool("r", true, "Recursive mode for PPT, default true")
	target           = flag.String("t", "http://localhost:8080", "Target URL")
)

func main() {
	flag.Parse()
	inputData, err := csv.ReadFromCSV(*csvInputPathFlag)
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(buildFileName(*outFileFlag, *modeFlag, *recursiveFlag))
	if err != nil {
		panic(err)
	}

	defer outFile.Close()

	for _, endpoint := range inputData {
		_, err := outFile.WriteString(fmt.Sprintf(parseMode(*modeFlag, *recursiveFlag), *target, endpoint[0]))
		if err != nil {
			log.Fatal(err)
		}
	}

}

func parseMode(input string, recursive bool) string {
	switch input {
	case "ppt", "parent", "p":
		if recursive {
			return pptRecursiveFormat
		}
		return pptDefaultFormat
	case "nested", "n", "nest":
		return nestedFormat
	default:
		return pptDefaultFormat
	}
}

func buildFileName(fileName, mode string, recursive bool) string {
	defaultFormat := "%s-%s.txt"
	recursiveFormat := "%s-%s-recursive.txt"
	if recursive {
		return fmt.Sprintf(recursiveFormat, fileName, mode)
	}
	return fmt.Sprintf(defaultFormat, fileName, mode)
}
