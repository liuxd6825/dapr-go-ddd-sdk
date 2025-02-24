package mongoutils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func Test_GetMQL(t *testing.T) {

	t.Run("GetMQL", func(t *testing.T) {
		project := bson.D{{
			"$project", bson.M{
				"_id": "$_id",
				"data": bson.M{
					"$slice": bson.A{
						"$data",
						20,
						100,
					},
				},
				"total_rows": "$total_rows",
			},
		}}
		var pipeline mongo.Pipeline
		pipeline = append(pipeline, project)
		PrintPipeline(pipeline)
	})
}
