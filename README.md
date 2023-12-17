# Backend Mytodolist

Repositori ini adalah Backend dari aplikasi Mytodolist yang ditulis dengan bahasa Go.

## model
- `type.go`
    1. Import Standard Library
        - `go.mongodb.org/mongo-driver/bson/primitive` adalah bagian dari MongoDB Go Driver yang menyediakan tipe data dasar dan fungsi konversi untuk mengoperasikan BSON (Binary JSON) dalam Go.
            Contoh penggunaannya:
            ```go
            ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
            ```

## modul
- `controller.go`
    1. Import Standard Library
        - `context` adalah package yang menyediakan tool untuk mentransmisikan informasi antar goroutine.
            Contoh penggunaannya:
            ```go
            result, err := cols.InsertOne(context.Background(), docs)
            ```
        - `errors` adalah package yang digunakan untuk menangani kesalahan atau error.
            Contoh penggunaannya:
            ```go
            if errors.Is(err, mongo.ErrNoDocuments) {
                fmt.Println("no data found for ID", _id)
            }
            ```
        - `fmt` adalah package yang digunakan untuk mencetak teks.
            Contoh penggunaannya:
            ```go
            fmt.Printf("InsertOneDoc: %v\n", err)
            ```
        - `os` adalah package yang digunakan untuk berinteraksi dengan operation sistem, seperti mengakses environment variables, dll.
            Contoh penggunaannya:
            ```go
            DBString: os.Getenv(MONGOCONNSTRINGENV),
            ``` 
        - `strings` adalah package yang digunakan untuk memanipulasi dan bekerja dengan string.
            Contoh penggunaannya:
            ```go
            if strings.Contains(userdata.Username, " ") {
                return fmt.Errorf("Username tidak boleh mengandung spasi")
            }
            ```
        - `time` adalah package dari Go yang digunakan untuk memanipulasi dan bekerja dengan waktu.
            Contoh penggunaannya:
            ```go
            {Key: "timestamp.updatedat", Value: time.Now()},
            ```
        - `crypto/rand` adalah package yang digunakan untuk menghasilkan bilangan acak yang aman dari segi kriptografi.
            Contoh penggunaannya:
            ```go
            _, err := rand.Read(bytes)
            ```
        - `encoding/hex` adalah package yang digunakan untuk mengubah data biner menjadi representasi heksadesimal (hex) dan sebaliknya.
            Contoh penggunaannya:
            ```go
            hex.EncodeToString(bytes), nil
            ```
        - `go.mongodb.org/mongo-driver/bson` adalah package yang digunakan untuk melakukan marshaling (konversi objek ke BSON) dan unmarshaling (konversi BSON ke objek).
            Contoh penggunaannya:
            ```go
            filter := bson.M{"_id": userdata.ID}
            ```
        - `go.mongodb.org/mongo-driver/bson/primitive` adalah package yang menyediakan tipe data dasar dan fungsi konversi untuk mengoperasikan BSON.
            Contoh penggunaannya:
            ```go
            func GetUserFromID(db *mongo.Database, col string, _id primitive.ObjectID){}
            ```
        - `go.mongodb.org/mongo-driver/mongo` adalah package yang menyediakan fungsionalitas untuk berinteraksi dengan MongoDB.
            Contoh penggunaannya:
            ```go
            func MongoConnect(MONGOCONNSTRINGENV, dbname string) *mongo.Database{}
            ```
        - `go.mongodb.org/mongo-driver/mongo/options` adalah package yang menyediakan berbagai opsi yang dapat dikonfigurasi saat berinteraksi dengan server MongoDB.
            Contoh penggunaannya:
            ```go
            options := options.Update().SetUpsert(true)
            ```

    2. Import External Library
        - `github.com/aiteung/atdb` adalah modul eksternal yang di import untuk menggunakan structnya.
            Contoh penggunaannya:
            ```go
            var DBmongoinfo = atdb.DBInfo{}
            ```
        - `github.com/badoux/checkmail` adalah modul eksternal yang di import untuk validasi email.
            Contoh penggunaannya:
            ```go
            err := checkmail.ValidateFormat(userdata.Email)
            ```

    3. Import this Module Repository
        - `github.com/mytodolist1/be_p3/model` adalah modul yang dibuat pada repositori ini dan di import karena berbeda folder untuk mengambil structnya.
            Contoh penggunaannya:
            ```go
            func Register(db *mongo.Database, col string, userdata model.User) error {}
            ```

- `handler.go`
    1. Import Standard Library
        - `encoding/json` adalah package yang menyediakan fungsi-fungsi untuk mengkodekan (marshal) dan mendekodekan (unmarshal) data JSON.
            Contoh penggunaannya:
            ```go
            err = json.NewDecoder(r.Body).Decode(&datatodo)
            ```
        - `net/http` adalah package yang menyediakan dukungan untuk membangun layanan web (HTTP) dan mengirim permintaan HTTP.
            Contoh penggunaannya:
            ```go
            func GCFHandlerGetUserFromToken(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {}
            ```
        - `os` adalah package yang digunakan untuk berinteraksi dengan operation sistem, seperti mengakses environment variables, dll.
            Contoh penggunaannya:
            ```go
            tokenstring, err := paseto.Encode(user.UID, user.Role, os.Getenv(PASETOPRIVATEKEYENV))
            ```
        - `go.mongodb.org/mongo-driver/bson/primitive` adalah package yang menyediakan tipe data dasar dan fungsi konversi untuk mengoperasikan BSON.
            Contoh penggunaannya:
            ```go
            ID, err := primitive.ObjectIDFromHex(id)
            ```

    2. Import this Module Repository
        - `github.com/mytodolist1/be_p3/model` adalah modul yang digunakan untuk mengambil structnya.
            Contoh penggunaannya:
            ```go
            var datauser model.User
            ```
        - `github.com/mytodolist1/be_p3/paseto` adalah modul yang digunakan untuk encode dan decode token.
            Contoh penggunaannya:
            ```go
            userInfo, err := paseto.Decode(os.Getenv(PASETOPUBLICKEY), token)
            ```

- `hash.go`
    1. Import External Library
        - `golang.org/x/crypto/bcrypt` adalah package yang digunakan untuk mengenkripsi dan memverifikasi kata sandi menggunakan fungsi bcrypt.
            Contoh penggunaannya:
            ```go
            bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
            ```

## paseto
- `paseto.go`
    1. Import Standard Library
        - `encoding/json` adalah package yang menyediakan fungsi-fungsi untuk mengkodekan (marshal) dan mendekodekan (unmarshal) data JSON.
            Contoh penggunaannya:
            ```go
            json.Unmarshal(token.ClaimsJSON(), &payload)
            ```
        - `fmt` adalah package yang digunakan untuk mencetak teks.
            Contoh penggunaannya:
            ```go
            fmt.Println("Decode ParseV4Public : ", err)
            ```
        - `time` adalah package dari Go yang digunakan untuk memanipulasi dan bekerja dengan waktu.
            Contoh penggunaannya:
            ```go
            token.SetIssuedAt(time.Now())
            ```
    
    2. Import External Library
        - `aidanwoods.dev/go-paseto` adalah package yang digunakan untuk membuat token paseto.
            Contoh penggunaannya:
            ```go
            secretKey := paseto.NewV4AsymmetricSecretKey()
            ```

- `paseto_test.go`
    1. Import Standard Library
        - `fmt` adalah package yang digunakan untuk mencetak teks.
            Contoh penggunaannya:
            ```go
            fmt.Println("hasil: ", hasil, err)
            ```
        - `testing` adalah package yang menyediakan alat dan infrastruktur untuk menulis dan menjalankan pengujian (testing).
            Contoh penggunaannya:
            ```go
            func TestGenerateTokenPaseto(t *testing.T) {}
            ```
    
    2. Import External Library
        - `github.com/stretchr/testify/require` adalah package yang digunakan untuk menghentikan eksekusi pengujian jika gagal.
            Contoh penggunaannya:
            ```go
            require.NoError(t, err)
            ```

    3. Import this Module Repository
        - `github.com/mytodolist1/be_p3/paseto` adalah modul yang digunakan untuk encode dan decode token.
            Contoh penggunaannya:
            ```go
            privateKey, publicKey := paseto.GenerateKey()
            ```

#### `dev_test.go`
-   1. Import Standard Library
        - `fmt` adalah package yang digunakan untuk mencetak teks.
            Contoh penggunaannya:
            ```go
            fmt.Println("Status", status)
            ```
        - `testing` adalah package yang menyediakan alat dan infrastruktur untuk menulis dan menjalankan pengujian (testing).
            Contoh penggunaannya:
            ```go
            func TestLogIn(t *testing.T) {}
            ```
        - `go.mongodb.org/mongo-driver/bson/primitive` adalah package yang menyediakan tipe data dasar dan fungsi konversi untuk mengoperasikan BSON.
            Contoh penggunaannya:
            ```go
            ID, err := primitive.ObjectIDFromHex(id)
            ```

-   2. Import this Module Repository
        - `github.com/mytodolist1/be_p3/model` adalah modul yang digunakan untuk mengambil structnya.
            Contoh penggunaannya:
            ```go
            var data model.User
            ```
        - `github.com/mytodolist1/be_p3/modul` adalah modul yang digunakan untuk memanggil controller.
            Contoh penggunaannya:
            ```go
            var mconn = modul.MongoConnect("MONGOSTRING", "mytodolist")
            ```

## tags
- Release Version Go
    ```bash
    git tag v0.0.1
    git push origin --tags
    go list -m github.com/mytodolist1/be_p3@v0.0.1
    ```