package parser

import (
	"fmt"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	treeUtil "github.com/F-Amaral/tcc/scripts/utils/level"
	"sort"
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
		level := treeUtil.GetNodeLevel(node.ParentId, nodesPtr)
		node.Level = level
		nodes = append(nodes, *node)
		print(fmt.Sprintf("%v", node))
	}

	return nodes, nil

}

func ConvertToNestedSet(nodes []entity.Node) []entity.NestedNode {
	var rootNode *entity.Node
	for i := range nodes {
		if nodes[i].ParentId == "" {
			rootNode = &nodes[i]
			break
		}
	}

	treeMap := buildTreeMap(rootNode.Id, nodes)
	_, nestedSet := buildNestedSet(rootNode, 1, treeMap)
	return sortNestedSet(nestedSet)
}

func buildTreeMap(rootId string, nodes []entity.Node) map[string][]*entity.Node {
	treeMap := make(map[string][]*entity.Node)
	rootChildren := removeParentFromNodeSlice(rootId, nodes)

	for _, node := range rootChildren {
		treeMap[node.ParentId] = append(treeMap[node.ParentId], node)
	}
	return treeMap
}

func removeParentFromNodeSlice(parentId string, nodes []entity.Node) []*entity.Node {
	var childNodes []*entity.Node
	for _, node := range nodes {
		if node.Id != parentId {
			n2 := node
			childNodes = append(childNodes, &n2)
		}
	}

	return childNodes
}

func createNestedNode(node entity.Node, left, right int) entity.NestedNode {
	return entity.NestedNode{
		Id:       node.Id,
		ParentId: node.ParentId,
		Left:     left,
		Right:    right,
		Level:    node.Level,
	}
}

func buildNestedSet(node *entity.Node, left int, treeMap map[string][]*entity.Node) (int, []*entity.NestedNode) {
	var nestedNodes []*entity.NestedNode

	leftValue := left
	rightValue := left + 1

	// Process each child node recursively
	for _, childNode := range treeMap[node.Id] {
		right, nestedChildNodes := buildNestedSet(childNode, rightValue, treeMap)
		rightValue = right
		nestedNodes = append(nestedNodes, nestedChildNodes...)
	}

	// Create a NestedNode for this node and append it to the result slice
	nestedNode := createNestedNode(*node, leftValue, rightValue)
	nestedNodes = append(nestedNodes, &nestedNode)

	// Update the right value for this node
	rightValue++

	return rightValue, nestedNodes
}

func sortNestedSet(ptrs []*entity.NestedNode) []entity.NestedNode {
	nestedNodes := make([]entity.NestedNode, 0)
	for _, ptr := range ptrs {
		nestedNodes = append(nestedNodes, *ptr)
	}

	sort.Slice(nestedNodes, func(i, j int) bool {
		return nestedNodes[i].Left < nestedNodes[j].Left
	})
	return nestedNodes
}
