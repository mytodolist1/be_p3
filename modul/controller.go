package modul

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aiteung/atdb"
	"github.com/badoux/checkmail"
	model "github.com/mytodolist1/be_p3/model"
)

func MongoConnect(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		// DBString: "mongodb+srv://admin:admin@projectexp.pa7k8.gcp.mongodb.net", //os.Getenv(MONGOCONNSTRINGENV),
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func InsertOneDoc(db *mongo.Database, col string, docs interface{}) (insertedID primitive.ObjectID, err error) {
	cols := db.Collection(col)
	result, err := cols.InsertOne(context.Background(), docs)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, err
}

// user
func Register(db *mongo.Database, col string, userdata model.User) error {
	if userdata.Username == "" || userdata.Password == "" || userdata.Email == "" {
		return fmt.Errorf("Data tidak lengkap")
	}

	// Periksa apakah email valid
	if err := checkmail.ValidateFormat(userdata.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	}

	// Periksa apakah email dan username sudah terdaftar
	userExists, _ := GetUserFromEmail(db, col, userdata.Email)
	if userExists.Email != "" {
		return fmt.Errorf("Email sudah terdaftar")
	}
	userExists, _ = GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username != "" {
		return fmt.Errorf("Username sudah terdaftar")
	}

	// Periksa apakah password memenuhi syarat
	if len(userdata.Password) < 6 {
		return fmt.Errorf("Password minimal 6 karakter")
	}
	if strings.Contains(userdata.Password, " ") {
		return fmt.Errorf("Password tidak boleh mengandung spasi")
	}

	// Periksa apakah username memenuhi syarat
	if strings.Contains(userdata.Username, " ") {
		return fmt.Errorf("Username tidak boleh mengandung spasi")
	}

	// Simpan pengguna ke basis data
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

	// Periksa apakah pengguna dengan username tertentu ada
	userExists, _ := GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username == "" {
		err = fmt.Errorf("Username tidak ditemukan")
		return user, false, err
	}

	// Periksa apakah kata sandi benar
	if !CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("Password salah")
		return user, false, err
	}
	return userExists, true, nil
}

func UpdateUser(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	if userdata.Username == "" || userdata.Email == "" {
		err = fmt.Errorf("Data tidak boleh kosong")
		return user, false, err
	}

	userExists, err := GetUserFromID(db, col, userdata.ID)
	if err != nil {
		return user, false, err
	}

	// Periksa apakah data yang akan diupdate sama dengan data yang sudah ada
	if userdata.Username == userExists.Username && userdata.Email == userExists.Email {
		err = fmt.Errorf("Data yang ingin diupdate tidak boleh sama")
		return user, false, err
	}

	checkmail.ValidateFormat(userdata.Email)
	if err != nil {
		err = fmt.Errorf("Email tidak valid")
		return user, false, err
	}

	// Periksa apakah username memenuhi syarat
	if strings.Contains(userdata.Username, " ") {
		err = fmt.Errorf("Username tidak boleh mengandung spasi")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	filter := bson.M{"_id": userdata.ID}
	update := bson.M{
		"$set": bson.M{
			"email":    userdata.Email,
			"username": userdata.Username,
		},
	}
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return user, false, err
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("Data tidak berhasil diupdate")
		return user, false, err
	}
	return user, true, nil
}

func ChangePassword(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	// Periksa apakah pengguna dengan username tertentu ada
	userExists, err := GetUserFromUsername(db, col, userdata.Username)
	if err != nil {
		return user, false, err
	}

	// Periksa apakah password memenuhi syarat
	if userdata.Password == "" {
		err = fmt.Errorf("Password tidak boleh kosong")
		return user, false, err
	}

	if len(userdata.Password) < 6 {
		err = fmt.Errorf("Password minimal 6 karakter")
		return user, false, err
	}

	if strings.Contains(userdata.Password, " ") {
		err = fmt.Errorf("Password tidak boleh mengandung spasi")
		return user, false, err
	}

	// Periksa apakah password sama dengan password lama
	if CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("Password tidak boleh sama")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	userExists.Password = hash
	filter := bson.M{"username": userdata.Username}
	update := bson.M{
		"$set": bson.M{
			"password": userExists.Password,
		},
	}
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return user, false, err
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("PAssword tidak berhasil diupdate")
		return user, false, err
	}
	return user, true, nil
}

func DeleteUser(db *mongo.Database, col string, userdata model.User) (status bool, err error) {
	_, err = GetUserFromUsername(db, col, userdata.Username)
	if err != nil {
		err = fmt.Errorf("Username tidak ditemukan")
		return false, err
	}

	filter := bson.M{"username": userdata.Username}
	cols := db.Collection(col)

	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		err = fmt.Errorf("Error deleting document: %v", err)
		return false, err
	}

	if result.DeletedCount == 0 {
		return false, fmt.Errorf("Failed to delete user")
	}

	return true, nil
}

func GetUserFromID(db *mongo.Database, col string, _id primitive.ObjectID) (user model.User, err error) {
	// userID := _id.Hex()

	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err := fmt.Errorf("no data found for ID %s", _id)
			return user, err
		}
		err := fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
		return user, err
	}
	return user, nil
}

func GetUserFromUsername(db *mongo.Database, col string, username string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"username": username}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err := fmt.Errorf("no data found for username %s", username)
			return user, err
		}
		err := fmt.Errorf("error retrieving data for username %s: %s", username, err.Error())
		return user, err
	}
	return user, nil
}

func GetUserFromEmail(db *mongo.Database, col string, email string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"email": email}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err := fmt.Errorf("no data found for email %s", email)
			return user, err
		}
		err := fmt.Errorf("error retrieving data for email %s: %s", email, err.Error())
		return user, err
	}
	return user, nil
}

func GetAllUser(db *mongo.Database, col string) (userlist []model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{}

	cur, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetAllUser in colection", col, ":", err)
		return userlist, err
	}

	err = cur.All(context.Background(), &userlist)
	if err != nil {
		fmt.Println("Error reading documents:", err)
		return userlist, err
	}

	return userlist, nil
}

// todo
func InsertTodo(db *mongo.Database, col string, todo model.Todo) (insertedID primitive.ObjectID, err error) {
	todo.TimeStamp.CreatedAt = time.Now()
	todo.TimeStamp.UpdatedAt = time.Now()

	insertedID, err = InsertOneDoc(db, col, todo)
	if err != nil {
		fmt.Printf("InsertTodo: %v\n", err)
	}
	return insertedID, nil
}

func GetTodoFromID(db *mongo.Database, col string, _id primitive.ObjectID) (todo model.Todo, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	err = cols.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no data found for ID", _id)
		} else {
			fmt.Println("error retrieving data for ID", _id, ":", err.Error())
		}
	}
	return todo, nil
}

func GetTodoList(db *mongo.Database, col string) (todo []model.Todo, err error) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetTodoList in colection", col, ":", err)
		return nil, err
	}
	err = cursor.All(context.Background(), &todo)
	if err != nil {
		fmt.Println(err)
	}
	return todo, nil
}

func UpdateTodo(db *mongo.Database, col string, todo model.Todo) (todos model.Todo, status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": todo.ID}
	update := bson.M{
		"$set": bson.M{
			"title":               todo.Title,
			"description":         todo.Description,
			"deadline":            todo.Deadline,
			"timestamp.updatedat": time.Now(),
		},
		"$setOnInsert": bson.M{
			"timestamp.createdat": todo.TimeStamp.CreatedAt,
		},
	}

	options := options.Update().SetUpsert(true)

	result, err := cols.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return todos, false, err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		err = fmt.Errorf("Data tidak berhasil diupdate")
		return todos, false, err
	}

	err = cols.FindOne(context.Background(), filter).Decode(&todos)
	if err != nil {
		return todos, false, err
	}

	return todos, true, nil
}

func DeleteTodo(db *mongo.Database, col string, _id primitive.ObjectID) (status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("Data tidak berhasil dihapus")
		return false, err
	}
	return true, nil
}
