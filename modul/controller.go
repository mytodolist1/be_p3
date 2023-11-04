package modul

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/badoux/checkmail"
	model "github.com/mytodolist1/be_p3/model"
)

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

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
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs
}

// user
func Register(db *mongo.Database, col string, userdata model.User) error {
	if userdata.Username == "" || userdata.Password == "" || userdata.Email == "" {
		return fmt.Errorf("Data tidak lengkap")
	}
	if err := checkmail.ValidateFormat(userdata.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	}
	userExists, _ := GetUserFromEmail(db, col, userdata.Email)
	if userExists.Email != "" {
		return fmt.Errorf("Email sudah terdaftar")
	}
	userExists, _ = GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username != "" {
		return fmt.Errorf("Username sudah terdaftar")
	}
	if len(userdata.Password) < 6 {
		return fmt.Errorf("Password minimal 6 karakter")
	}
	if strings.Contains(userdata.Password, " ") {
		return fmt.Errorf("Password tidak boleh mengandung spasi")
	}
	if strings.Contains(userdata.Username, " ") {
		return fmt.Errorf("Username tidak boleh mengandung spasi")
	}
	hash, _ := HashPassword(userdata.Password)
	user := bson.M{
		"_id":      primitive.NewObjectID(),
		"email":    userdata.Email,
		"username": userdata.Username,
		"password": hash,
		"role":     "user",
	}
	_, err := InsertOneDoc(db, col, user)
	if err != nil {
		return fmt.Errorf("SignUp: %v", err)
	}
	return nil
}

func LogIn(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	if userdata.Username == "" || userdata.Password == "" {
		err = fmt.Errorf("Data tidak lengkap")
		return user, false, err
	}
	userExists, _ := GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username == "" {
		err = fmt.Errorf("Username tidak ditemukan")
		return user, false, err
	}
	if !CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("Password salah")
		return user, false, err
	}
	return userExists, true, nil
}

func GetUserFromID(db *mongo.Database, col string, _id primitive.ObjectID) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromID: %v\n", err)
	}
	return user, nil
}

func GetUserFromUsername(db *mongo.Database, col string, username string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"username": username}
	err = cols.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromUsername: \n%v", err)
		// return user, err
	}
	return user, nil
}

func GetUserByUsername(db *mongo.Database, col string, username string) (user model.User, err error) {
	allUsers := GetAllUser(db, col)
	for _, u := range allUsers {
		if u.Username == username {
			return u, nil
		}
	}
	return user, fmt.Errorf("User with username %s not found", username)
}

func GetUserFromEmail(db *mongo.Database, col string, email string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"email": email}
	err = cols.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromEmail: %v\n", err)
	}
	return user, nil
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
		fmt.Println("Error GetTodoList in colection", col, ":", err)
	}
	err = cursor.All(context.Background(), &todolist.Items)
	if err != nil {
		fmt.Println(err)
	}
	return todolist
}
