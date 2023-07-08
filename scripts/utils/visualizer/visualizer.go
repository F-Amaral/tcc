package main

import (
	"flag"
	"fmt"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
)

type Flags struct {
	fileName string
}

const path = "/Users/famaral/go/src/github.com/f-amaral/tcc/scripts/utils/loader/output/dataset_20_csv.csv"

func parseFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.fileName, "file", path, "CSV file to parse")
	flag.Parse()
	return flags
}

type Node struct {
	Id       string
	Children []*Node
}

func (n *Node) print(indent string) {
	fmt.Println(n.Id)
	for _, child := range n.Children {
		fmt.Printf("%s└─ ", indent)
		child.print(indent + "   ")
	}
}

func main() {
	flags := parseFlags()

	data, err := csv.ReadFromCSV(flags.fileName)
	if err != nil {
		panic(err)
	}

	nodes := make(map[string]*Node)

	// Initialize all nodes.
	for _, line := range data[1:] {
		parentId, childId := line[0], line[1]
		if _, exists := nodes[parentId]; !exists {
			nodes[parentId] = &Node{Id: parentId}
		}
		if _, exists := nodes[childId]; !exists {
			nodes[childId] = &Node{Id: childId}
		}
	}

	// Set parent-child relationships.
	for _, line := range data[1:] {
		parentId, childId := line[0], line[1]
		nodes[parentId].Children = append(nodes[parentId].Children, nodes[childId])
	}

	// Assume the root is the parent in the last element in data.
	root := nodes[data[len(data)-1][0]]

	// Print the tree.
	root.print("")
}
