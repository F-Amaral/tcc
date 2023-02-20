package main

import (
	"github.com/F-Amaral/tcc/pkg/adjlist/domain/entity"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
	"github.com/google/uuid"
	"math/rand"
)

func main() {
	numNodes := 20
	avgDepth := 3
	maxDepth := 6
	prob := 0.6

	var nodes []entity.Node

	// create root node
	rootId := uuid.New()
	nodes = append(nodes, entity.Node{Id: rootId.String()})

	// create child nodes
	for i := 0; i < numNodes-1; i++ {
		// select random parent node
		parent := nodes[randInt(0, len(nodes)-1)]
		parentDepth := getNodeDepth(parent.Id, nodes)

		// calculate depth for new node
		depth := randInt(1, maxDepth-parentDepth)
		if depth > avgDepth {
			depth = randInt(1, avgDepth)
		}

		// create new node with random depth
		childId := uuid.New()
		child := entity.Node{
			Id:       childId.String(),
			ParentId: parent.Id,
		}
		nodes = append(nodes, child)

		// recursively create children nodes
		if depth > 1 && randFloat64() < prob {
			createChildren(childId, depth-1, avgDepth, maxDepth, prob, nodes)
		}
	}

	// write nodes to csv file
	var records [][]string
	records = append(records, []string{"id", "parent_id"})
	for _, node := range nodes {
		record := []string{node.Id, node.ParentId}
		records = append(records, record)
	}

	filename := "tree.csv"
	if err := csv.WriteCSVFile(filename, records); err != nil {
		panic(err)
	}
}

func createChildren(parentId uuid.UUID, depth, avgDepth, maxDepth int, prob float64, nodes []entity.Node) {
	for i := 0; i < randInt(1, 3); i++ {
		childId := uuid.New()
		child := entity.Node{
			Id:       childId.String(),
			ParentId: parentId.String(),
		}
		nodes = append(nodes, child)

		// recursively create children nodes
		if depth > 1 && randFloat64() < prob {
			createChildren(childId, depth-1, avgDepth, maxDepth, prob, nodes)
		}
	}
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

func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randFloat64() float64 {
	return rand.Float64()
}
