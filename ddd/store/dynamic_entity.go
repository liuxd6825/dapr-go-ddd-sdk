package store

type DataMap map[string]any

type IDataMap interface {
	GetDataMap() map[string]any
	SetDataMap(data map[string]any)
}
