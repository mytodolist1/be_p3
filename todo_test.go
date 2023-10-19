package bep3

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// paseto

func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("Private Key: ", privateKey)
	fmt.Println("Public Key: ", publicKey)
	hasil, err := watoken.Encode("mytodolist", privateKey)
	fmt.Println("hasil: ", hasil, err)
}

func TestValidateToken(t *testing.T) {
	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEwLTE5VDE0OjE0OjQ0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMC0xOVQxMjoxNDo0NCswNzowMCIsImlkIjoibXl0b2RvbGlzdCIsIm5iZiI6IjIwMjMtMTAtMTlUMTI6MTQ6NDQrMDc6MDAifUpIr_FRgF_teFsWe1zvDUP5jgjYfR_MLph9CwuElISzwjr0LI546Sw5v7FV7_8eAtSNw5hiypWkU6woUlth3gs" // Gantilah dengan token PASETO yang sesuai
	publicKey := "72654165d09b9f0a8b4f0c5815775ed5fc933069ce2e006b4e62a65bea6f06e3"
	payload, _err := watoken.Decode(publicKey, tokenstring)
	if _err != nil {
		fmt.Println("expired token", _err)
	} else {
		fmt.Println("ID: ", payload.Id)
		fmt.Println("Di mulai: ", payload.Nbf)
		fmt.Println("Di buat: ", payload.Iat)
		fmt.Println("Expired: ", payload.Exp)
	}

}

// hash password
func TestGenerateHashPassword(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckHashPassword(password, hash)
	fmt.Println("Match:   ", match)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var userdata User
	userdata.Username = "budi"
	userdata.Password = "secret"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckHashPassword(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var userdata User
	userdata.Username = "budi"
	userdata.Password = "secret"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

// user
func TestInsertUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var userdata User
	userdata.Username = "budiman"
	userdata.Role = "admin"
	userdata.Password = "secret"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

// login
// func TestLoginUser(t *testing.T) {
// 	username := "user"
// 	password := "user"

// 	authenticated, err := LoginUser(MongoConn, "user", username, password)
// 	if err != nil {
// 		t.Errorf("Error authenticating user: %v", err)
// 	}

// 	if authenticated {
// 		fmt.Println("User authenticated successfully")
// 	} else {
// 		t.Errorf("User authentication failed")
// 	}
// }

// func TestInsertUser(t *testing.T) {
// 	username := "user"
// 	password := "user"

// 	insertedID, err := InsertUser(MongoConn, "user", username, password)
// 	if err != nil {
// 		t.Errorf("Error inserting data: %v", err)
// 	}
// 	fmt.Printf("Data berhasil disimpan dengan id %s", insertedID.Hex())
// }
