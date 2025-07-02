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
	DataType_Bool   = gormschema.Bool
	DataType_Int    = gormschema.Int
	DataType_Uint   = gormschema.Uint
	DataType_Float  = gormschema.Float
	DataType_String = gormschema.String
	DataType_Date   = gormschema.Date
	DataType_Time   = gormschema.Time
	DataType_Bytes  = gormschema.Bytes
	DataType_Json   = gormschema.Json
	DataType_Array  = gormschema.Array
)
