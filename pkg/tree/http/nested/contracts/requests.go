package contracts

type GetTreeRequest struct {
	Id string `uri:"id"`
}

type AddToParentRequest struct {
	NodeId   string `uri:"childId"`
	ParentId string `uri:"parentId"`
}

type RemoveFromParentRequest struct {
	NodeId   string `uri:"childId"`
	ParentId string `uri:"parentId"`
}
