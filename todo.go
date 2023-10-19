package bep3

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/mongo"
)

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func InsertUser(db *mongo.Database, collection string, userdata User) string {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(db, collection, userdata)
	return "Username : " + userdata.Username + "\nPassword : " + userdata.Password
}

// var MongoString string = os.Getenv("MONGOSTRING")

// var MongoInfo = atdb.DBInfo{
// 	DBString: MongoString,
// 	DBName:   "db_proyek3",
// }

// var MongoConn = atdb.MongoConnect(MongoInfo)

// func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
// 	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
// 	if err != nil {
// 		fmt.Printf("InsertOneDoc: %v\n", err)
// 	}
// 	return insertResult.InsertedID
// }

// // login
// func LoginUser(db *mongo.Database, col string, username string, password string) (authenticated bool, err error) {
// 	filter := bson.M{
// 		"username": username,
// 		"password": password,
// 	}

// 	result, err := db.Collection(col).CountDocuments(context.Background(), filter)
// 	if err != nil {
// 		fmt.Printf("LoginUser: %v\n", err)
// 		return false, err
// 	}

// 	if result == 0 {
// 		return true, nil
// 	}

// 	return false, nil
// }

// func InsertUser(db *mongo.Database, col string, username string, password string) (insertedID primitive.ObjectID, err error) {
// 	user := bson.M{
// 		"username": username,
// 		"password": password,
// 	}
// 	result, err := db.Collection(col).InsertOne(context.Background(), user)
// 	if err != nil {
// 		fmt.Printf("InsertUser: %v\n", err)
// 		return
// 	}
// 	insertedID = result.InsertedID.(primitive.ObjectID)
// 	return insertedID, nil
// }
