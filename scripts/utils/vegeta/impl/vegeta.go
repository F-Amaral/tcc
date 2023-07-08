package vegeta

import (
	"encoding/json"
	"fmt"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
	"github.com/F-Amaral/tcc/scripts/utils/vegeta/enums"
	vegeta "github.com/tsenart/vegeta/lib"
	"github.com/tsenart/vegeta/lib/plot"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Target struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

func (t Target) ToText() string {
	return fmt.Sprintf("%s %s\n", t.Method, t.URL)
}

func (t Target) ToJson() string {
	if m, err := json.Marshal(t); err == nil {
		return fmt.Sprintf("%s\n", m)
	}
	return ""
}

func Run(csvFolder, csvInputPath, outFile, mode, targets string, execution bool) {
	if csvFolder != "" {
		files, err := os.ReadDir(path.Clean(csvFolder))
		if err != nil {
			panic(err)
		}
		for i, file := range files {
			if strings.HasSuffix(file.Name(), ".csv") {
				inputData, err := csv.ReadFromCSV(csvFolder + "/" + file.Name())
				if err != nil {
					panic(err)
				}
				if !execution {
					outFile := path.Join(csvFolder, fmt.Sprintf("%s-%d", outFile, len(inputData)))
					parseToFile(inputData, mode, targets, outFile)
					os.Exit(0)
				} else {
					exec(inputData, fmt.Sprintf("report-%s-%d", outFile, i), mode, targets)
				}
			}
		}
		os.Exit(0)
	} else {
		inputData, err := csv.ReadFromCSV(csvInputPath)
		if err != nil {
			panic(err)
		}

		if !execution {
			parseToFile(inputData, mode, targets, outFile)
			os.Exit(0)
		} else {
			exec(inputData, fmt.Sprintf("report-%s", outFile), mode, targets)
		}
	}
}

func exec(inputData [][]string, reportFileName, modeStr, targetStr string) {
	mode, err := enums.NameOf(modeStr)
	if err != nil {
		log.Fatal(err)
	}

	rate := vegeta.Rate{Freq: 100, Per: time.Second}

	attacker := vegeta.NewAttacker(
		vegeta.Timeout(60*time.Second),
		vegeta.Workers(1))
	targeters := buildTargetersFromMap(parseToModeTarget(inputData, mode, targetStr))
	//targeter := buildTargeter(parseToModeTarget(inputData, mode, targetStr))
	for mode, targeter := range targeters {
		var metrics vegeta.Metrics
		plotter := plot.New(plot.Title("test"))
		for res := range attacker.Attack(targeter, rate, 10*time.Minute, mode.String()) {
			metrics.Add(res)
			fmt.Println(fmt.Sprintf("CODE: %d, ERROR: %s", res.Code, res.Error))
			plotter.Add(res)
		}
		metrics.Close()
		outFile := createFile(reportFileName, &mode)
		//plotter.WriteTo(outFile)
		reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)
		reporter.Report(outFile)
		//reporter := vegeta.NewJSONReporter(&metrics)
		//err = reporter.Report(outFile)
		//if err != nil {
		//	panic(err)
		//}
	}

}

func parseToFile(inputData [][]string, modeStr, targetStr, outFile string) {
	mode, err := enums.NameOf(modeStr)
	if err != nil {
		log.Fatal(err)
	}

	subModes := mode.Expand()
	fileMap := make(map[enums.Mode]*os.File)
	for _, subMode := range subModes {
		file := createFile(outFile, &subMode)
		fileMap[subMode] = file
	}

	targetMap := parseToModeTarget(inputData, mode, targetStr)
	for mode, targets := range targetMap {
		for _, target := range targets {
			_, err := fileMap[mode].WriteString(target.ToText())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	//file := createFile(outFile, mode)
	//
	//targets := parseToModeTarget(inputData, mode, targetStr)
	//for _, target := range targets {
	//	_, err := file.WriteString(target.ToText())
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

}

func createFile(outFileName string, mode *enums.Mode) *os.File {
	outFile, err := os.Create(buildFileName(outFileName, mode))
	if err != nil {
		panic(err)
	}
	return outFile
}

func buildFileName(fileName string, mode *enums.Mode) string {
	outFormat := "%s-%s-write.txt"
	if mode.IsRead() {
		outFormat = "%s-%s-read.txt"
	}
	return fmt.Sprintf(outFormat, fileName, mode.String())
}

func buildTarget(mode enums.Mode, method, targets, parentId, nodeId string) Target {
	addPath := fmt.Sprintf("%s/%s", parentId, nodeId)

	url := fmt.Sprintf(mode.Template(), targets, addPath)
	if method == "GET" {
		url = fmt.Sprintf(mode.Template(), targets, parentId)
	}

	return Target{
		Method: method,
		URL:    url,
	}
}

//func parseToModeTarget(inputData [][]string, mode *enums.Mode, targetStr string) []Target {
//	targets := make([]Target, 0)
//	for _, endpoint := range inputData[1:] {
//		for _, mode := range mode.Expand() {
//			if !mode.Is(enums.Recursive) {
//				target := buildTarget(mode, "POST", targetStr, endpoint[0], endpoint[1])
//				targets = append(targets, target)
//			}
//			targets = append(targets, buildTarget(mode, "GET", targetStr, endpoint[0], endpoint[1]))
//		}
//	}
//	return targets
//}

func parseToModeTarget(inputData [][]string, mode *enums.Mode, targetStr string) map[enums.Mode][]Target {
	modeTargets := make(map[enums.Mode][]Target)
	for _, endpoint := range inputData[1:] {
		for _, mode := range mode.Expand() {
			if !mode.IsRead() {
				target := buildTarget(mode, "POST", targetStr, endpoint[0], endpoint[1])
				modeTargets[mode] = append(modeTargets[mode], target)
			} else {
				getEndpoint := buildTarget(mode, "GET", targetStr, endpoint[0], endpoint[1])
				modeTargets[mode] = append(modeTargets[mode], getEndpoint)
			}
		}
	}
	return modeTargets
}

func buildTargetersFromMap(input map[enums.Mode][]Target) map[enums.Mode]vegeta.Targeter {
	targeters := make(map[enums.Mode]vegeta.Targeter)
	for mode, targets := range input {
		targetstr := ""
		for _, target := range targets {
			targetstr += target.ToJson()
		}
		reader := strings.NewReader(targetstr)
		targeters[mode] = vegeta.NewJSONTargeter(reader, []byte{}, http.Header{})
	}
	return targeters
}
func buildTargeter(input []Target) vegeta.Targeter {
	targetstr := ""
	for _, targets := range input {
		targetstr += targets.ToText()
	}
	reader := strings.NewReader(targetstr)
	return vegeta.NewHTTPTargeter(reader, []byte{}, http.Header{})
}
