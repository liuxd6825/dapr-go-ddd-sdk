package class

type Package struct {
	namespaceName string
	fileLines     []string
}

func NewPackage(namespaceName string, fileLines []string) *Package {
	return &Package{
		namespaceName: namespaceName,
		fileLines:     fileLines,
	}
}
