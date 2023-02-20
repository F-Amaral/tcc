package level

import "github.com/F-Amaral/tcc/pkg/tree/domain/entity"

func GetNodeLevel(nodeID string, nodes []*entity.Node) int {
	if nodeID == "" {
		// the root node has level 0
		return 0
	}

	for _, node := range nodes {
		if node.Id == nodeID {
			return GetNodeLevel(node.ParentId, nodes) + 1
		}
	}

	// if the node was not found, return -1 to indicate an error
	return -1
}
