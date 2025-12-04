package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIDIndexForAllCollections 遍历指定数据库中的所有集合，
// 并确保每个集合的 `_id` 字段上都存在一个主键唯一索引。
// MongoDB 默认会自动创建此索引，所以此函数主要用于验证和修复可能存在的异常情况。
// EnsureIDIndexForAllCollections 遍历指定数据库中的所有集合，
// 显式检查 `_id` 主键索引是否存在，如果不存在则创建它。
func EnsureIDIndexForAllCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // 增加超时时间以防集合过多
	defer cancel()

	// 1. 获取数据库中所有集合的名称
	collectionNames, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("无法获取集合列表: %w", err)
	}

	if len(collectionNames) == 0 {
		fmt.Println("数据库中没有找到任何集合。")
		return nil
	}

	fmt.Printf("在数据库 '%s' 中找到 %d 个集合，开始检查 id 索引...\n", db.Name(), len(collectionNames))

	// 2. 遍历每个集合
	for _, collName := range collectionNames {
		collection := db.Collection(collName)
		idIndexExists := false

		// 3. 获取当前集合的所有索引
		cursor, err := collection.Indexes().List(ctx)
		if err != nil {
			logs.Printf(ctx, nil, "警告：无法获取集合 '%s' 的索引列表: %v。将尝试创建索引。\n", collName, err)
			// 即使无法列出索引，我们仍然尝试创建，因为 CreateOne 是幂等的
		} else {
			// 4. 遍历索引列表进行检查
			var results []bson.M
			if err = cursor.All(ctx, &results); err != nil {
				logs.Printf(ctx, nil, "警告：无法解码集合 '%s' 的索引: %v。将尝试创建索引。\n", collName, err)
			}

			for _, index := range results {
				if name, ok := index["name"].(string); ok && name == "id_" {
					idIndexExists = true
					break // 找到后即可退出循环
				}
			}
		}

		// 5. 根据检查结果决定是否创建索引
		if idIndexExists {
			//fmt.Printf("✅ 集合 '%s' 的 'id_' 索引已存在。\n", collName)
		} else {
			//fmt.Printf("⚠️ 集合 '%s' 未找到 'id_' 索引，正在创建...\n", collName)

			indexModel := mongo.IndexModel{
				Keys:    bson.D{{Key: "id", Value: 1}},
				Options: options.Index().SetUnique(true),
			}

			_, err := collection.Indexes().CreateOne(ctx, indexModel)
			if err != nil {
				fmt.Printf("❌ 集合 '%s' 创建 id 索引失败: %s\n", collName, err.Error())
				continue // 继续处理下一个集合
			}
			//fmt.Printf("✅ 已为集合 '%s' 成功创建索引 '%s'。\n", collName, indexName)
		}
	}

	//fmt.Println("\n所有集合的 id 索引检查和创建操作已完成。")
	return nil
}
