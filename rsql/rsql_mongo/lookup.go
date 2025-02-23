package rsql_mongo

import "go.mongodb.org/mongo-driver/bson"

type Lookup struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
}

func NewLookup() *Lookup {
	return &Lookup{}
}

func (l *Lookup) GetLookup() bson.D {
	return bson.D{
		{"$lookup",
			bson.D{
				{"foreignField", l.ForeignField},
				{"localField", l.LocalField},
				{"from", l.From},
				{"as", l.As},
			},
		},
	}
}

func (l *Lookup) GetUnwind() bson.D {
	return bson.D{{"$unwind", "$" + l.As}}
}
