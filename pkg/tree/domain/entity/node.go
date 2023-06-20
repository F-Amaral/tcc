package entity

import (
	"encoding/json"
	"fmt"
)

type Node struct {
	Id       string
	ParentId string
	Level    int
	Children []*Node
}

type NestedNode struct {
	Id       string
	ParentId string
	Left     int
	Right    int
	Level    int
}

func (n Node) String() string {
	return fmt.Sprintf("%s (Level %d)", n.Id, n.Level)
}

func (n NestedNode) String() string {
	bytes, _ := json.Marshal(n)
	return string(bytes)

}
