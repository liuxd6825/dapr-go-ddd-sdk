package xtest

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

type MongoConn struct {
	Client *mongo.Client
}

func NewMongo() *mongodb.MongoDB {
	var client *mongo.Client

	// 设置MongoDB连接URL
	clientOptions := options.Client().ApplyURI("mongodb://192.168.65.5:27018,192.168.65.5:27019,192.168.65.5:27020/?retryWrites=false&replicaSet=mongors&readPreference=primary&serverSelectionTimeoutMS=5000&connectTimeoutMS=10000")

	// 连接到MongoDB
	clientVal, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	client = clientVal
	// 确保连接成功
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	db := mongodb.NenMongoDBWithClient("test", client)
	return db
}

type MySQLConn struct {
	DB        *gorm.DB
	UserName  string
	Password  string
	Host      string
	Port      string
	Charset   string
	DbName    string
	ParseTime bool
	Loc       string
}

func NewMySQL() *gorm.DB {
	c := &MySQLConn{
		UserName:  "root",
		Password:  "11111111",
		Host:      "localhost",
		Port:      "3306",
		DbName:    "test",
		ParseTime: true,
	}

	dsn := c.dsn()
	dialector := mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  false, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true,  // 根据当前 MySQL 版本自动配置
	})
	db, err := gorm.Open(dialector)
	if err != nil {
		panic(err)
	}

	return db
}

// DSN 构造 MySQL 连接字符串
func (c *MySQLConn) dsn() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?", c.UserName, c.Password, c.Host, c.Port, c.DbName)
	if c.Charset != "" {
		dsn = fmt.Sprintf("%s&charset=%s", dsn, c.Charset)
	}
	dsn = fmt.Sprintf("%s&parseTime=%v", dsn, c.ParseTime)
	if c.Loc != "" {
		dsn = fmt.Sprintf("%s&loc=%s", dsn, c.Loc)
	}
	dsn = strings.ReplaceAll(dsn, "?&", "?")
	return dsn
}

func NewSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db
}

func NewNeo4j() neo4j.DriverWithContext {
	ctx := context.Background()
	uri := fmt.Sprintf("bolt://%v:%v", "127.0.0.1", 7687)
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth("neo4j", "12345678", ""))
	if err != nil {
		panic(err)
	}
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(fmt.Sprintf("连接neo4j失败, error:%s。%s  ", uri, err.Error()))
	}
	return driver
}
