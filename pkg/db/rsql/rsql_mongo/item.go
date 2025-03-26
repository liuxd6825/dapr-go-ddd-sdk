package rsql_mongo

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"

type Item struct {
	parent *Item
	name   string
	value  interface{}
	items  []*Item
}

func NewItem(parent *Item, name string) *Item {
	n := name
	if n == "id" {
		n = _id
	}
	return &Item{
		name:   n,
		parent: parent,
		value:  nil,
		items:  make([]*Item, 0),
	}
}

func (i *Item) addChildItem(name string, value interface{}) *Item {
	newItem := NewItem(i, name)
	newItem.value = value
	i.items = append(i.items, newItem)
	return newItem
}

func (i *Item) getAndItem() {
}

func (i *Item) getValues(data map[string]interface{}) {
	if len(i.items) != 0 {
		array := make([]interface{}, len(i.items))
		for i, v := range i.items {
			item := ddd_utils.NewMap()
			item[v.name] = v.value
			if len(v.items) > 0 {
				m := ddd_utils.NewMap()
				v.getValues(m)
				item[v.name] = m[v.name]
			}
			array[i] = item
		}
		data[i.name] = array
	} else if i.value != nil {
		data[i.name] = i.value
	}
}

func (i *Item) setValue(name string, value interface{}) {
	i.name = name
	i.value = value
}
