package modul

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	model "github.com/mytodolist1/be_p3/model"
	"github.com/whatsauth/watoken"
)

// user
func GCFHandlerGetAllUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	userlist, err := GetAllUser(mconn, collectionname)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Get User Success"
	Response.Data = userlist

	return GCFReturnStruct(Response)
}

func GCFHandlerGetUserByUsername(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	user, err := GetUserFromUsername(mconn, collectionname, datauser.Username)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Hello " + " " + user.Username

	return GCFReturnStruct(Response)
}

func GCFHandlerRegister(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	err = Register(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Register success"

	return GCFReturnStruct(Response)
}

func GCFHandlerLogIn(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	user, _, err := LogIn(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token :" + err.Error()
	} else {
		Response.Message = "Selamat Datang" + " " + user.Username
		Response.Token = tokenstring
	}

	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	user, status, err := UpdateUser(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Update success " + " " + user.Username + " " + strconv.FormatBool(status)

	return GCFReturnStruct(Response)
}

func GCFHandlerChangePassword(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	user, status, err := ChangePassword(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Password change success for user" + " " + user.Username + " " + strconv.FormatBool(status)

	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteUser(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	}
	status, err := DeleteUser(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Delete user success" + " " + datauser.Username + " " + strconv.FormatBool(status)

	return GCFReturnStruct(Response)
}

// todo
func GCFHandlerGetAllTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.TodoResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.Todo

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	_, err = GetTodoList(mconn, collectionname)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Get todo success"
	Response.Data = datauser
	return GCFReturnStruct(Response)
}

func GCFHandlerGetTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.TodoResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.Todo

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	_, err = GetTodoFromID(mconn, collectionname, datauser.ID)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Get todo success"
	return GCFReturnStruct(Response)
}

func GCFHandlerInsertTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.TodoResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.Todo

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	_, err = InsertTodo(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Insert todo success"
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.TodoResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.Todo

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	_, status, err := UpdateTodo(mconn, collectionname, datauser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Update todo success" + " " + strconv.FormatBool(status)
	Response.Data = datauser
	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteTodo(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.TodoResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var datauser model.Todo

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	status, err := DeleteTodo(mconn, collectionname, datauser.ID)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Delete todo success" + " " + strconv.FormatBool(status)
	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
