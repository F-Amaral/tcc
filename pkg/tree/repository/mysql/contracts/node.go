package contracts

import (
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"gorm.io/gorm"
	"time"
)

type Node struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ParentId  *string        `gorm:"index"`
	Level     int
	Children  []*Node `gorm:"foreignKey:ParentId;constraint:OnDelete:SET NULL"`
}

func MapFromEntity(node *entity.Node) *Node {
	return &Node{
		ID:        node.Id,
		CreatedAt: time.Now(),
		ParentId:  &node.ParentId,
		Level:     node.Level,
	}
}

func MapToEntity(node *Node) *entity.Node {
	parentId := *node.ParentId
	if *node.ParentId == node.ID {
		parentId = ""
	}
	return &entity.Node{
		Id:       node.ID,
		ParentId: parentId,
		Level:    node.Level,
		Children: MapToEntityList(node.Children),
	}
}

func MapToEntityList(nodes []*Node) []*entity.Node {
	var result []*entity.Node
	for _, node := range nodes {
		result = append(result, MapToEntity(node))
	}
	return result
}
