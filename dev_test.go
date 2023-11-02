package bep3

import (
	"fmt"
	"testing"

	model "github.com/mytodolist1/be_p3/model"
	modul "github.com/mytodolist1/be_p3/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// validasi
func TestRegister(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var userdata model.User
	userdata.Email = "tejoko@gmail.com"
	userdata.Username = "tejoko"
	userdata.Role = "user"
	userdata.Password = "secret"

	err := modul.Register(mconn, "user", userdata)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	} else {
		fmt.Println("Register success")
	}
}

func TestLogIn(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var userdata model.User
	userdata.Username = "tejoko"
	userdata.Password = "secret"
	user, status, err := modul.LogIn(mconn, "user", userdata)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error logging in user: %v", err)
	} else {
		fmt.Println("Login success", user)
	}
}

// user
// func TestInsertUser(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "mytodolist")
// 	var userdata User
// 	userdata.Email = "budiman@gmail.com"
// 	userdata.Username = "budiman"
// 	userdata.Role = "admin"
// 	userdata.Password = "secret"

// 	nama := InsertUser(mconn, "user", userdata)
// 	fmt.Println(nama)
// }

func TestGetAllUserFromUsername(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	anu, _ := modul.GetUserFromUsername(mconn, "user", "budiman")
	fmt.Println(anu)
}

func TestGetAllUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	anu := modul.GetAllUser(mconn, "user")
	fmt.Println(anu)
}

// todo
func TestInsertTodo(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var tododata model.Todo
	tododata.Title = "Perjalanan"
	tododata.Description = "pergi ke bali"
	tododata.IsDone = true

	nama, err := modul.InsertTodo(mconn, "todo", tododata)
	if err != nil {
		t.Errorf("Error inserting todo: %v", err)
	}
	fmt.Println(nama)
}

func TestGetTodoFromID(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	id, _ := primitive.ObjectIDFromHex("653e02ab28597c2c37171d44")
	anu := modul.GetTodoFromID(mconn, "todo", id)
	fmt.Println(anu)
}
