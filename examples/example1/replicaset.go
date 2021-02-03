package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/romnn/mongotypes/replicaset"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDBConfig struct {
	Host       string
	Port       uint
	User       string
	Database   string
	Password   string
	ReplicaSet string
}

func run(config *mongoDBConfig) error {
	databaseName := config.Database
	var databaseAuth string
	if config.User != "" && config.Password != "" {
		databaseAuth = fmt.Sprintf("%s:%s@", config.User, config.Password)
	}
	databaseHost := fmt.Sprintf("%s:%d", config.Host, config.Port)
	databaseConnectionURI := fmt.Sprintf("mongodb://%s%s/?connect=direct", databaseAuth, databaseHost)
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseConnectionURI))
	if err != nil {
		return fmt.Errorf("Failed to create database client: %v (%s:%s)", err, databaseConnectionURI, databaseName)
	}
	mctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(mctx)

	adminDatabase := client.Database("admin")
	replConfig := bson.D{
		{Key: "_id", Value: config.ReplicaSet},
		{Key: "version", Value: 1},
		{Key: "protocolVersion", Value: 1},
		{Key: "members", Value: bson.A{
			bson.D{
				{Key: "_id", Value: 0},
				{Key: "host", Value: databaseHost},
			},
		}},
	}

	initiateCommand := bson.D{{Key: "replSetInitiate", Value: replConfig}}
	reconfigCommand := bson.D{{Key: "replSetReconfig", Value: replConfig}}
	statusCommand := bson.D{{Key: "replSetGetStatus", Value: nil}}

	var initiateResult bson.M
	var reconfigResult replicaset.OkResponse
	var statusResult replicaset.Status

	if err := adminDatabase.RunCommand(context.TODO(), initiateCommand).Decode(&initiateResult); err != nil {
		if !strings.Contains(err.Error(), "already initialized") {
			log.Println(err)
		}
	}
	if err := adminDatabase.RunCommand(context.TODO(), reconfigCommand).Decode(&reconfigResult); err != nil {
		if !strings.Contains(err.Error(), "configuration version must be greater than old") {
			log.Println(err)
		}
	}
	if err := adminDatabase.RunCommand(context.TODO(), statusCommand).Decode(&statusResult); err != nil {
		log.Println(err)
	}
	if !(reconfigResult.Ok == 1 || statusResult.Primary() != nil) {
		log.Println(initiateResult)
		log.Println(reconfigResult)
		log.Println(statusResult)
		return fmt.Errorf("Reconfig failed")
	}

	err = client.Ping(mctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("Could not ping database within 10 seconds: %s (%s:%s)", err.Error(), databaseConnectionURI, databaseName)
	}
	return nil
}

func main() {
	cfg := &mongoDBConfig{
		Host:       "localhost",
		Port:       27017,
		User:       "example",
		Database:   "mydb",
		Password:   "123",
		ReplicaSet: "my-replicaset",
	}
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("Sucessfully connected to replicaset")
}
