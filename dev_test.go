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
	var data model.User
	data.Email = "nopal@gmail.com"
	data.Username = "nopal"
	data.Role = "user"
	data.Password = "secret"

	err := modul.Register(mconn, "user", data)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	} else {
		fmt.Println("Register success", data)
	}
}

func TestLogIn(t *testing.T) {
	var data model.User
	data.Username = "nopal1"
	data.Password = "secret111"

	user, status, err := modul.LogIn(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error logging in user: %v", err)
	} else {
		fmt.Println("Login success", user)
	}
}

func TestUpdateUser(t *testing.T) {
	var data model.User
	data.Email = "nopal@gmail.com"
	data.Username = "nopal1"
	data.Role = "user"

	data.Password = "secret"

	id, err := primitive.ObjectIDFromHex("654d9264ffbef4e358812043")
	data.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := modul.UpdateUser(mconn, "user", data)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updateting document: %v", err)
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
	}
}

func TestChangePassword(t *testing.T) {
	var data model.User
	data.Email = "nopal@gmail.com"
	data.Username = "nopal1"
	data.Role = "user"

	data.Password = "secret111"

	_, status, err := modul.ChangePassword(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error updateting document: %v", err)
	} else {
		fmt.Println("Password berhasil diubah dengan username:", data.Username)
	}
}

func TestDeleteUser(t *testing.T) {
	username := "tejo1"

	err := modul.DeleteUser(mconn, "user", username)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	} else {
		fmt.Println("Delete user success")
	}
}

func TestGetUserFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("6549e6252174254280d650af")
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
