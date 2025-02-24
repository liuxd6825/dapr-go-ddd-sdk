package rsql_mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Filter struct {
	Lookups []*Lookup
	Match   map[string]any
}

func NewMongoFilter() *Filter {
	return &Filter{
		Match:   make(map[string]any),
		Lookups: make([]*Lookup, 0),
	}
}

func (l *Filter) IsAggregate() bool {
	if len(l.Lookups) == 0 {
		return false
	}
	return true
}

func (l *Filter) NewPipeline() mongo.Pipeline {
	pipeline := mongo.Pipeline{}
	pipeline = l.AddLookups(pipeline)
	pipeline = l.AddMatch(pipeline)
	return pipeline
}

func (l *Filter) AddPipelines(pipeline mongo.Pipeline) mongo.Pipeline {
	pipeline = l.AddLookups(pipeline)
	pipeline = l.AddMatch(pipeline)
	return pipeline
}

func (l *Filter) AddLookups(pipeline mongo.Pipeline) mongo.Pipeline {
	for _, lookup := range l.Lookups {
		pipeline = append(pipeline, lookup.GetLookup())
		pipeline = append(pipeline, lookup.GetUnwind())
	}
	return pipeline
}

func (l *Filter) AddMatch(pipeline mongo.Pipeline) mongo.Pipeline {
	if l.Match != nil {
		pipeline = append(pipeline, bson.D{{"$match", l.Match}})
	}
	return pipeline
}
