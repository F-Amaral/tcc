package main

import (
	"fmt"
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

	var tree *widgets.Tree
	switch nn := interface{}(nestedNodes[0]).(type) {
	case entity.NestedNode:
		tree = createNestedTree(nestedNodes)
	case entity.Node:
		tree = createNodeTree(adjNodes)
	default:
		panic(fmt.Sprintf("unknown tree type: %T", nn))
	}

	x, y := termui.TerminalDimensions()

	tree.SetRect(0, 0, x, y)
	termui.Render(tree)

	handleUiInput(tree)
}

func createNodeTree(nodes []entity.Node) *widgets.Tree {
	tree := widgets.NewTree()
	tree.Title = "Tree"
	tree.TextStyle = termui.NewStyle(termui.ColorGreen)
	tree.SetNodes(createTreeNodes(nodes, ""))

	return tree
}

func createNestedTree(nodes []entity.NestedNode) *widgets.Tree {
	tree := widgets.NewTree()
	tree.Title = "Nested Set Tree"
	tree.TextStyle = termui.NewStyle(termui.ColorGreen)
	tree.SetNodes(createNestedTreeNodes(nodes, 1, 2*len(nodes)))

	return tree
}

func createNestedTreeNodes(nodes []entity.NestedNode, left, right int) []*widgets.TreeNode {
	var treeNodes []*widgets.TreeNode

	for _, node := range nodes {
		if node.Left >= left && node.Right < right {
			children := createNestedTreeNodes(nodes, node.Left, node.Right)
			treeNode := &widgets.TreeNode{
				Value:    node,
				Expanded: false,
				Nodes:    children,
			}

			treeNodes = append(treeNodes, treeNode)
		} else if node.Left == left && node.Right == right {
			treeNode := &widgets.TreeNode{
				Value:    node,
				Expanded: false,
				Nodes:    nil,
			}

			treeNodes = append(treeNodes, treeNode)
		}
	}

	return treeNodes
}

func createTreeNodes(nodes []entity.Node, parentID string) []*widgets.TreeNode {
	var treeNodes []*widgets.TreeNode

	for _, node := range nodes {
		if node.ParentId == parentID {
			children := createTreeNodes(nodes, node.Id)
			treeNode := &widgets.TreeNode{
				Value:    node,
				Expanded: false,
				Nodes:    children,
			}

			treeNodes = append(treeNodes, treeNode)
		}
	}

	return treeNodes
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
