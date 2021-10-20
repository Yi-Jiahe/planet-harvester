package models

type Resource struct {
	Name        string
	Description string
}

type ResourceNode struct {
	Resource Resource
}

var (
	Wood = Resource{
		Name:        "Wood",
		Description: "Basic resource gathered from trees",
	}
	Iron = Resource{
		Name:        "Iron",
		Description: "Basic resource mined from the ground",
	}
	Coal = Resource{
		Name:        "Coal",
		Description: "",
	}
)

func NewTree() ResourceNode {
	tree := ResourceNode{
		Resource: Wood,
	}

	return tree
}
