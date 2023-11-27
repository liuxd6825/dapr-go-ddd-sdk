package types

type DataType string

const (
	DataTypeString   DataType = "string"
	DataTypeInt      DataType = "int"
	DataTypeFloat    DataType = "float"
	DataTypeMoney    DataType = "money"
	DataTypeDate     DataType = "date"
	DataTypeDateTime DataType = "dateTime"
	DataTypeBool     DataType = "bool"
	DataTypeArray    DataType = "array"
	DataTypeObject   DataType = "object"
)

func (d DataType) Name() string {
	return string(d)
}

func (d DataType) IsDate() bool {
	return d.Name() == DataTypeDate.Name()
}

func (d DataType) IsDateTime() bool {
	return d.Name() == DataTypeDateTime.Name()
}

func (d DataType) IsString() bool {
	return d.Name() == DataTypeString.Name()
}

func (d DataType) IsInt() bool {
	return d.Name() == DataTypeInt.Name()
}

func (d DataType) IsFloat() bool {
	return d.Name() == DataTypeFloat.Name()
}
