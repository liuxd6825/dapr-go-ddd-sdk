package ddd_neo4j

func GetNodes[T any](list []T, lables ...string) []Node {
	var nodes []Node
	count := len(lables)
	for _, i := range list {
		var a any = i
		if node, ok := a.(Node); ok {
			if count>0{
				node.SetLabels(lables)
			}
			nodes = append(nodes, node)
		}

	}
	return nodes
}

func GetRelations[T any](list []T) []Relation {
	var nodes []Relation
	for _, i := range list {
		var a any = i
		if rel, ok := a.(Relation); ok {
			nodes = append(nodes, rel)
		}
	}
	return nodes
}
