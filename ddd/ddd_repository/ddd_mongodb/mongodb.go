package ddd_mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strconv"
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
	operationTimeout       = "OperationTimeout"
	params                 = "Params"
	id                     = "_id"
	value                  = "value"
	etag                   = "_etag"

	defaultTimeout = 5 * time.Second

	// mongodb://<UserName>:<Password@<Host>/<database><Params>
	connectionURIFormatWithAuthentication = "mongodb://%s:%s@%s/%s%s"

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
	Host             string
	UserName         string
	Password         string
	DatabaseName     string
	server           string
	WriteConcern     string
	ReadConcern      string
	Params           string
	OperationTimeout time.Duration
}

// NewMongoDB returns a new MongoDB state store.
func NewMongoDB() *MongoDB {
	s := &MongoDB{}
	return s
}

// Init establishes connection to the store based on the config.
func (m *MongoDB) Init(config *Config) error {
	m.operationTimeout = config.OperationTimeout

	client, err := getMongoDBClient(config)
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

func (m *MongoDB) NewCollection(collectionName string) *mongo.Collection {
	return m.database.Collection(collectionName, m.collectionOptions)
}

func (m *MongoDB) Ping() error {
	if err := m.client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("mongoDB store: error connecting to mongoDB at %s: %s", m.config.Host, err)
	}

	return nil
}

func getMongoURI(metadata *Config) string {
	if len(metadata.server) != 0 {
		return fmt.Sprintf(connectionURIFormatWithSrv, metadata.server, metadata.Params)
	}

	if metadata.UserName != "" && metadata.Password != "" {
		return fmt.Sprintf(connectionURIFormatWithAuthentication, metadata.UserName, metadata.Password, metadata.Host, metadata.DatabaseName, metadata.Params)
	}

	return fmt.Sprintf(connectionURIFormat, metadata.Host, metadata.DatabaseName, metadata.Params)
}

func getMongoDBClient(metadata *Config) (*mongo.Client, error) {
	uri := getMongoURI(metadata)

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), metadata.OperationTimeout)
	defer cancel()

	mongoLog := "mongodb-"
	if clientOptions.AppName != nil {
		clientOptions.SetAppName(mongoLog + ":" + *clientOptions.AppName)
	} else {
		clientOptions.SetAppName(mongoLog)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
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
