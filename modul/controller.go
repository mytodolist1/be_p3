package modul

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoString string = os.Getenv("MONGOSTRING")

var MongoInfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "db_proyek3",
}

var MongoConn = atdb.MongoConnect(MongoInfo)

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

// login
func LoginAdmin(db *mongo.Database, col string, username string, password string) (authenticated bool, err error) {
	filter := bson.M{
		"username": username,
		"password": password,
	}

	result, err := db.Collection(col).CountDocuments(context.Background(), filter)
	if err != nil {
		fmt.Printf("LoginAdmin: %v\n", err)
		return false, err
	}

	if result == 0 {
		return true, nil
	}

	return false, nil
}

func InsertAdmin(db *mongo.Database, col string, username string, password string) (insertedID primitive.ObjectID, err error) {
	admin := bson.M{
		"username": username,
		"password": password,
	}
	result, err := db.Collection(col).InsertOne(context.Background(), admin)
	if err != nil {
		fmt.Printf("InsertAdmin: %v\n", err)
		return
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}
