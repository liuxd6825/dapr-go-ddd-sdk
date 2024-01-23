package restapp

import "github.com/liuxd6825/dapr-go-ddd-sdk/logs"

type RunOptions struct {
	runType *RunType
	tables  *Tables
	sqlFile *string
	prefix  *string
	dbKey   *string
	level   *logs.Level
	status  *bool
	stop    *bool
}

func NewRunOptions(opts ...*RunOptions) *RunOptions {
	o := &RunOptions{
		tables: nil,
	}
	for _, item := range opts {
		if item.tables != nil {
			o.tables = item.tables
		}
		if item.runType != nil {
			o.runType = item.runType
		}
		if item.prefix != nil {
			o.prefix = item.prefix
		}
		if item.dbKey != nil {
			o.dbKey = item.dbKey
		}
		if item.sqlFile != nil {
			o.sqlFile = item.sqlFile
		}
		if item.status != nil {
			o.status = item.status
		}
		if item.stop != nil {
			o.stop = item.stop
		}
	}
	return o
}

func (o *RunOptions) SetStatus(v bool) *RunOptions {
	o.status = &v
	return o
}

func (o *RunOptions) GetStatus() bool {
	if o.status == nil {
		return false
	}
	return *o.status
}

func (o *RunOptions) SetStop(v bool) *RunOptions {
	o.stop = &v
	return o
}

func (o *RunOptions) GetStop() bool {
	if o.stop == nil {
		return false
	}
	return *o.status
}

func (o *RunOptions) GetRunType() RunType {
	if o.runType == nil {
		return *o.runType
	}
	return RunTypeStart
}

func (o *RunOptions) SetInit(v RunType) *RunOptions {
	o.runType = &v
	return o
}

func (o *RunOptions) SetSqlFile(v string) *RunOptions {
	o.sqlFile = &v
	return o
}

func (o *RunOptions) GetSqlFile() string {
	if o.sqlFile == nil {
		return ""
	}
	return *o.sqlFile
}

func (o *RunOptions) GetLevel() *logs.Level {
	return o.level
}

func (o *RunOptions) SetLevel(v *logs.Level) *RunOptions {
	o.level = v
	return o
}

func (o *RunOptions) SetPrefix(v string) *RunOptions {
	o.prefix = &v
	return o
}

func (o *RunOptions) GetPrefix() string {
	if o.prefix == nil {
		return ""
	}
	return *o.prefix
}

func (o *RunOptions) SetTable(v *Tables) *RunOptions {
	o.tables = v
	return o
}

func (o *RunOptions) GetTable() *Tables {
	return o.tables
}

func (o *RunOptions) SetDbKey(v string) *RunOptions {
	o.dbKey = &v
	return o
}

func (o *RunOptions) GetDbKey() string {
	if o.dbKey == nil {
		return ""
	}
	return *o.dbKey
}

func (o *RunOptions) SetFlag(flag *RunFlag) *RunOptions {
	o.SetPrefix(flag.Prefix)
	o.SetDbKey(flag.DbKey)
	o.SetSqlFile(flag.SqlFile)
	o.runType = &flag.RunType
	return o
}
