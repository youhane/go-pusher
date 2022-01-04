package main

import (
	"database/sql"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pusher/pusher-http-go"
)

// blueprintnya
// instantiate client appnya
// remember to replace the values inside here
// var client = pusher.Client{
// 	AppID:   "PUSHER_APP_ID",
// 	Key:     "PUSHER_APP_KEY",
// 	Secret:  "PUSHER_APP_SECRET",
// 	Cluster: "PUSHER_APP_CLUSTER",
// 	Secure:  true,
// }

// instantiate client appnya
// remember to replace the values inside here
var client = pusher.Client{
	AppID:   "1326078",
	Key:     "3fe7405204654449a6d1",
	Secret:  "f0e96eae0a5534a0219b",
	Cluster: "ap1",
	Secure:  true,
}

// bikin struct yang bakal dipake terus
type Photo struct {
	// ID, yang tipenya int64, yang disebutnya id
	ID int64 `json:"id"`
	// Src, yang tipenya string yang disebutnya src
	Src string `json:"src"`
}

// bikin route handler, yang ngereturn apa yang diminta dari endpoint yang udah dibikin di main
// tipe data disamping () parameter, itu return type, kalo misal mau dibatesin
func getPhotos(db *sql.DB) echo.HandlerFunc {
	// ini function yang ngereturn function lagi
	return func(c echo.Context) error {
		// run query, dari db yang udah dibuat
		rows, err := db.Query("SELECT * FROM photos")

		// kalo ada error apapun, panic
		if err != nil {
			panic(err)
		}

		// dia bakal jalan, kalo 1 scope code yang dia berada udah selesai dijalanin semua dulu,
		// baru dia bakal jalan
		defer rows.Close()

		// bikin result yang bentuknya PhotoCollection
		result := PhotoCollection{}

		// for rows.Next(), ngeloop sampe semua record yang ada udah dilewatin
		for rows.Next() {
			// bikin Photo
			photo := Photo{}

			// cek ada error apa ngga, die ngecek alamatnya pake pointer ada apa ngga
			err2 := rows.Scan(&photo.ID, &photo.Src)

			// dan kalo ada, maka itu error
			if err2 != nil {
				panic(err2)
			}

			// dari resultsnya, kalo ga error, append datanya ke photo yang ada disini
			result.Photos = append(result.Photos, photo)
		}
		// return si resultnya dalam bentuk JSON, sama HTTP Code Status OK
		return c.JSON(http.StatusOK, result)
	}
}

func uploadPhoto(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ambil filenya dari context,
		// dan return resultnya berupa object multipart, informasi header filenya, dan error kalo ada
		// si "file", itu nama input tagnya di frontendnya
		file, err := c.FormFile("file")

		// kalo ada error, return si error
		if err != nil {
			return err
		}

		// dari source, open konten dari formnya, yang dalam hal ini file
		// jadi open filenya, terus kalo ada error, return errorrnya juga
		src, err := file.Open()

		// kalo ada errorrnya, return errorrnya
		if err != nil {
			return err
		}

		// kalo udah kelar semua close si src
		defer src.Close()

		// bikin filepath, yang bakal dipake buat nyimpen filenya
		filePath := "./public/uploads/" + file.Filename

		// bikin file baru, yang bakal dipake buat akses filenya via url
		fileSrc := "http://127.0.0.1:9000/uploads/" + file.Filename

		// bikin filenya itu sendiri, di path yang udah dibuat
		dst, err := os.Create(filePath)

		// kalo ada error, panic
		if err != nil {
			panic(err)
		}

		// terus close dstnya kalo misal udah
		defer dst.Close()

		// copy semua file dari src, ke dst, sampe abis, ato sampe error
		// terus direturn number of bytes yang tercopy, dan errorrnya kalo ada
		// kalo gaada error, error = nil, kalo error = EOF or anything else
		// pake _ namanya, karena gapeduli si number of bytesnya, jadi gabakal dipake lagi
		if _, err = io.Copy(dst, src); err != nil {
			panic(err)
		}

		// bikin query baru buat masukin ke db
		stmt, err := db.Prepare("INSERT INTO photos (src) VALUES(?)")

		// kalo error ngepanic
		if err != nil {
			panic(err)
		}

		// close stmt kalo udah semua
		defer stmt.Close()

		// execute querynya, dimana yang dimasukin itu filesrcnya yang udah ada
		result, err := stmt.Exec(fileSrc)

		// kalo error ngepanic
		if err != nil {
			panic(err)
		}

		// ambil idnya dia, pake lastinsertid()
		insertedId, err := result.LastInsertId()

		// kalo error ngepanic
		if err != nil {
			panic(err)
		}

		// bikin photo, yang isinya dari filesrc, sama idnya
		photo := Photo{
			Src: fileSrc,
			ID:  insertedId,
		}
		// ini cara pake client pushernya
		// jadi triggernya, kalo ada photo baru yang keupload
		// upload via photo-stream, via new-photo, isinya photo
		client.Trigger("photo-stream", "new-photo", photo)

		// kalo berhasil dia return status OK, sama photonya
		return c.JSON(http.StatusOK, photo)
	}
}

type PhotoCollection struct {
	// array of Photo, yang bentuknya itu json, yang disebutnya items
	Photos []Photo `json:"items"`
}

func initializeDatabase(filepath string) *sql.DB {
	// instantiate sqlitenya
	db, err := sql.Open("sqlite3", filepath)

	// kalo error, dia ngepanic
	if err != nil || db == nil {
		panic("Error connecting to database")
	}

	// terus ngereturn si sqlite yang bisa dipake
	return db
}

func migrateDatabase(db *sql.DB) {
	// ini terima sqlite buat bikin dbnya

	// ini query sqlite buat DDL (create, insert dll)
	sql := `
		CREATE TABLE IF NOT EXISTS photos(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			src VARCHAR NOT NULL
		);
	`

	// pake method Exec, yang udah ada, kita run query diatas
	_, err := db.Exec(sql)

	// tapi kalo error, tampilin errornya
	if err != nil {
		panic(err)
	}

	// kalo ga error, udah ke migrate, gaada yang direturn tapi
}

func main() {
	//ngebikin sqlite yang baru, terus taro di path yang kita minta
	// abis udah jadi sqlitenya, dipake di migrateDatabase dibawah
	db := initializeDatabase("database/database.sqlite")

	// ngebikin migration baru, pake sqlite yang udah dibikin tadi
	migrateDatabase(db)

	// instantiate Echo
	e := echo.New()

	// log info, buat HTTP requests
	e.Use(middleware.Logger())

	// Error handler, yang ke trigger kalo ada panic
	e.Use(middleware.Recover())

	// route /, dari public/index.html
	e.File("/", "public/index.html")

	// ini bikin endpoint
	e.GET("/photos", getPhotos(db))
	e.POST("/photos", uploadPhoto(db))

	// sama kayak get, tapi buat return file2 static yang gabakal berubah2
	e.Static("/uploads", "public/uploads")

	// ngestart server
	e.Logger.Fatal(e.Start(":9000"))
}
