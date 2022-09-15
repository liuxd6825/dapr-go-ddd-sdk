package ddd_neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	neo4jURL = "bolt://192.168.64.4:7687"
	username = "neo4j"
	password = "123456"
	driver   neo4j.Driver
)

func init() {
	d, err := CreateDriver(neo4jURL, username, password)
	if err != nil {
		panic(err)
	}
	driver = d
}

func CreateDriver(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}

func CloseDriver(driver neo4j.Driver) error {
	return driver.Close()
}
