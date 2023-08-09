package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"technical/dto"
	"technical/saveimage"

	_ "github.com/go-sql-driver/mysql"
)

func CreateEmploye(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	imageFile, handler, err := r.FormFile("faceid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer imageFile.Close()

	imageName := handler.Filename

	err = saveimage.SaveImage(imageFile, imageName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/technicaltest")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	createTableSQL := `create table if not exists employee (
		id int auto_increment primary key,
		name varchar(255),
		email varchar(255),
		password varchar(255),
		faceid BLOB
		)`

	_, err = db.Exec(createTableSQL)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("succes create table")

	insertValueSQL := "INSERT INTO employee (name, email,password, faceid) VALUES (?, ?, ?, ?)"

	_, err = db.Exec(insertValueSQL, name, email, password, imageName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Data berhasil di masukan")

}

func CreateLoginHistory(email string) error {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/technicaltest")
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `create table if not exists login_history (
		id int auto_increment primary key,
		email varchar(255),
		login_time timestamp
	)`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	insertValueSQL := "INSERT INTO login_history (email, login_time) VALUES (?, NOW())"
	_, err = db.Exec(insertValueSQL, email)
	if err != nil {
		return err
	}

	return nil
}

func CreateLogoutHistory(email string) error {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/technicaltest")
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `create table if not exists logout_history (
		id int auto_increment primary key,
		email varchar(255),
		logout_time timestamp
	)`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	insertValueSQL := "INSERT INTO logout_history (email, logout_time) VALUES (?, NOW())"
	_, err = db.Exec(insertValueSQL, email)
	if err != nil {
		return err
	}

	return nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	imageFile, _, err := r.FormFile("faceid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	isValid := IsLoginValid(email, password, imageFile)
	if isValid {
		err = CreateLoginHistory(email)
		if err != nil {
			http.Error(w, "Login berhasil, tetapi gagal membuat log", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Login berhasil")
	} else {
		http.Error(w, "Login gagal", http.StatusUnauthorized)
	}

}

func IsLoginValid(email, password string, imageFile multipart.File) bool {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/technicaltest")
	if err != nil {
		fmt.Println("Gagal membuka koneksi database:", err)
		return false
	}
	defer db.Close()

	var storedEmployee dto.User
	err = db.QueryRow("SELECT * FROM employee WHERE email = ?", email).Scan(
		&storedEmployee.Id,
		&storedEmployee.Name,
		&storedEmployee.Email,
		&storedEmployee.Password,
		&storedEmployee.Faceid,
	)
	if err != nil {
		fmt.Println("Gagal mengambil data employee:", err)
		return false
	}

	if password != storedEmployee.Password {
		return false
	}

	storedImageBytes, err := ioutil.ReadFile(`C:\Users\LENOVO\Desktop\technical\image\` + storedEmployee.Faceid)
	if err != nil {
		fmt.Println("Gagal membaca data gambar tersimpan:", err)
		return false
	}

	uploadedImageBytes, err := ioutil.ReadAll(imageFile)
	if err != nil {
		fmt.Println("Gagal membaca data gambar yang diunggah:", err)
		return false
	}

	if bytes.Compare(storedImageBytes, uploadedImageBytes) != 0 {
		return false
	}

	return true
}

func IsLogoutValid(email, password string, imageFile multipart.File) bool {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/technicaltest")
	if err != nil {
		fmt.Println("Gagal membuka koneksi database:", err)
		return false
	}
	defer db.Close()

	var storedEmployee dto.User
	err = db.QueryRow("SELECT * FROM employee WHERE email = ?", email).Scan(
		&storedEmployee.Id,
		&storedEmployee.Name,
		&storedEmployee.Email,
		&storedEmployee.Password,
		&storedEmployee.Faceid,
	)
	if err != nil {
		fmt.Println("Gagal mengambil data employee:", err)
		return false
	}

	if password != storedEmployee.Password {
		return false
	}

	storedImageBytes, err := ioutil.ReadFile(`C:\Users\LENOVO\Desktop\technical\image\` + storedEmployee.Faceid)
	if err != nil {
		fmt.Println("Gagal membaca data gambar tersimpan:", err)
		return false
	}

	uploadedImageBytes, err := ioutil.ReadAll(imageFile)
	if err != nil {
		fmt.Println("Gagal membaca data gambar yang diunggah:", err)
		return false
	}

	if bytes.Compare(storedImageBytes, uploadedImageBytes) != 0 {
		return false
	}

	return true
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	imageFile, _, err := r.FormFile("faceid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	isValid := IsLogoutValid(email, password, imageFile)
	if isValid {
		err = CreateLogoutHistory(email)
		if err != nil {
			http.Error(w, "Logout berhasil, tetapi gagal membuat log", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Logout berhasil")
	} else {
		http.Error(w, "Logout gagal", http.StatusUnauthorized)
	}

}

func main() {

	http.HandleFunc("/signup", CreateEmploye)
	http.HandleFunc("/signin", Login)
	http.HandleFunc("/signout", Logout)

	fmt.Println("running on 8080")

	http.ListenAndServe(":8080", nil)

}
