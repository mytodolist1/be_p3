package main

import (
	"fmt"
	"testing"

	modul "github.com/mytodolist1/be_p3/modul"
)

// login
func TestLoginAdmin(t *testing.T) {
	username := "admin"
	password := "admin"

	authenticated, err := modul.LoginAdmin(modul.MongoConn, "admin", username, password)
	if err != nil {
		t.Errorf("Error authenticating admin: %v", err)
	}

	if authenticated {
		fmt.Println("Admin authenticated successfully")
	} else {
		t.Errorf("Admin authentication failed")
	}
}

func TestInsertAdmin(t *testing.T) {
	username := "admin"
	password := "admin"

	insertedID, err := modul.InsertAdmin(modul.MongoConn, "admin", username, password)
	if err != nil {
		t.Errorf("Error inserting data: %v", err)
	}
	fmt.Printf("Data berhasil disimpan dengan id %s", insertedID.Hex())
}
