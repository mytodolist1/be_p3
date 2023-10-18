package modul

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	model "github.com/mytodolist1/be_p3/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	}

	var admin model.Admin
	err = db.Collection(col).FindOne(context.Background(), filter).Decode(&admin)
	if err != nil {
		fmt.Printf("LoginAdmin: %v\n", err)
		return false, err
	}

	// Verifikasi kata sandi menggunakan bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err == nil {
		return true, nil // Kata sandi cocok, otentikasi berhasil
	}

	return false, nil // Kata sandi tidak cocok, otentikasi gagal
}
