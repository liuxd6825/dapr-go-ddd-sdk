package ddd_mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strconv"
	"strings"
	"time"
)

const (
	host                   = "Host"
	username               = "UserName"
	password               = "Password"
	databaseName           = "DatabaseName"
	eventCollectionName    = "eventCollectionName"
	snapshotCollectionName = "snapshotCollectionName"
	server                 = "server"
	writeConcern           = "writeConcern"
	readConcern            = "readConcern"
	operationTimeout       = "operationTimeout"
	params                 = "Params"
	id                     = "_id"
	value                  = "value"
	etag                   = "_etag"

	defaultTimeout = 30 * time.Second

	// mongodb://<UserName>:<Password@<Host>/<database><Params>
	connectionURIFormatWithAuthentication = "mongodb://%s:%s@%s/%s?%s"

	// mongodb://<Host>/<database><Params>
	connectionURIFormat = "mongodb://%s/%s%s"

	// mongodb+srv://<server>/<Params>
	connectionURIFormatWithSrv = "mongodb+srv://%s/%s"
)

// MongoDB is a state store implementation for MongoDB.
type MongoDB struct {
	config            Config
	client            *mongo.Client
	operationTimeout  time.Duration
	database          *mongo.Database
	collectionOptions *options.CollectionOptions
}

type Config struct {
	AppName      string
	Host         string
	UserName     string
	Password     string
	DatabaseName string
	server       string
	WriteConcern string
	ReadConcern  string
	Options      string
	ReplicaSet   string

	MaxPoolSize            uint64
	MinPoolSize            uint64
	MaxConnecting          uint64
	OperationTimeout       time.Duration
	HeartbeatInterval      time.Duration
	LocalThreshold         time.Duration
	MaxConnIdleTime        time.Duration //设置服务器端口等待响应的最长时间，单位是毫秒，指期望建立连接的最大空闲时间.
	ServerSelectionTimeout time.Duration //控制客户端在选择服务器之前可以逗留在服务器上的时间，单位是毫秒.
	ConnectTimeout         time.Duration //控制客户端向MongoDB实例发出连接请求的最长时间.
	SocketTimeout          time.Duration //控制客户端在接收响应之前可以逗留在服务器上的时间，单位是毫秒.
}

type ObjectId string

type InitOptionsFunc = func(options *options.ClientOptions) error

func (i ObjectId) String() string {
	return string(i)
}

// ServerCount
// @Description: 获取服务器的数量
// @receiver c
// @return int
func (c *Config) ServerCount() int {
	list := strings.Split(c.Host, ",")
	return len(list)
}

// NewMongoDB
// @Description:  新建MongoDB
// @param config  配置类
// @return *MongoDB MongoDB对象
// @return error 错误信息
func NewMongoDB(config *Config, optionsFunc InitOptionsFunc) (*MongoDB, error) {
	mongodb := &MongoDB{}
	if err := mongodb.init(config, optionsFunc); err != nil {
		return nil, err
	}
	return mongodb, nil
}

// init
// @Description: 初始化
// @receiver m  *MongoDB
// @param config  配置类
// @return error 错误信息
func (m *MongoDB) init(config *Config, optionsFunc InitOptionsFunc) error {

	if config == nil {
		return errors.New("NewMongoDB() error,config is nil")
	}

	if config.OperationTimeout == 0 {
		config.OperationTimeout = defaultTimeout
	}

	client, err := getMongoDBClient(config, optionsFunc)
	if err != nil {
		return fmt.Errorf("error in creating mongodb client: %s", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("error in connecting to mongodb, Host: %s error: %s", config.Host, err)
	}

	m.client = client
	// get the write concern
	wc, err := getWriteConcernObject(config.WriteConcern)
	if err != nil {
		return fmt.Errorf("error in getting write concern object: %s", err)
	}

	// get the read concern
	rc, err := getReadConcernObject(config.ReadConcern)
	if err != nil {
		return fmt.Errorf("error in getting read concern object: %s", err)
	}

	m.config = *config
	m.collectionOptions = options.Collection().SetWriteConcern(wc).SetReadConcern(rc)
	m.database = m.client.Database(config.DatabaseName)
	return nil
}

func (m *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	return m.database.Collection(collectionName, m.collectionOptions)
}

func (m *MongoDB) Name() string {
	return m.database.Name()
}

func (m *MongoDB) Drop(ctx context.Context) error {
	return m.database.Drop(ctx)
}

func (m *MongoDB) GetCollectionNames(ctx context.Context) ([]string, error) {
	return m.database.ListCollectionNames(ctx, bson.D{})
}

func (m *MongoDB) ExistCollection(ctx context.Context, name any) (bool, error) {
	names, err := m.database.ListCollectionNames(ctx, bson.D{{"name", name}})
	if err != nil {
		return false, err
	}
	if len(names) == 1 {
		return true, err
	}
	return false, err
}

func (m *MongoDB) CreateCollection(collectionName string) error {
	ops := &options.CreateCollectionOptions{}
	return m.database.CreateCollection(context.Background(), collectionName, ops)
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoDB) Ping() error {
	if err := m.client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("mongoDB store: error connecting to mongoDB at %s: %s", m.config.Host, err)
	}

	return nil
}

func getMongoURI(metadata *Config) string {
	if len(metadata.server) != 0 {
		return fmt.Sprintf(connectionURIFormatWithSrv, metadata.server, metadata.Options)
	}
	if metadata.UserName != "" && metadata.Password != "" {
		return fmt.Sprintf(connectionURIFormatWithAuthentication, metadata.UserName, metadata.Password, metadata.Host, metadata.DatabaseName, metadata.Options)
	}
	return fmt.Sprintf(connectionURIFormat, metadata.Host, metadata.DatabaseName, metadata.Options)
}

func getMongoDBClient(config *Config, optionsFunc InitOptionsFunc) (*mongo.Client, error) {

	uri := getMongoURI(config)
	//fmt.Println(uri)
	// Set client options
	opts := options.Client().ApplyURI(uri).SetAppName(config.AppName)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), config.OperationTimeout)
	defer cancel()

	if len(config.ReplicaSet) != 0 {
		opts.SetReplicaSet(config.ReplicaSet)
	}
	if config.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(config.MaxPoolSize)
	}
	if config.ConnectTimeout > 0 {
		opts.SetConnectTimeout(config.ConnectTimeout)
	}
	if config.ServerSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(config.ServerSelectionTimeout)
	}
	if config.SocketTimeout > 0 {
		opts.SetSocketTimeout(config.SocketTimeout)
	}
	if config.HeartbeatInterval > 0 {
		opts.SetHeartbeatInterval(config.HeartbeatInterval)
	}
	if config.LocalThreshold > 0 {
		opts.SetLocalThreshold(config.LocalThreshold)
	}
	if config.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(config.MaxConnIdleTime)
	}

	/*
		// 解决mongo不是本地时区的问题
		builder := bsoncodec.NewRegistryBuilder()

		// 注册默认的编码和解码器
		bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(builder)
		bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(builder)

		// 注册时间解码器
		tTime := reflect.TypeOf(time.Time{})
		tCodec := bsoncodec.NewTimeCodec(bsonoptions.TimeCodec().SetUseLocalTimeZone(true))
	*/

	opts.SetBSONOptions(&options.BSONOptions{
		UseLocalTimeZone: false,
	})

	if optionsFunc != nil {
		_ = optionsFunc(opts)
	}
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v uri:%v", err.Error(), uri))
	}
	return client, nil
}

func getWriteConcernObject(cn string) (*writeconcern.WriteConcern, error) {
	var wc *writeconcern.WriteConcern
	if cn != "" {
		if cn == "majority" {
			wc = writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))
		} else {
			w, err := strconv.Atoi(cn)
			wc = writeconcern.New(writeconcern.W(w), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))

			return wc, err
		}
	} else {
		wc = writeconcern.New(writeconcern.W(1), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))
	}

	return wc, nil
}

func getReadConcernObject(cn string) (*readconcern.ReadConcern, error) {
	switch cn {
	case "local":
		return readconcern.Local(), nil
	case "majority":
		return readconcern.Majority(), nil
	case "available":
		return readconcern.Available(), nil
	case "linearizable":
		return readconcern.Linearizable(), nil
	case "snapshot":
		return readconcern.Snapshot(), nil
	case "":
		return readconcern.Local(), nil
	}

	return nil, fmt.Errorf("readConcern %s not found", cn)
}

// AsFieldName
// @Description: 转换为mongodb规范的字段名称
// @param name
// @return string
func AsFieldName(name string) string {
	return stringutils.SnakeString(name)
}
