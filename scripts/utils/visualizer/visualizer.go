package main

import (
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/scripts/utils/csv"
	"github.com/F-Amaral/tcc/scripts/utils/parser"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	csvData, err := csv.ReadFromCSV("tree.csv")
	if err != nil {
		panic(err)
	}

	adjNodes, err := parser.ParseData(csvData)
	if err != nil {
		panic(err)
	}

	nestedNodes := parser.ConvertToNestedSet(adjNodes)

	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	tree := createNestedTree(nestedNodes)

	x, y := termui.TerminalDimensions()

	tree.SetRect(0, 0, x, y)
	termui.Render(tree)

	handleUiInput(tree)
}

func createNestedTree(nodes []entity.NestedNode) *widgets.Tree {
	tree := widgets.NewTree()
	tree.Title = "Tree"
	tree.TextStyle = termui.NewStyle(termui.ColorGreen)
	tree.SetNodes(createNestedTreeNodes(nodes))

	return tree
}

func createNestedTreeNodes(nodes []entity.NestedNode) []*widgets.TreeNode {
	var rootNodes []*widgets.TreeNode
	for _, node := range nodes {
		if node.ParentId == "" {
			rootNode := &widgets.TreeNode{
				Value:    node,
				Expanded: true,
				Nodes:    getNestedNodeChildren(nodes, node.Id),
			}
			rootNodes = append(rootNodes, rootNode)
		}
	}
	return rootNodes
}

func getNestedNodeChildren(nodes []entity.NestedNode, parentId string) []*widgets.TreeNode {
	var children []*widgets.TreeNode
	for _, node := range nodes {
		if node.ParentId == parentId {
			treeNode := &widgets.TreeNode{
				Value:    node,
				Expanded: true,
			}
			treeNode.Nodes = getNestedNodeChildren(nodes, node.Id)
			children = append(children, treeNode)
		}
	}
	return children
}

func handleUiInput(tree *widgets.Tree) {
	previousKey := ""
	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			tree.ScrollDown()
		case "k", "<Up>":
			tree.ScrollUp()
		case "<C-d>":
			tree.ScrollHalfPageDown()
		case "<C-u>":
			tree.ScrollHalfPageUp()
		case "<C-f>":
			tree.ScrollPageDown()
		case "<C-b>":
			tree.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				tree.ScrollTop()
			}
		case "<Home>":
			tree.ScrollTop()
		case "<Enter>":
			tree.ToggleExpand()
		case "G", "<End>":
			tree.ScrollBottom()
		case "E":
			tree.ExpandAll()
		case "C":
			tree.CollapseAll()
		case "<Resize>":
			x, y := termui.TerminalDimensions()
			tree.SetRect(0, 0, x, y)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		termui.Render(tree)
	}
}
