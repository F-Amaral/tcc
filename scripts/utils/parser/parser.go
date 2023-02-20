package parser

import (
	"github.com/F-Amaral/tcc/pkg/adjlist/domain/entity"
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
