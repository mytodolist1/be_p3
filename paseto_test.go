package bep3

import (
	"fmt"
	"testing"

	"github.com/whatsauth/watoken"
)

// paseto
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("Private Key: ", privateKey)
	fmt.Println("Public Key: ", publicKey)
	hasil, err := watoken.Encode("mytodolist", privateKey)
	fmt.Println("hasil: ", hasil, err)
}

// func TestGeneratePasswordHash(t *testing.T) {
// 	password := "secret"
// 	hash, _ := modul.HashPassword(password) // ignore error for the sake of simplicity

// 	fmt.Println("Password:", password)
// 	fmt.Println("Hash:    ", hash)

// 	match := modul.CheckPasswordHash(password, hash)
// 	fmt.Println("Match:   ", match)
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
