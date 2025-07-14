package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	host                   = "Host"
	username               = "User"
	password               = "Pwd"
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

	// mongodb://<User>:<Pwd@<Host>/<database><Params>
	connectionURIFormatWithAuthentication = "mongodb://%s:%s@%s/%s"

	// mongodb://<Host>/<database><Params>
	connectionURIFormat = "mongodb://%s/%s%s"

	// mongodb+srv://<server>/<Params>
	connectionURIFormatWithSrv = "mongodb+srv://%s/%s"
)

func getMongoDBClient(config *Config) (*mongo.Client, error) {

	uri := getMongoURI(&config.DB)
	//fmt.Println(uri)
	// Set client options
	opts := options.Client().ApplyURI(uri).SetAppName(config.Tool.Name)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.DB.OperationTimeout))
	defer cancel()

	if len(config.DB.ReplicaSet) != 0 {
		opts.SetReplicaSet(config.DB.ReplicaSet)
	}
	if config.DB.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(config.DB.MaxPoolSize)
	}
	if config.DB.ConnectTimeout > 0 {
		opts.SetConnectTimeout(config.DB.ConnectTimeout)
	}
	if config.DB.ServerSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(config.DB.ServerSelectionTimeout)
	}
	if config.DB.SocketTimeout > 0 {
		opts.SetSocketTimeout(config.DB.SocketTimeout)
	}
	if config.DB.HeartbeatInterval > 0 {
		opts.SetHeartbeatInterval(config.DB.HeartbeatInterval)
	}
	if config.DB.LocalThreshold > 0 {
		opts.SetLocalThreshold(config.DB.LocalThreshold)
	}
	if config.DB.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(config.DB.MaxConnIdleTime)
	}
	opts.Direct = config.DB.Direct

	authSource := "admin"
	if config.DB.AuthSource != "" {
		authSource = config.DB.AuthSource
	}

	authMechanism := "SCRAM-SHA-256"
	if config.DB.AuthMechanism != "" {
		authMechanism = config.DB.AuthMechanism
	}

	if opts.Auth == nil {
		opts.Auth = &options.Credential{}
	}
	opts.Auth.AuthSource = authSource
	opts.Auth.AuthMechanism = authMechanism

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
		UseLocalTimeZone: times.IsLocalTimeZone(),
	})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v uri:%v", err.Error(), uri))
	}

	rp, err := readpref.New(readpref.PrimaryMode)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s", err.Error()))
	}
	err = client.Ping(ctx, rp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s", err.Error()))
	}
	return client, nil
}

func getMongoURI(metadata *DBConfig) string {
	if metadata.UserName != "" && metadata.Password != "" {
		return fmt.Sprintf(connectionURIFormatWithAuthentication, metadata.UserName, metadata.Password, metadata.Host, metadata.DBName)
	}
	return fmt.Sprintf(connectionURIFormat, metadata.Host, metadata.DBName, metadata.Options)
}
