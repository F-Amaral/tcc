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
	width    int
	numNodes int
}

func NewNodeGenerator(numNodes, depth, width int) *NodeGenerator {
	return &NodeGenerator{
		numNodes: numNodes,
		depth:    depth,
		width:    width,
	}
}

func (ng *NodeGenerator) GenerateRoot() entity.Node {
	rootId := uuid.New()
	root := &entity.Node{Id: rootId.String(), Level: 0, Children: []*entity.Node{}}
	ng.numNodes--
	ng.Generate(root, 1)
	return *root
}

func (ng *NodeGenerator) Generate(parent *entity.Node, currentDepth int) {
	if ng.numNodes <= 0 {
		return
	}

	if currentDepth < ng.depth {
		for i := 0; i < ng.width && ng.numNodes > 0; i++ {
			ng.GenerateChild(parent, currentDepth)
		}
	}
}

func (ng *NodeGenerator) GenerateChild(parent *entity.Node, currentDepth int) {
	if ng.numNodes <= 0 {
		return
	}

	childId := uuid.New()
	child := &entity.Node{
		Id:       childId.String(),
		ParentId: parent.Id,
		Level:    currentDepth,
		Children: []*entity.Node{},
	}
	parent.Children = append(parent.Children, child)
	ng.numNodes--

	ng.Generate(child, currentDepth+1)
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
