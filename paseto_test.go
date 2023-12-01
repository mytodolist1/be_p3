package bep3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whatsauth/watoken"
)

// paseto
// func TestGeneratePrivateKeyPaseto(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	fmt.Println("Private Key: ", privateKey)
// 	fmt.Println("Public Key: ", publicKey)

// 	uid := "81381f10-cd45-42e4-a72c-642f34bdd53d"
// 	hasil, err := watoken.Encode(uid, privateKey)
// 	fmt.Println("hasil: ", hasil, err)
// }

func TestGenerateTokenPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("privateKey : ", privateKey)
	fmt.Println("publicKey : ", publicKey)
	userid := "81381f10-cd45-42e4-a72c-642f34bdd53d"

	tokenstring, err := watoken.Encode(userid, privateKey)
	require.NoError(t, err)
	body, err := watoken.Decode(publicKey, tokenstring)
	fmt.Println("signed : ", tokenstring)
	fmt.Println("isi : ", body)
	require.NoError(t, err)
}

func TestDecoedTokenPaseto(t *testing.T) {
	publicKey := "3fca58bcee37564ae23005b9aefe51b93cda7327a0831f533cae57f26ae70398"
	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEyLTAxVDA2OjMzOjMyWiIsImlhdCI6IjIwMjMtMTItMDFUMDQ6MzM6MzJaIiwiaWQiOiJxaXFpIiwibmJmIjoiMjAyMy0xMi0wMVQwNDozMzozMloifWNdG8-O7zBRsXlT78B8T5QEH-UlvYqUWBgSa22gAIs2noox_o5QZ-gj4if8gOYurkLa2oU7T7wHWBNwOxI0sAU"

	uid, err := watoken.Decode(publicKey, tokenstring)
	require.NoError(t, err)
	fmt.Println("uid : ", uid)
}

func TestDecodeToken(t *testing.T) {
	//privateKey := "461ce0e87748fd656c518b870da217dc200fc8d3b6275dda8cf14943424bf8c49e2ece1954df1ea8b151fba59cc7cbd4fb810b69716149e1c26169227bd5b6868ac78b29e58b97d4018d66ad9aed4c608028f8e188dd976fa5f61fb46b47c37365d8d07b2b8d915ec9771904b608e6ba1a91b815f9e8aece8255a660b528287e"
	publicKey := "3fca58bcee37564ae23005b9aefe51b93cda7327a0831f533cae57f26ae70398"
	//userid := "awangga"
	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEyLTAxVDA2OjMzOjMyWiIsImlhdCI6IjIwMjMtMTItMDFUMDQ6MzM6MzJaIiwiaWQiOiJxaXFpIiwibmJmIjoiMjAyMy0xMi0wMVQwNDozMzozMloifWNdG8-O7zBRsXlT78B8T5QEH-UlvYqUWBgSa22gAIs2noox_o5QZ-gj4if8gOYurkLa2oU7T7wHWBNwOxI0sAU"
	idstring := watoken.DecodeGetId(publicKey, tokenstring)
	if idstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println("TestWaTokenDecodewithStaticKey idstring : ", idstring)
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
