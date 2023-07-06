package generator

import (
	"fmt"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
	"github.com/google/uuid"
	"os"
	"path"
	"path/filepath"
)

type NodeGenerator struct {
	depth    int
	numNodes int
}

func NewNodeGenerator(numNodes, depth int) *NodeGenerator {
	return &NodeGenerator{
		numNodes: numNodes,
		depth:    depth,
	}
}

func (ng *NodeGenerator) GenerateRoot() entity.Node {
	rootId := uuid.New()
	root := &entity.Node{Id: rootId.String(), Level: 0, Children: []*entity.Node{}}
	ng.numNodes--
	ng.Generate(root, 0)
	return *root
}

func (ng *NodeGenerator) Generate(parent *entity.Node, depth int) {
	if depth == ng.depth || ng.numNodes == 0 {
		return
	}

	for i := 0; i < ng.depth && ng.numNodes > 0; i++ {
		ng.numNodes--
		childId := uuid.New()
		child := &entity.Node{
			Id:       childId.String(),
			ParentId: parent.Id,
			Level:    depth + 1,
			Children: []*entity.Node{},
		}
		parent.Children = append(parent.Children, child)
		ng.Generate(child, depth+1)
	}
}

func PostOrderTraversal(root *entity.Node) []entity.Node {
	var nodes []entity.Node

	for _, child := range root.Children {
		nodes = append(nodes, PostOrderTraversal(child)...)
	}

	nodes = append(nodes, *root)

	return nodes
}

func SaveNodesToFile(node entity.Node, basePath, filename, format string) {
	nodes := PostOrderTraversal(&node)
	if !filepath.IsAbs(basePath) {
		wd, _ := os.Getwd()
		basePath = filepath.Join(wd, filepath.Clean(basePath))
	}

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	records := buildDataToFormat(nodes, format)

	fullName := path.Join(basePath, filename)
	if err := csv.WriteCSVFile(fullName, records); err != nil {
		panic(err)
	}
}

func buildDataToFormat(nodes []entity.Node, format string) [][]string {
	var records [][]string
	switch format {
	case "vegeta":
		for _, node := range nodes {
			if node.ParentId != "" {
				record := []string{fmt.Sprintf("%s/%s", node.ParentId, node.Id)}
				records = append(records, record)
			}
		}
	default:
		records = [][]string{[]string{"parentId", "id"}}
		for _, node := range nodes {
			if node.ParentId != "" {
				record := []string{node.ParentId, node.Id}
				records = append(records, record)
			}
		}
	}

	return records
}
