package modul

import (
	"encoding/json"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mytodolist1/be_p3/model"
	"github.com/mytodolist1/be_p3/paseto"
)

var (
	Responsed model.Credential
	Response  model.TodoResponse
	resp      model.LogTodoResponse
	datauser  model.User
	datatodo  model.Todo
	logtodo   model.LogTodo
	isdone    model.TodoClear
)

// for user
// user
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

func GCFHandlerGetUserByID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

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

	user, err := GetUserFromID(mconn, collectionname, ID)
	if err != nil {
		Responsed.Message = "Error retrieving user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Hello user " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerGetUserFromToken(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	user, err := GetUserFromToken(mconn, collectionname, userInfo.Id)
	if err != nil {
		Responsed.Message = err.Error()
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
	tokenstring, err := paseto.Encode(user.UID, user.Role, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Responsed.Message = "Gagal Encode Token :" + err.Error()

	} else {
		Responsed.Message = "Selamat Datang " + user.Role + " " + user.Username
		Responsed.Token = tokenstring
		Responsed.Data = []model.User{user}
	}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerUpdateUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
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
	Responsed.Message = "Update success " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerChangePassword(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
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
	Responsed.Message = "Password change success for user " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerDeleteUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
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
func GCFHandlerGetTodoListByUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	// err = json.NewDecoder(r.Body).Decode(&datatodo)
	// if err != nil {
	// 	Response.Message = "error parsing application/json3: " + err.Error()
	// 	return GCFReturnStruct(Response)
	// }

	todo, err := GetTodoFromToken(mconn, collectionname, userInfo.Id)
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
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

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

	// err = json.NewDecoder(r.Body).Decode(&datatodo)
	// if err != nil {
	// 	Response.Message = "error parsing application/json3: " + err.Error()
	// 	return GCFReturnStruct(Response)
	// }

	todo, err := GetTodoFromID(mconn, collectionname, ID)
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
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&datatodo)
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
	Response.Message = "Insert todo success for " + datatodo.Title
	Response.Data = []model.Todo{datatodo}

	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

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
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

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

	// err = json.NewDecoder(r.Body).Decode(&datatodo)
	// if err != nil {
	// 	Response.Message = "error parsing application/json3: " + err.Error()
	// 	return GCFReturnStruct(Response)
	// }

	_, err = DeleteTodo(mconn, collectionname, ID)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Delete todo success"

	return GCFReturnStruct(Response)
}

// for admin
// user
func GCFHandlerGetAllUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	if userInfo.Role != "admin" {
		Responsed.Message = "You are not admin"
		return GCFReturnStruct(Responsed)

	} else {
		userlist, err := GetUserFromRole(mconn, collectionname, "user")
		if err != nil {
			Responsed.Message = err.Error()
			return GCFReturnStruct(Responsed)
		}

		Responsed.Status = true
		Responsed.Message = "Get User Success"
		Responsed.Data = userlist

		return GCFReturnStruct(Responsed)
	}
}

// todo
func GCFHandlerGetAllTodoList(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	if userInfo.Role != "admin" {
		Responsed.Message = "You are not admin"
		return GCFReturnStruct(Responsed)

	} else {
		todolist, err := GetTodoList(mconn, collectionname)
		if err != nil {
			Response.Message = err.Error()
			return GCFReturnStruct(Response)
		}

		Response.Status = true
		Response.Message = "Get todo success"
		Response.Data = todolist

		return GCFReturnStruct(Response)
	}
}

// log todo
func GCFHandlerGetLogTodoList(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)

	token := r.Header.Get("Authorization")
	if token == "" {
		resp.Message = "error parsing application/json1:"
		return GCFReturnStruct(resp)
	}

	userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		resp.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(resp)
	}

	if userInfo.Role != "admin" {
		Responsed.Message = "You are not admin"
		return GCFReturnStruct(Responsed)

	} else {
		loglist, err := GetLogTodoList(mconn, collectionname)
		if err != nil {
			resp.Message = err.Error()
			return GCFReturnStruct(resp)
		}

		resp.Status = true
		resp.Message = "Get log todo success"
		resp.Data = loglist

		return GCFReturnStruct(Response)
	}
}

// log for user
func GCFHandlerGetLogTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)

	token := r.Header.Get("Authorization")
	if token == "" {
		resp.Message = "error parsing application/json1:"
		return GCFReturnStruct(resp)
	}

	userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		resp.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(resp)
	}

	loglist, err := GetLogTodoFromUID(mconn, collectionname, userInfo.Id)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get log todo success"
	resp.Data = []model.LogTodo{loglist}

	return GCFReturnStruct(Response)
}

// isDone
func GCFHandlerIsDone(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

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

	ID = datatodo.ID

	_, err = TodoClear(mconn, collectionname, ID, isdone)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "IsDone success"

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// user
// func GCFHandlerGetAllUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	Responsed.Status = false

// 	userlist, err := GetAllUser(mconn, collectionname)
// 	if err != nil {
// 		Responsed.Message = err.Error()
// 		return GCFReturnStruct(Responsed)
// 	}

// 	Responsed.Status = true
// 	Responsed.Message = "Get User Success"
// 	Responsed.Data = userlist

// 	return GCFReturnStruct(Responsed)
// }
