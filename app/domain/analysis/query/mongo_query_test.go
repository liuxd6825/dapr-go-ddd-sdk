package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	{
	    "options":{
	      "sort":{
	        "timeKey":"time",
	        "timeType": "month"
	      }
	    },
	    "collection":"master_record_day_view",
	    "pipeline":[
	      {
	        "$match":{
	           "case_id":"1001",
	           "date":{"$gte": "ISODate(2023-11-15T10:00:00Z)"}
	        }
	      },
	      {
	        "$group": {
	          "_id": { "year": "$year", "month": "$month" },
	          "amount": { "$sum": "$amount" },
	          "record": { "$sum": 1 },
	          "opp": { "$addToSet": "$opp_name" }
	        }
	      },
	      {
	        "$addFields": {
	          "opp": { "$size": "$opp" },
	          "time": {
	            "$toDate": {
	              "$concat": [
	                { "$toString": "$_id.year" },
	                "-",
	                { "$toString": "$_id.month" },
	                "-01"
	              ]
	            }
	          }
	        }
	      },
	      {
	        "$project": {
	          "_id": 0,
	          "time": 1,
	          "amount": 1,
	          "record": 1,
	          "opp": 1
	        }
	      },
	      {
	        "$sort": {
	          "time": 1
	        }
	      }
	    ]
	}
*/
func Test_AggregateQueryInit(t *testing.T) {
	q := AggregateQuery{
		Collection: "test",
		Params: map[string]AggregateParam{
			"startYear": {
				Value:     "2019",
				DataType:  "date",
				Formatter: "%s-01-01T10:00:00Z",
			},
			"endYear": {
				Value:     "2032",
				DataType:  "date",
				Formatter: "%s-01-01T10:00:00Z",
			},
		},
		Pipeline: []map[string]any{
			{
				"$match": map[string]any{
					"case_id": "1001",
					"date":    map[string]any{"$gte": "@startYear", "$lt": "@endYear"},
				},
			},
		},
	}
	q.Init()
	t.Log(q)
}

func Test_AggregateQuery_NullRemove(t *testing.T) {
	q := AggregateQuery{
		Collection: "test",
		Params: map[string]AggregateParam{
			"startYear": {
				Value:      "2019",
				DataType:   "date",
				Formatter:  "%s-01-01T10:00:00Z",
				NullRemove: true,
			},
			"endYear": {
				Value:      "",
				DataType:   "date",
				Formatter:  "%s-01-01T10:00:00Z",
				NullRemove: true,
			},
		},
		Pipeline: []map[string]any{
			{
				"$match": map[string]any{
					"case_id": "1001",
					"date":    map[string]any{"$gte": "@startYear", "$lt": "@endYear"},
				},
			},
		},
	}
	q.Init()
	match := q.GetMatch()
	if dateMap, ok := match["date"].(map[string]any); ok {
		assert.Equal(t, 2, len(dateMap))
	}
}

func Test_AggregateQuery_NullRemove2(t *testing.T) {
	q := AggregateQuery{
		Collection: "test",
		Params: map[string]AggregateParam{
			"startYear": {
				Value:      "",
				DataType:   "date",
				Formatter:  "%s-01-01T00:00:00Z",
				NullRemove: true,
			},
			"endYear": {
				Value:      "",
				DataType:   "date",
				Formatter:  "%s-01-01T00:00:00Z",
				NullRemove: true,
			},
		},
		Pipeline: []map[string]any{
			{
				"$match": map[string]any{
					"case_id": "1001",
					"date":    map[string]any{"$gte": "@startYear", "$lte": "@endYear"},
				},
			},
		},
	}
	q.Init()
	match := q.GetMatch()
	_, ok := match["date"]
	assert.Equal(t, false, ok)
}
