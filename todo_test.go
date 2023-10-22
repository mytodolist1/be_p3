package bep3

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGeneratePasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}

func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("Private Key: ", privateKey)
	fmt.Println("Public Key: ", publicKey)
	hasil, err := watoken.Encode("mytodolist", privateKey)
	fmt.Println("hasil: ", hasil, err)
}

// func TestValidateToken(t *testing.T) {
// 	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEwLTIyVDE1OjU0OjQyKzA3OjAwIiwiaWF0IjoiMjAyMy0xMC0yMlQxMzo1NDo0MiswNzowMCIsImlkIjoibXl0b2RvbGlzdCIsIm5iZiI6IjIwMjMtMTAtMjJUMTM6NTQ6NDIrMDc6MDAifZDQc54uatLh_mbG9PeBXjvxoCmLFrEnpj5Ach9ysg8-OP8SRoIVXKBxsLtmzsGEP_DJOXEqnW65j9Rtr8S8DAI" // Gantilah dengan token PASETO yang sesuai
// 	publicKey := "8459635e3ed946b66df39f5a30633f3e5e426a768cf0d69c3265cd7079c3f173"
// 	payload, _err := watoken.Decode(publicKey, tokenstring)
// 	if _err != nil {
// 		fmt.Println("expired token", _err)
// 	} else {
// 		fmt.Println("ID: ", payload.Id)
// 		fmt.Println("Di mulai: ", payload.Nbf)
// 		fmt.Println("Di buat: ", payload.Iat)
// 		fmt.Println("Expired: ", payload.Exp)
// 	}
// }

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")

	var userdata User
	userdata.Username = "budiman"
	userdata.Password = "secret"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)
}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var userdata User
	userdata.Username = "budiman"
	userdata.Password = "secret"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestInsertUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "mytodolist")
	var userdata User
	userdata.Username = "budiman"
	userdata.Role = "admin"
	userdata.Password = "secret"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

// func TestGCFPostHandler(t *testing.T) {

// 	// Membuat body request sebagai string
// 	requestBody := `{"username": "budiman", "password": "secret"}`

// 	// Membuat objek http.Request
// 	r := httptest.NewRequest("POST", "https://contoh.com/path", strings.NewReader(requestBody))
// 	r.Header.Set("Content-Type", "application/json")

// 	resp := GCFPostHandler("PASETOPRIVATEKEY", "MONGOSTRING", "mytodolist", "user", r)
// 	fmt.Println(resp)
// }
