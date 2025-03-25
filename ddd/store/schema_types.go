package store

import gormschema "gorm.io/gorm/schema"

type (
	// DataType GORM data type
	DataType = gormschema.DataType
	// TimeType GORM time type
	TimeType int64
)

// GORM time types
const (
	UnixTime        TimeType = 1
	UnixSecond      TimeType = 2
	UnixMillisecond TimeType = 3
	UnixNanosecond  TimeType = 4
)

// GORM fields types
const (
	Bool   DataType = "bool"
	Int    DataType = "int"
	Uint   DataType = "uint"
	Float  DataType = "float"
	String DataType = "string"
	Date   DataType = "date"
	Time   DataType = "time"
	Bytes  DataType = "bytes"
	Object DataType = "object"
	Array  DataType = "array"
)
