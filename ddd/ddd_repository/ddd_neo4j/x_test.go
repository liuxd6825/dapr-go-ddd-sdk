package ddd_neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	neo4jURL = "bolt://192.168.64.4:7687"
	username = "neo4j"
	password = "123456"
	driver   neo4j.DriverWithContext
)

func init() {
	d, err := CreateDriver(neo4jURL, username, password)
	if err != nil {
		panic(err)
	}
	driver = d
}

func CreateDriver(uri, username, password string) (neo4j.DriverWithContext, error) {
	return neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
}

func CloseDriver(ctx context.Context, driver neo4j.DriverWithContext) error {
	return driver.Close(ctx)
}
