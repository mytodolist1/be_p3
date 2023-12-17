package bep3

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mytodolist1/be_p3/model"
	"github.com/mytodolist1/be_p3/modul"
)

var mconn = modul.MongoConnect("MONGOSTRING", "mytodolist")

// user
func TestRegister(t *testing.T) {
	var data model.User
	data.Email = "nopal@gmail.com"
	data.Username = "nopal"
	// data.Role = "user"
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

// update with log
func TestUpdateUser(t *testing.T) {
	var data model.User
	data.Email = "dimas@gmail.com"
	data.Username = "dimas"

	id := "657437ffb905cf734635c9a8"
	ID, err := primitive.ObjectIDFromHex(id)
	data.ID = ID
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := modul.UpdateUser(mconn, "user", data)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updating user with id: %v", err)
			return
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
		fmt.Println(data)
	}
}

func TestChangePassword(t *testing.T) {
	var data model.User
	data.Password = "secret"
	data.ConfirmPassword = "secret"

	username := "dimass"
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
	id := "656877141a72c5656be2662d"
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
	anu, err := modul.GetUserFromUsername(mconn, "user", "qiqi1")
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

// admin
func TestGetUserFromRole(t *testing.T) {
	anu, err := modul.GetUserFromRole(mconn, "user", "user")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(anu)
}

// todo
func TestInsertTodo(t *testing.T) {
	var data model.Todo
	data.Title = "Pergi aja"
	data.Description = "Pergi ke pasar"
	data.Deadline = "12/18/2023"
	data.Time = "12:30 PM"
	data.Tags.Category = "personal"
	// data.IsDone = false

	uid := "c742a1aeebfa6cc8"

	_, err := modul.InsertTodo(mconn, "todo", data, uid)
	if err != nil {
		t.Errorf("Error inserting todo: %v", err)
	}
	fmt.Println(data)
}

func TestGetTodoFromCategory(t *testing.T) {
	anu, err := modul.GetTodoFromCategory(mconn, "todo", "personal")
	if err != nil {
		t.Errorf("Error getting todo: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestCheckCategory(t *testing.T) {
	anu, err := modul.CheckCategory(mconn, "category", "Personal")
	if err != nil {
		t.Errorf("Error getting todo: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestGetTodoFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("657eb9fb1e8a11c92dcf749f")
	anu, err := modul.GetTodoFromID(mconn, "todo", id)
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

// isDone
func TestTodoClear(t *testing.T) {
	var data model.TodoClear

	id := "657bd0fccd023400d63b49c4"

	ID, err := primitive.ObjectIDFromHex(id)
	data.Todo.ID = ID
	if err != nil {
		fmt.Printf("Data tidak berhasil di selesaikan")
	} else {

		status, err := modul.TodoClear(mconn, "todoclear", data)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error cleared todo with id: %v", err)
			return
		} else {
			fmt.Printf("Data berhasil di selesaikan untuk: %s\n", ID)
		}
		fmt.Println(data)
	}
}

func TestGetTodoDone(t *testing.T) {
	anu, err := modul.GetTodoDone(mconn, "todoclear")
	if err != nil {
		t.Errorf("Error getting todo: %v", err)
		return
	}
	fmt.Println(anu)
}

// with log
func TestUpdateTodo(t *testing.T) {
	var data model.Todo
	data.Title = "Belajar javascript"
	data.Description = "Hari ini belajar javascript"
	data.Deadline = "12/18/2023"
	data.Time = "10:00 PM"
	data.Tags.Category = "Personal"

	id := "657e6db23e913ed6f8dc4909"
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

// log
func TestGetLogTodoFromUID(t *testing.T) {
	uid := "c742a1aeebfa6cc8"

	anu, err := modul.GetLogTodoFromUID(mconn, "logtodo", uid)
	if err != nil {
		t.Errorf("Error getting log todo: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestGetLogUserFromUID(t *testing.T) {
	uid := "657437ffb905cf734635c9a8"

	anu, err := modul.GetLogUserFromUID(mconn, "loguser", uid)
	if err != nil {
		t.Errorf("Error getting log user: %v", err)
		return
	}
	fmt.Println(anu)
}

// func TestUpdateUser(t *testing.T) {
// 	var data model.User
// 	data.Email = "fulan12@gmail.com"
// 	data.Username = "fulan12"

// 	id := "6568592b8012346866b0ea1e"
// 	ID, err := primitive.ObjectIDFromHex(id)
// 	data.ID = ID
// 	if err != nil {
// 		fmt.Printf("Data tidak berhasil diubah")
// 	} else {

// 		_, status, err := modul.UpdateUser(mconn, "user", data)
// 		fmt.Println("Status", status)
// 		if err != nil {
// 			t.Errorf("Error updateting document: %v", err)
// 		} else {
// 			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
// 		}
// 	}
// }

// func TestGetAllUser(t *testing.T) {
// 	anu, err := modul.GetAllUser(mconn, "user")
// 	if err != nil {
// 		t.Errorf("Error getting user: %v", err)
// 		return
// 	}
// 	fmt.Println(anu)
// }

// func TestUpdateTodo(t *testing.T) {
// 	var data model.Todo
// 	data.Title = "Belajar Golang"
// 	data.Description = "Hari ini belajar golang"
// 	data.Deadline = "02/02/2021"

// 	id := "655c5047370b53741a9705d8"
// 	ID, err := primitive.ObjectIDFromHex(id)
// 	data.ID = ID
// 	if err != nil {
// 		fmt.Printf("Data tidak berhasil diubah")
// 	} else {

// 		_, status, err := modul.UpdateTodo(mconn, "todo", data)
// 		fmt.Println("Status", status)
// 		if err != nil {
// 			t.Errorf("Error updating todo with id: %v", err)
// 			return
// 		} else {
// 			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
// 		}
// 		fmt.Println(data)
// 	}
// }
