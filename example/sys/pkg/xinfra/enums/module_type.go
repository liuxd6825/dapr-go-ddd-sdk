package enums

type ModuleType int

const (
	ModuleType_Enterprise ModuleType = iota
	ModuleType_Project
)

func (mt ModuleType) String() string {
	switch mt {
	case ModuleType_Enterprise:
		return "企业"
	case ModuleType_Project:
		return "项目"
	}
	return "NA"
}
