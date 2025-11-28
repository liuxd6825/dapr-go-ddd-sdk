package service

import (
	"context"
	"fmt"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoService struct {
	env      *env.Env
	rootPath string
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoService(env *env.Env, rootPath string) *MongoService {
	dbKey := env.GetDBKeyValue(config.DBKey)
	dbItem := env.GetDB(dbKey)
	if dbItem == nil {
		panic("mongoDbItem is nil")
	}
	client := dbItem.GetMongo().Client()
	database := dbItem.GetMongo().GetDatabase()
	return &MongoService{
		env:      env,
		rootPath: rootPath,
		client:   client,
		database: database,
	}
}

func (s *MongoService) Aggregate(ctx context.Context, qry *query.AggregateQuery) (any, error) {
	qry.Init()

	// 1. 从请求体中解码聚合管道
	var pipeline = qry.Pipeline
	// 2. 检查管道是否为空
	if len(pipeline) == 0 {
		return nil, errors.New("聚合管道不能为空")
	}

	// 3. 获取集合句柄
	// 为简单起见，我们硬编码了数据库和集合名称
	// 在实际应用中，可以考虑从URL路径中动态获取
	collection := s.database.Collection(qry.Collection)

	// 4. 执行聚合查询
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("数据库聚合查询失败:%s ", err.Error())
	}
	defer cursor.Close(ctx)

	// 5. 将查询结果解码到通用结构中
	// bson.M 是 map[string]interface{} 的别名，可以表示任何BSON文档
	var results []map[string]any
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.New("解析查询结果失败: %s ", err.Error())
	}

	// 如果没有结果，返回一个空数组，这比返回 null 更友好
	if results == nil {
		results = make([]map[string]any, 0)
	}
	if qry.Options.Sort != nil {
		querySort := qry.Options.Sort
		results, err = CompleteAndSortData(results, querySort.TimeKey, querySort.TimeType)
	}
	// 6. 设置响应头并返回JSON结果
	return results, err
}

// CompleteAndSortData 对数据按指定字段排序并补全时间序列。
// data: 原始数据
// timeKey: 作为时间依据的 map key
// timeType: 补全的时间单位 (year, month, day)
func CompleteAndSortData(data []map[string]any, timeKey string, timeType model.SortTimeType) ([]map[string]any, error) {
	if len(data) == 0 {
		return data, nil
	}

	const timeLayout = "2006-01-02 15:04:05"

	// 1. 排序
	/*
		sort.Slice(data, func(i, j int) bool {
			tI, okI := data[i][timeKey].(primitive.DateTime)
			tJ, okJ := data[j][timeKey].(primitive.DateTime)
			if !okI || !okJ {
				// 在实际应用中可能需要更复杂的错误处理
				return false
			}
			ti := tI.Time()
			tj := tJ.Time()
			return ti.Before(tj)
		})
	*/
	for _, item := range data {
		item[timeKey] = item[timeKey].(primitive.DateTime).Time()
	}

	// 2. 补全
	var completedData []map[string]any
	if len(data) < 2 {
		return data, nil
	}

	for i := 0; i < len(data)-1; i++ {
		// 添加当前元素
		completedData = append(completedData, data[i])

		// 解析当前和下一个时间
		currentTime := data[i][timeKey].(time.Time)
		nextTime := data[i+1][timeKey].(time.Time)

		// 检查并填充缺失的时间点
		tempTime := currentTime
		for {
			switch timeType {
			case model.Year:
				tempTime = tempTime.AddDate(1, 0, 0)
			case model.Month:
				tempTime = tempTime.AddDate(0, 1, 0)
			case model.Day:
				tempTime = tempTime.AddDate(0, 0, 1)
			default:
				return nil, fmt.Errorf("无效的时间类型: %s", timeType)
			}

			// 如果补全的时间点在下一个真实数据点之前，则添加
			if tempTime.Before(nextTime) {
				missingEntry := map[string]any{
					timeKey: tempTime,
				}
				completedData = append(completedData, missingEntry)
			} else {
				break
			}
		}
	}

	// 添加最后一个元素
	completedData = append(completedData, data[len(data)-1])

	return completedData, nil
}
