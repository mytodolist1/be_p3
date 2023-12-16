package modul

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aiteung/atdb"
	"github.com/badoux/checkmail"
	"github.com/mytodolist1/be_p3/model"
)

func MongoConnect(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
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
func GenerateUID(len int) (string, error) {
	bytes := make([]byte, len)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func Register(db *mongo.Database, col string, userdata model.User) error {
	if userdata.Email == "" || userdata.Username == "" || userdata.Password == "" || userdata.ConfirmPassword == "" {
		return fmt.Errorf("Data tidak lengkap")
	}

	// Periksa apakah email valid
	err := checkmail.ValidateFormat(userdata.Email)
	if err != nil {
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

	// Periksa apakah password dan konfirmasi password sama
	if userdata.Password != userdata.ConfirmPassword {
		return fmt.Errorf("Password dan konfirmasi password tidak sama")
	}

	// uid := GenerateUID(&userdata)

	uid, err := GenerateUID(8)
	if err != nil {
		return fmt.Errorf("GenerateUID: %v", err)
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	user := bson.D{
		{Key: "_id", Value: primitive.NewObjectID()},
		{Key: "uid", Value: uid},
		{Key: "email", Value: userdata.Email},
		{Key: "username", Value: userdata.Username},
		{Key: "password", Value: hash},
		{Key: "role", Value: "user"},
	}

	_, err = InsertOneDoc(db, col, user)
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

	// // Periksa apakah role benar
	// if userdata.Role != userExists.Role {
	// 	err = fmt.Errorf("Role tidak sesuai")
	// 	return user, false, err
	// }

	return userExists, true, nil
}

// update user with log
func UpdateUser(db *mongo.Database, col string, userdata model.User) (model.User, bool, error) {
	if userdata.Username == "" || userdata.Email == "" {
		err := fmt.Errorf("Data tidak boleh kosong")
		return model.User{}, false, err
	}

	userExists, err := GetUserFromID(db, col, userdata.ID)
	if err != nil {
		return model.User{}, false, err
	}

	// Periksa apakah data yang akan diupdate sama dengan data yang sudah ada
	if userdata.Username == userExists.Username && userdata.Email == userExists.Email {
		err = fmt.Errorf("Data yang ingin diupdate tidak boleh sama")
		return model.User{}, false, err
	}

	checkmail.ValidateFormat(userdata.Email)
	if err != nil {
		err = fmt.Errorf("Email tidak valid")
		return model.User{}, false, err
	}

	// Periksa apakah username memenuhi syarat
	if strings.Contains(userdata.Username, " ") {
		err = fmt.Errorf("Username tidak boleh mengandung spasi")
		return model.User{}, false, err
	}

	cols := db.Collection(col)

	// Simpan pengguna ke basis data
	filter := bson.M{"_id": userdata.ID}

	var originalUser model.User
	err = cols.FindOne(context.Background(), filter).Decode(&originalUser)
	if err != nil {
		return model.User{}, false, err
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "email", Value: userdata.Email},
			{Key: "username", Value: userdata.Username},
		}},
	}

	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	result := cols.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&userdata)
	if result != nil {
		return model.User{}, false, err
	}

	err = LogUser(db, "loguser", userdata.ID, originalUser, userdata)
	if err != nil {
		// Handle error logging
		return model.User{}, false, err
	}

	return userdata, true, nil
}

func ChangePassword(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	// Periksa apakah pengguna dengan username tertentu ada
	userExists, err := GetUserFromUsername(db, col, userdata.Username)
	if err != nil {
		return user, false, err
	}

	// Periksa apakah password memenuhi syarat
	if userdata.Password == "" || userdata.ConfirmPassword == "" {
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

	// Periksa apakah password dan konfirmasi password sama
	if userdata.Password != userdata.ConfirmPassword {
		err = fmt.Errorf("Password dan konfirmasi password tidak sama")
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

func DeleteUser(db *mongo.Database, col string, userdata model.User) (bool, error) {
	_, err := GetUserFromUsername(db, col, userdata.Username)
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

func GetUserFromToken(db *mongo.Database, col string, uid string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"uid": uid}

	err = cols.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no data found for ID", uid)
		} else {
			fmt.Println("error retrieving data for ID", uid, ":", err.Error())
		}
	}

	return user, nil
}

// admin
func GetUserFromRole(db *mongo.Database, col string, role string) (userlist []model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"role": role}

	cur, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetUserFromRole in colection", col, ":", err)
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
func InsertTodo(db *mongo.Database, col string, todoDoc model.Todo, uid string) (insertedID primitive.ObjectID, err error) {
	if todoDoc.Title == "" || todoDoc.Description == "" || todoDoc.Deadline == "" || todoDoc.Time == "" {
		err = fmt.Errorf("Data tidak boleh kosong")
		return insertedID, err
	}

	objectId := primitive.NewObjectID()

	time := time.Now().UnixMilli()
	fmt.Println(time)

	todo := bson.D{
		{Key: "_id", Value: objectId},
		{Key: "title", Value: todoDoc.Title},
		{Key: "description", Value: todoDoc.Description},
		{Key: "deadline", Value: todoDoc.Deadline},
		{Key: "time", Value: todoDoc.Time},
		{Key: "timestamps", Value: bson.D{
			{Key: "createdat", Value: time},
			{Key: "updatedat", Value: time},
		}},
		{Key: "user", Value: bson.D{
			{Key: "uid", Value: uid},
		}},
	}

	insertedID, err = InsertOneDoc(db, col, todo)
	if err != nil {
		fmt.Printf("InsertTodo: %v\n", err)
	}

	return insertedID, nil
}

// update todo with log
func UpdateTodo(db *mongo.Database, col string, todo model.Todo) (model.Todo, bool, error) {
	if todo.Title == "" || todo.Description == "" || todo.Deadline == "" || todo.Time == "" {
		err := fmt.Errorf("Data tidak lengkap")
		return model.Todo{}, false, err
	}

	todoExists, err := GetTodoFromID(db, col, todo.ID)
	if err != nil {
		return model.Todo{}, false, err
	}

	// Periksa apakah data yang akan diupdate sama dengan data yang sudah ada
	if todo.Title == todoExists.Title && todo.Description == todoExists.Description && todo.Deadline == todoExists.Deadline && todo.Time == todoExists.Time {
		err = fmt.Errorf("Silahkan update data anda")
		return model.Todo{}, false, err
	}

	cols := db.Collection(col)
	filter := bson.M{"_id": todo.ID}

	var originalTodo model.Todo
	err = cols.FindOne(context.Background(), filter).Decode(&originalTodo)
	if err != nil {
		return model.Todo{}, false, err
	}

	time := time.Now().UnixMilli()

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: todo.Title},
			{Key: "description", Value: todo.Description},
			{Key: "deadline", Value: todo.Deadline},
			{Key: "time", Value: todo.Time},
			{Key: "timestamps.updatedat", Value: time},
		}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "timestamps.createdat", Value: todo.TimeStamps.CreatedAt},
		}},
	}

	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	result := cols.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&todo)
	if result != nil {
		return model.Todo{}, false, err
	}

	err = LogTodo(db, "logtodo", todo.ID, originalTodo, todo)
	if err != nil {
		// Handle error logging
		return model.Todo{}, false, err
	}

	return todo, true, nil
}

func DeleteTodo(db *mongo.Database, col string, _id primitive.ObjectID) (bool, error) {
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

// isDone
func TodoClear(db *mongo.Database, col string, done model.TodoClear) (bool, error) {
	cols := db.Collection(col)

	todo, err := GetTodoFromID(db, "todo", done.Todo.ID)
	if err != nil {
		fmt.Println("Error GetTodoFromID in colection", col, ":", err)
		return false, err
	}

	fmt.Println("todo: ", todo)

	time := time.Now().UnixMilli()

	insert := bson.D{
		{Key: "isdone", Value: true},
		{Key: "timeclear", Value: time},
		{Key: "todo", Value: bson.D{
			{Key: "_id", Value: todo.ID},
			{Key: "title", Value: todo.Title},
			{Key: "description", Value: todo.Description},
			{Key: "deadline", Value: todo.Deadline},
			{Key: "time", Value: todo.Time},
			{Key: "user", Value: bson.D{
				{Key: "uid", Value: todo.User.UID},
			}},
		}},
	}

	fmt.Println("insert: ", insert)

	_, err = cols.InsertOne(context.Background(), insert)
	if err != nil {
		return false, err
	}

	status, err := DeleteTodo(db, "todo", todo.ID)
	if err != nil {
		fmt.Println("Error DeleteTodo in colection", col, ":", err)
		return false, err
	}

	if !status {
		fmt.Println("Data tidak berhasil dihapus")
		return false, err
	}

	return true, nil
}

func GetTodoDone(db *mongo.Database, col string) (todo []model.TodoClear, err error) {
	cols := db.Collection(col)
	filter := bson.M{"isdone": true}

	cur, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetTodoDone in colection", col, ":", err)
		return todo, err
	}

	err = cur.All(context.Background(), &todo)
	if err != nil {
		fmt.Println("Error reading documents:", err)
		return todo, err
	}

	return todo, nil
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

func GetTodoFromToken(db *mongo.Database, col string, uid string) (todo []model.Todo, err error) {
	cols := db.Collection(col)
	filter := bson.M{"user.uid": uid}

	cursor, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetTodoFromToken in colection", col, ":", err)
		return nil, err
	}

	err = cursor.All(context.Background(), &todo)
	if err != nil {
		fmt.Println(err)
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

// log
func LogTodo(db *mongo.Database, col string, todoID primitive.ObjectID, originalTodo model.Todo, updatedTodo model.Todo) error {
	cols := db.Collection(col)

	// Filter untuk menemukan dokumen log dengan ID yang sesuai
	filter := bson.M{"todoid": todoID.Hex()}

	// Data baru yang akan ditambahkan ke dalam array Change
	time := time.Now().UnixMilli()

	newChangeEntry := bson.D{
		{Key: "timestamp", Value: time},
		{Key: "dataold", Value: originalTodo},
		{Key: "datanew", Value: updatedTodo},
	}

	// Setelah menemukan dokumen log, perbarui array Change dengan $addToSet
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "change", Value: newChangeEntry},
		}},
	}

	// Jalankan operasi update
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
		return err
	}

	// Jika dokumen log belum ada, tambahkan dokumen baru
	if result.MatchedCount == 0 {
		logUpdate := bson.D{
			{Key: "timestamp", Value: time},
			{Key: "action", Value: "update"},
			{Key: "todoid", Value: todoID.Hex()},
			{Key: "userid", Value: updatedTodo.User.UID},
			{Key: "change", Value: []map[string]interface{}{
				{
					"timestamp": time,
					"dataold":   originalTodo,
					"datanew":   updatedTodo,
				},
			}},
		}

		_, err := cols.InsertOne(context.Background(), logUpdate)
		if err != nil {
			fmt.Printf("InsertOneDoc: %v\n", err)
			return err
		}
	} else {
		// Jika dokumen log sudah ada, tampilkan dokumen log yang sudah diperbarui
		err = cols.FindOne(context.Background(), filter).Decode(&originalTodo)
		if err != nil {
			fmt.Printf("FindOneDoc: %v\n", err)
			return err
		}
	}

	fmt.Printf("Result (Before Update): %+v\n", originalTodo)
	fmt.Printf("Result (After Update): %+v\n", updatedTodo)

	return nil
}

func LogUser(db *mongo.Database, col string, id primitive.ObjectID, originalUser model.User, updatedUser model.User) error {
	cols := db.Collection(col)

	// Filter untuk menemukan dokumen log dengan ID yang sesuai
	filter := bson.M{"id": id.Hex()}

	// Data baru yang akan ditambahkan ke dalam array Change
	time := time.Now().UnixMilli()

	newChangeEntry := bson.D{
		{Key: "timestamp", Value: time},
		{Key: "dataold", Value: originalUser},
		{Key: "datanew", Value: updatedUser},
	}

	// Setelah menemukan dokumen log, perbarui array Change dengan $addToSet
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "change", Value: newChangeEntry},
		}},
	}

	// Jalankan operasi update
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
		return err
	}

	// Jika dokumen log belum ada, tambahkan dokumen baru
	if result.MatchedCount == 0 {
		logUpdate := bson.D{
			{Key: "timestamp", Value: time},
			{Key: "action", Value: "update"},
			{Key: "id", Value: id.Hex()},
			{Key: "userid", Value: updatedUser.UID},
			{Key: "change", Value: []map[string]interface{}{
				{
					"timestamp": time,
					"dataold":   originalUser,
					"datanew":   updatedUser,
				},
			}},
		}

		_, err := cols.InsertOne(context.Background(), logUpdate)
		if err != nil {
			fmt.Printf("InsertOneDoc: %v\n", err)
			return err
		}
	} else {
		// Jika dokumen log sudah ada, tampilkan dokumen log yang sudah diperbarui
		err = cols.FindOne(context.Background(), filter).Decode(&originalUser)
		if err != nil {
			fmt.Printf("FindOneDoc: %v\n", err)
			return err
		}
	}

	fmt.Printf("Result (Before Update): %+v\n", originalUser)
	fmt.Printf("Result (After Update): %+v\n", updatedUser)

	return nil
}

// get log
func GetLogTodoList(db *mongo.Database, col string) (log []model.LogTodo, err error) {
	cols := db.Collection(col)
	filter := bson.M{}

	cur, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetLogTodoList in colection", col, ":", err)
		return log, err
	}

	err = cur.All(context.Background(), &log)
	if err != nil {
		fmt.Println("Error reading documents:", err)
		return log, err
	}

	return log, nil
}

func GetLogTodoFromUID(db *mongo.Database, col, userid string) (log []model.LogTodo, err error) {
	cols := db.Collection(col)
	filter := bson.M{"userid": userid}

	cur, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetLogTodoList in colection", col, ":", err)
		return log, err
	}

	err = cur.All(context.Background(), &log)
	if err != nil {
		fmt.Println("Error reading documents:", err)
		return log, err
	}

	return log, nil
}

// not used yet
func GetLogAllUser(db *mongo.Database, col string) (log []model.LogUser, err error) {
	cols := db.Collection(col)
	filter := bson.M{}

	cur, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetLogAllUser in colection", col, ":", err)
		return log, err
	}

	err = cur.All(context.Background(), &log)
	if err != nil {
		fmt.Println("Error reading documents:", err)
		return log, err
	}

	return log, nil
}

// not used yet
func GetLogUserFromUID(db *mongo.Database, col, userid string) (log []model.LogUser, err error) {
	cols := db.Collection(col)
	filter := bson.M{"userid": userid}

	err = cols.FindOne(context.Background(), filter).Decode(&log)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no data found for ID", userid)
		} else {
			fmt.Println("error retrieving data for ID", userid, ":", err.Error())
		}
	}

	return log, nil
}

// // update user
// func UpdateUser(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
// 	if userdata.Username == "" || userdata.Email == "" {
// 		err = fmt.Errorf("Data tidak boleh kosong")
// 		return user, false, err
// 	}

// 	userExists, err := GetUserFromID(db, col, userdata.ID)
// 	if err != nil {
// 		return user, false, err
// 	}

// 	// Periksa apakah data yang akan diupdate sama dengan data yang sudah ada
// 	if userdata.Username == userExists.Username && userdata.Email == userExists.Email {
// 		err = fmt.Errorf("Data yang ingin diupdate tidak boleh sama")
// 		return user, false, err
// 	}

// 	checkmail.ValidateFormat(userdata.Email)
// 	if err != nil {
// 		err = fmt.Errorf("Email tidak valid")
// 		return user, false, err
// 	}

// 	// Periksa apakah username memenuhi syarat
// 	if strings.Contains(userdata.Username, " ") {
// 		err = fmt.Errorf("Username tidak boleh mengandung spasi")
// 		return user, false, err
// 	}

// 	// Simpan pengguna ke basis data
// 	filter := bson.M{"_id": userdata.ID}
// 	update := bson.D{
// 		{Key: "$set", Value: bson.D{
// 			{Key: "email", Value: userdata.Email},
// 			{Key: "username", Value: userdata.Username},
// 		}},
// 	}

// 	cols := db.Collection(col)
// 	result, err := cols.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		return user, false, err
// 	}

// 	if result.ModifiedCount == 0 {
// 		err = fmt.Errorf("Data tidak berhasil diupdate")
// 		return user, false, err
// 	}

// 	return user, true, nil
// }

// func GetAllUser(db *mongo.Database, col string) (userlist []model.User, err error) {
// 	cols := db.Collection(col)
// 	filter := bson.M{}

// 	cur, err := cols.Find(context.Background(), filter)
// 	if err != nil {
// 		fmt.Println("Error GetAllUser in colection", col, ":", err)
// 		return userlist, err
// 	}

// 	err = cur.All(context.Background(), &userlist)
// 	if err != nil {
// 		fmt.Println("Error reading documents:", err)
// 		return userlist, err
// 	}

// 	return userlist, nil
// }

// // update todo
// func UpdateTodo(db *mongo.Database, col string, todo model.Todo) (todos model.Todo, status bool, err error) {
// 	if todo.Title == "" || todo.Description == "" || todo.Deadline == "" || todo.Time == "" {
// 		err = fmt.Errorf("Data tidak lengkap")
// 		return todos, false, err
// 	}

// 	cols := db.Collection(col)
// 	filter := bson.M{"_id": todo.ID}
// 	update := bson.D{
// 		{Key: "$set", Value: bson.D{
// 			{Key: "title", Value: todo.Title},
// 			{Key: "description", Value: todo.Description},
// 			{Key: "deadline", Value: todo.Deadline},
// 			{Key: "timestamps.updatedat", Value: time.Now()},
// 		}},
// 		{Key: "$setOnInsert", Value: bson.D{
// 			{Key: "timestamps.createdat", Value: todo.TimeStamps.CreatedAt},
// 		}},
// 	}

// 	options := options.Update().SetUpsert(true)

// 	result, err := cols.UpdateOne(context.Background(), filter, update, options)
// 	if err != nil {
// 		return todos, false, err
// 	}

// 	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
// 		err = fmt.Errorf("Data tidak berhasil diupdate")
// 		return todos, false, err
// 	}

// 	err = cols.FindOne(context.Background(), filter).Decode(&todos)
// 	if err != nil {
// 		return todos, false, err
// 	}

// 	return todos, true, nil
// }
