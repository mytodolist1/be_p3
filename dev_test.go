package bep3

import (
	"fmt"
	"testing"

	model "github.com/mytodolist1/be_p3/model"
	modul "github.com/mytodolist1/be_p3/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mconn = modul.MongoConnect("MONGOSTRING", "mytodolist")

// user
func TestRegister(t *testing.T) {
	var data model.User
	data.Email = "fulan1@gmail.com"
	data.Username = "fulan1"
	data.Role = "user"
	data.Password = "secret"
	data.ConfirmPassword = "secret"

	err := modul.Register(mconn, "user", data)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	} else {
		fmt.Println("Register success", data.Username)
	}
}

func TestLogIn(t *testing.T) {
	var data model.User
	data.Username = "fulan1"
	data.Password = "secret"
	data.Role = "user"

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
	data.Email = "fulan12@gmail.com"
	data.Username = "fulan12"

	id := "6568592b8012346866b0ea1e"
	ID, err := primitive.ObjectIDFromHex(id)
	data.ID = ID
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
	data.Password = "secret"

	username := "tejoko"
	data.Username = username

	_, status, err := modul.ChangePassword(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error updateting document: %v", err)
	} else {
		fmt.Println("Password berhasil diubah dengan username:", username)
	}
}

func TestDeleteUser(t *testing.T) {
	var data model.User
	data.Username = "admin"

	status, err := modul.DeleteUser(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error deleting document: %v", err)
	} else {
		fmt.Println("Delete user" + data.Username + "success")
	}
}

func TestGetUserFromID(t *testing.T) {
	id := "6550a8f0b5b4a0a7f89941aa"
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to ObjectID: %v", err)
		return
	}

	anu, err := modul.GetUserFromID(mconn, "user", ID)
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
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
	anu, err := modul.GetAllUser(mconn, "user")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(anu)
}

// todo
func TestInsertTodo(t *testing.T) {
	var data model.Todo
	data.Title = "Pergi ke sana"
	data.Description = "membeli itu ini"
	data.Deadline = "02/02/2021"
	// data.IsDone = false
	data.User.Username = "nopal"

	id, err := modul.InsertTodo(mconn, "todo", data)
	if err != nil {
		t.Errorf("Error inserting todo: %v", err)
	}
	fmt.Println(id)
}

func TestGetTodoFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("655c4408d06d3d2ddba5d1d7")
	anu, err := modul.GetTodoFromID(mconn, "todo", id)
	if err != nil {
		t.Errorf("Error getting todo: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestGetTodoFromUsername(t *testing.T) {
	anu, err := modul.GetTodoFromUsername(mconn, "todo", "nopal")
	if err != nil {
		t.Errorf("Error getting todo: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestGetTodoList(t *testing.T) {
	anu, err := modul.GetTodoList(mconn, "todo")
	if err != nil {
		t.Errorf("Error getting todo: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestUpdateTodo(t *testing.T) {
	var data model.Todo
	data.Title = "Belajar Golang"
	data.Description = "Hari ini belajar golang"
	data.Deadline = "02/02/2021"

	id := "655c5047370b53741a9705d8"
	ID, err := primitive.ObjectIDFromHex(id)
	data.ID = ID
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := modul.UpdateTodo(mconn, "todo", data)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updating todo with id: %v", err)
			return
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
		fmt.Println(data)
	}
}

func TestDeleteTodo(t *testing.T) {
	id := "655c4408d06d3d2ddba5d1d7"
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to ObjectID: %v", err)
		return
	} else {

		status, err := modul.DeleteTodo(mconn, "todo", ID)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error deleting document: %v", err)
			return
		} else {
			fmt.Println("Delete success")
		}
	}
}

// // paseto
// func TestGeneratePasswordHash(t *testing.T) {
// 	password := "secret"
// 	hash, _ := modul.HashPassword(password) // ignore error for the sake of simplicity

// 	fmt.Println("Password:", password)
// 	fmt.Println("Hash:    ", hash)

// 	match := modul.CheckPasswordHash(password, hash)
// 	fmt.Println("Match:   ", match)
// }

// func TestGeneratePrivateKeyPaseto(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	fmt.Println("Private Key: ", privateKey)
// 	fmt.Println("Public Key: ", publicKey)
// 	hasil, err := watoken.Encode("mytodolist", privateKey)
// 	fmt.Println("hasil: ", hasil, err)
// }

// func TestHashFunction(t *testing.T) {
// 	// mconn := SetConnection("MONGOSTRING", "mytodolist")

// 	var userdata model.User
// 	userdata.Username = "budiman"
// 	userdata.Password = "secret"

// 	filter := bson.M{"username": userdata.Username}
// 	res := atdb.GetOneDoc[model.User](mconn, "user", filter)
// 	fmt.Println("Mongo User Result: ", res)
// 	hash, _ := modul.HashPassword(userdata.Password)
// 	fmt.Println("Hash Password : ", hash)
// 	match := modul.CheckPasswordHash(userdata.Password, res.Password)
// 	fmt.Println("Match:   ", match)
// }

// func TestIsPasswordValid(t *testing.T) {
// 	// mconn := SetConnection("MONGOSTRING", "mytodolist")
// 	var userdata model.User
// 	userdata.Username = "budiman"
// 	userdata.Password = "secret"

// 	anu := modul.IsPasswordValid(mconn, "user", userdata)
// 	fmt.Println(anu)
// }
