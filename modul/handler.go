package modul

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	model "github.com/mytodolist1/be_p3/model"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Responsed model.Credential
	Response  model.TodoResponse
	datauser  model.User
	datatodo  model.Todo
)

// user
func GCFHandlerGetAllUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	userlist, err := GetAllUser(mconn, collectionname)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Get User Success"
	Responsed.Data = userlist

	return GCFReturnStruct(Responsed)
}

func GCFHandlerGetUserByUsername(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	user, err := GetUserFromUsername(mconn, collectionname, username)
	if err != nil {
		Responsed.Message = "Error retrieving user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Hello user"
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerRegister(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	err = Register(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Register success"

	return GCFReturnStruct(Responsed)
}

func GCFHandlerLogIn(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := LogIn(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	tokenstring, err := watoken.Encode(datauser.UID, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Responsed.Message = "Gagal Encode Token :" + err.Error()

	} else {
		Responsed.Message = "Selamat Datang" + " " + user.Username
		Responsed.Token = tokenstring
	}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerUpdateUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	// userlogin := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// if userlogin != existingUser.UID {
	// 	Responsed.Message = "Unauthorized access: User mismatch" + userlogin + " " + existingUser.UID
	// 	return GCFReturnStruct(Responsed)
	// }

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		Responsed.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Responsed.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.ID = ID

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := UpdateUser(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error updating user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Update success " + " " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerChangePassword(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	// userlogin := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// if userlogin != existingUser.UID {
	// 	Responsed.Message = "Unauthorized access: User mismatch" + userlogin + " " + existingUser.UID
	// 	return GCFReturnStruct(Responsed)
	// }

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := ChangePassword(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error changing password: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Password change success for user" + " " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerDeleteUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	// userlogin := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// if userlogin != existingUser.UID {
	// 	Responsed.Message = "Unauthorized access: User mismatch" + userlogin + " " + existingUser.UID
	// 	return GCFReturnStruct(Responsed)
	// }

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	_, err = DeleteUser(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error deleting user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Delete user " + datauser.Username + " success"

	return GCFReturnStruct(Responsed)
}

// todo
// func GCFHandlerGetTodoList(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	Response.Status = false

// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		Response.Message = "error parsing application/json1:"
// 		return GCFReturnStruct(Response)
// 	}

// 	userInfo := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)
// 	// if err != nil {
// 	// 	Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
// 	// 	return GCFReturnStruct(Response)
// 	// }

// 	if userInfo != datauser.Username {
// 		Response.Message = "Unauthorized access: User mismatch"
// 		return GCFReturnStruct(Response)
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&datatodo)
// 	if err != nil {
// 		Response.Message = "error parsing application/json3: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}

// 	todolist, err := GetTodoList(mconn, collectionname)
// 	if err != nil {
// 		Response.Message = err.Error()
// 		return GCFReturnStruct(Response)
// 	}

// 	Response.Status = true
// 	Response.Message = "Get todolist success"
// 	Response.Data = todolist

// 	return GCFReturnStruct(Response)
// }

func GCFHandlerGetTodoListByUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// userInfo := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// if userInfo != existingUser.UID {
	// 	Response.Message = "Unauthorized access: User mismatch"
	// 	return GCFReturnStruct(Response)
	// }

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	datauser.UID = userInfo.Id

	username := r.URL.Query().Get("username")
	if username == "" {
		Response.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	datauser.Username = username

	err = json.NewDecoder(r.Body).Decode(&datatodo)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	todo, err := GetTodoFromUsername(mconn, collectionname, datauser.Username)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Get todo success"
	Response.Data = todo

	return GCFReturnStruct(Response)
}

func GCFHandlerGetTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// userInfo := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// if userInfo != existingUser.UID {
	// 	Response.Message = "Unauthorized access: User mismatch"
	// 	return GCFReturnStruct(Response)
	// }

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	datauser.UID = userInfo.Id

	id := r.URL.Query().Get("_id")
	if id == "" {
		Response.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	datatodo.ID = ID

	err = json.NewDecoder(r.Body).Decode(&datatodo)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	todo, err := GetTodoFromID(mconn, collectionname, datatodo.ID)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Get todo success"
	Response.Data = []model.Todo{todo}

	return GCFReturnStruct(Response)
}

func GCFHandlerInsertTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	// token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// userInfo := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// datauser.UID = userInfo

	userInfo, _ := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	// if err != nil {
	// 	Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
	// 	return GCFReturnStruct(Response)
	// }

	if userInfo.Id != datauser.UID {
		Response.Message = "Unauthorized access: User mismatch" + ", " + datauser.UID + ", " + userInfo.Id
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&datatodo)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	_, err = InsertTodo(mconn, collectionname, datatodo, userInfo.Id)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Insert todo success"

	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// userInfo := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// if userInfo != existingUser.UID {
	// 	Response.Message = "Unauthorized access: User mismatch"
	// 	return GCFReturnStruct(Response)
	// }

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	datauser.UID = userInfo.Id

	id := r.URL.Query().Get("_id")
	if id == "" {
		Response.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	datatodo.ID = ID

	err = json.NewDecoder(r.Body).Decode(&datatodo)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	todo, _, err := UpdateTodo(mconn, collectionname, datatodo)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Update todo success"
	Response.Data = []model.Todo{todo}

	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// userInfo := watoken.DecodeGetId(os.Getenv(PASETOPUBLICKEY), token)

	// if userInfo != existingUser.UID {
	// 	Response.Message = "Unauthorized access: User mismatch"
	// 	return GCFReturnStruct(Response)
	// }

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	datatodo.User.UID = userInfo.Id

	id := r.URL.Query().Get("_id")
	if id == "" {
		Response.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	datatodo.ID = ID

	err = json.NewDecoder(r.Body).Decode(&datatodo)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	_, err = DeleteTodo(mconn, collectionname, datatodo.ID)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Delete todo success"

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
