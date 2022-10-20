package metadata

type DataType int64

const (
	DataTypeInt    DataType = 1
	DataTypeStr    DataType = 2
	DataTypeTime   DataType = 3
	DataTypeBool   DataType = 4
	DataTypeStruct DataType = 5
)
