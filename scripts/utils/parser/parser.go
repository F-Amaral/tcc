package parser

import (
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	treeUtil "github.com/F-Amaral/tcc/scripts/utils/level"
)

func ParseData(records [][]string) ([]entity.Node, error) {
	var nodesPtr []*entity.Node
	var nodes []entity.Node
	for i, record := range records {
		if i == 0 {
			// Skip header row
			continue
		}

		node := entity.Node{
			Id:       record[0],
			ParentId: record[1],
		}

		nodesPtr = append(nodesPtr, &node)
	}

	for _, node := range nodesPtr {
		level := treeUtil.GetNodeLevel(node.Id, nodesPtr)
		node.Level = level
		nodes = append(nodes, *node)
	}

	return nodes, nil
}

func ConvertToNestedSet(adjNodes []entity.Node) []entity.NestedNode {
	var nestedNodes []entity.NestedNode
	left := 1

	for _, adjNode := range adjNodes {
		nestedNode := entity.NestedNode{
			Id:    adjNode.Id,
			Level: adjNode.Level,
		}

		if adjNode.ParentId == "" {
			nestedNode.Left = left
			left++
		} else {
			parentIndex := getNodeIndex(nestedNodes, adjNode.ParentId)
			nestedNode.Left = nestedNodes[parentIndex].Right
			left = nestedNodes[parentIndex].Right + 1
			updateRightValues(nestedNodes, parentIndex, nestedNode.Left-1)
		}

		nestedNode.Right = left
		left++

		nestedNodes = append(nestedNodes, nestedNode)
	}

	return nestedNodes
}

func updateRightValues(nestedNodes []entity.NestedNode, start int, adjust int) {
	for i := start; i < len(nestedNodes); i++ {
		if nestedNodes[i].Left > adjust {
			nestedNodes[i].Left += 2
			nestedNodes[i].Right += 2
		}
	}
}

func getNodeIndex(nestedNodes []entity.NestedNode, nodeId string) int {
	for i, node := range nestedNodes {
		if node.Id == nodeId {
			return i
		}
	}

	return -1
}
