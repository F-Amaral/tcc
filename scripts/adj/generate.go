package main

import (
	"flag"
	"fmt"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"path"
	"sort"
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
	var nodes []entity.Node

	// create root node
	rootId := uuid.New()
	nodes = append(nodes, entity.Node{Id: rootId.String(), Level: 0})

	// create child nodes
	for i := 0; i < numNodes-1; i++ {
		// select random parent node
		parent := nodes[randInt(0, len(nodes)-1)]
		parentDepth := getNodeDepth(parent.Id, nodes)

		// calculate depth for new node
		var depth int
		if maxDepth > parentDepth {
			depth = randInt(1, maxDepth-parentDepth)
		} else {
			depth = 1 // fallback case if maxDepth is not more than parentDepth
		}
		if depth > avgDepth {
			depth = randInt(1, avgDepth)
		}

		// create new node with random depth
		childId := uuid.New()
		child := entity.Node{
			Id:       childId.String(),
			ParentId: parent.Id,
			Level:    parent.Level + 1,
		}
		nodes = append(nodes, child)

		// recursively create children nodes
		if depth > 1 && randFloat64() < prob {
			nodes = createChildren(childId, depth-1, avgDepth, maxDepth, prob, nodes)
		}
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Level > nodes[j].Level
	})

	// write nodes to csv file
	var records [][]string
	for _, node := range nodes[:len(nodes)-1] {
		record := []string{fmt.Sprintf("%s/%s", node.ParentId, node.Id)}
		records = append(records, record)
	}

	fullName := path.Join(basePath, filename)
	fmt.Println(fullName)
	if err := csv.WriteCSVFile(fullName, records); err != nil {
		panic(err)
	}
}

func createChildren(parentId uuid.UUID, depth, avgDepth, maxDepth int, prob float64, nodes []entity.Node) []entity.Node {
	for i := 0; i < randInt(1, 3); i++ {
		parent := getParent(parentId.String(), nodes)
		childId := uuid.New()
		child := entity.Node{
			Id:       childId.String(),
			ParentId: parentId.String(),
			Level:    parent.Level + 1,
		}
		nodes = append(nodes, child)

		// recursively create children nodes
		if depth > 1 && randFloat64() < prob {
			nodes = createChildren(childId, depth-1, avgDepth, maxDepth, prob, nodes)
		}
	}
	return nodes
}

func getNodeDepth(nodeId string, nodes []entity.Node) int {
	var depth int
	for _, node := range nodes {
		if node.Id == nodeId {
			depth++
			depth += getNodeDepth(node.ParentId, nodes)
			break
		}
	}
	return depth
}

func getParent(nodeId string, nodes []entity.Node) entity.Node {
	for _, node := range nodes {
		if node.Id == nodeId {
			return node
		}
	}
	return entity.Node{}
}

func randInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.Intn(max-min+1)
}

func randFloat64() float64 {
	return rand.Float64()
}
