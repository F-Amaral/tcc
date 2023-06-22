package contracts

import (
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Node struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ParentId  *string        `gorm:"index"`
	TreeId    string         `gorm:"index"`
	Level     int            `gorm:"-"`
	Left      int            `gorm:"column:lft"`
	Right     int            `gorm:"column:rgt"`
	Children  []*Node        `gorm:"foreignKey:ParentId;constraint:OnDelete:SET NULL"`
}

func MapFromEntity(node *entity.Node) *Node {
	var parentId *string
	parentId = nil
	if node.ParentId != "" {
		parentId = &node.ParentId
	}
	return &Node{
		ID:        node.Id,
		CreatedAt: time.Now(),
		ParentId:  parentId,
		Left:      1,
		Right:     2,
		TreeId:    uuid.New().String(),
		Level:     node.Level,
	}
}

func MapToEntity(node *Node) *entity.Node {
	parentId := ""
	if node.ParentId != nil {
		parentId = *node.ParentId
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
