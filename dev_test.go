package bep3

import (
	"fmt"
	"testing"

	model "github.com/mytodolist1/be_p3/model"
	modul "github.com/mytodolist1/be_p3/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mconn = SetConnection("MONGOSTRING", "mytodolist")

// user
func TestRegister(t *testing.T) {
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

func TestChangePassword(t *testing.T) {
	username := "tejoko"
	oldpassword := "secretbaruhehe"
	newpassword := "secret"

	var userdata model.User
	userdata.Username = username
	userdata.Password = newpassword

	userdata, status, err := modul.ChangePassword(mconn, "user", username, oldpassword, newpassword)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error changing password: %v", err)
	} else {
		fmt.Println("Password change success for user", userdata)
	}
}

func TestDeleteUser(t *testing.T) {
	username := "tejo_ko"

	err := modul.DeleteUser(mconn, "user", username)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	} else {
		fmt.Println("Delete user success")
	}

	_, err = modul.GetUserFromUsername(mconn, "user", username)
	if err == nil {
		fmt.Println("Data masih ada")
	}
}

func TestGetUserFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("653e03317043a1bb2588")
	anu, _ := modul.GetUserFromID(mconn, "user", id)
	fmt.Println(anu)
}

func TestGetUserFromUsername(t *testing.T) {
	anu, err := modul.GetUserFromUsername(mconn, "user", "budiman")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(anu)
}


func TestGetUserFromEmail(t *testing.T) {
	anu, _ := modul.GetUserFromEmail(mconn, "user", "tejo@gmail.com")
	fmt.Println(anu)
}

func TestGetAllUser(t *testing.T) {
	anu := modul.GetAllUser(mconn, "user")
	fmt.Println(anu)
}

// todo
func TestInsertTodo(t *testing.T) {
	var tododata model.Todo
	tododata.Title = "Belajar Golang"
	tododata.Description = "Hari ini belajar testing"
	tododata.IsDone = true

	nama, err := modul.InsertTodo(mconn, "todo", tododata)
	if err != nil {
		t.Errorf("Error inserting todo: %v", err)
	}
	fmt.Println(nama)
}

func TestGetTodoFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("653e02ab28597c2c37171d44")
	anu := modul.GetTodoFromID(mconn, "todo", id)
	fmt.Println(anu)
}

func TestGetTodoList(t *testing.T) {
	anu := modul.GetTodoList(mconn, "todo")
	fmt.Println(anu)
}
