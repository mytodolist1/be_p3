package modul

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	model "github.com/mytodolist1/be_p3/model"
)

// func MongoConnect(MongoString, dbname string) *mongo.Database {
// 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
// 	if err != nil {
// 		fmt.Printf("MongoConnect: %v\n", err)
// 	}
// 	return client.Database(dbname)
// }

func InsertOneDoc(db *mongo.Database, col string, docs interface{}) (insertedID primitive.ObjectID, err error) {
	cols := db.Collection(col)
	result, err := cols.InsertOne(context.Background(), docs)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return
}

func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs
}

func updateOneDoc(db *mongo.Database, id primitive.ObjectID, col string, docs interface{}) (err error) {
	cols := db.Collection(col)
	filter := bson.M{}
	update := bson.M{"$set": docs}
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("updateOneDoc: %v\n", err)
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}
	return
}

func deleteOneDoc(db *mongo.Database, id primitive.ObjectID, col string) (err error) {
	cols := db.Collection(col)
	filter := bson.M{}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Printf("deleteOneDoc: %v %v\n", id, err)
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("no data has been deleted with the specified id")
		return
	}
	return
}

// user
func InsertUser(db *mongo.Database, col string, userdata model.User) (insertedID primitive.ObjectID, err error) {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	insertedID, err = InsertOneDoc(db, col, userdata)
	if err != nil {
		fmt.Printf("InsertUser: %v\n", err)
	}
	return insertedID, err
}

func GetUserFromUsername(db *mongo.Database, col string, username string) (user model.User) {
	cols := db.Collection(col)
	filter := bson.M{"username": username}
	err := cols.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromUsername: %v\n", err)
	}
	return user
}

func GetAllUser(db *mongo.Database, col string) (userlist []model.User) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &userlist)
	if err != nil {
		fmt.Println(err)
	}
	return userlist
}

// todo
func InsertTodo(db *mongo.Database, col string, todo model.Todo) (insertedID primitive.ObjectID, err error) {
	insertedID, err = InsertOneDoc(db, col, todo)
	if err != nil {
		fmt.Printf("InsertTodo: %v\n", err)
	}
	return insertedID, err
}

func GetTodoFromID(db *mongo.Database, col string, id primitive.ObjectID) (todo model.Todo) {
	cols := db.Collection(col)
	filter := bson.M{"_id": id}
	err := cols.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		fmt.Printf("GetTodoFromID: %v\n", err)
	}
	return todo
}

func GetTodoList(db *mongo.Database, col string) (todolist model.TodoList) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.Background(), &todolist.Items)
	if err != nil {
		fmt.Println(err)
	}
	return todolist
}
