package bep3

import (
	"fmt"
	"net/http/httptest"
	"strings"
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

func TestValidateToken(t *testing.T) {
	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEwLTIyVDE1OjA4OjQwKzA3OjAwIiwiaWF0IjoiMjAyMy0xMC0yMlQxMzowODo0MCswNzowMCIsImlkIjoibXl0b2RvbGlzdCIsIm5iZiI6IjIwMjMtMTAtMjJUMTM6MDg6NDArMDc6MDAifbrXpe86NxnAX1ULNBeivaM53Mgrc1waKX2XGm1bt3JHiksBy7hK6d4TijzPNQcWuhLzeyNl7EwumfmIS7bj5gU" // Gantilah dengan token PASETO yang sesuai
	publicKey := "915ae398cfbc8902f33077c449bdbd5c9f475667fe79c2356e9e800798bb9839"
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

func TestGCFPostHandler(t *testing.T) {

	// Membuat body request sebagai string
	requestBody := `{"username": "dani", "password": "secret"}`

	// Membuat objek http.Request
	r := httptest.NewRequest("POST", "https://contoh.com/path", strings.NewReader(requestBody))
	r.Header.Set("Content-Type", "application/json")

	resp := GCFPostHandler("PASETOPRIVATEKEY", "MONGOSTRING", "trensentimen", "user", r)
	fmt.Println(resp)
}
