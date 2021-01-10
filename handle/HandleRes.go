package handle

import (
	"Go-Simple-Auth/db"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//PostRegister digunakan untuk mendapatkan nilai dari register
func PostRegister(w http.ResponseWriter, r *http.Request) {
	database, _ := db.Connect()
	defer database.Close()
	if r.Method != "POST" {
		http.ServeFile(w, r, "asset/register.html")
		return
	}
	//masukkan nilai dari template
	Username := r.FormValue("email")
	FName := r.FormValue("first_name")
	LName := r.FormValue("last_name")
	Password := r.FormValue("password")
	//cari Username kedalam database
	users := db.QueryUser(Username)
	//lakukan pengecekan terhadap users
	if (db.User{}) == users {
		//buat hash passwordnya
		HPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		//cek apakah ada error
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(HPassword) != 0 {
			//siapkan statement untuk nantinya dimasukkan ke db
			stmt, err := database.Prepare("INSERT INTO `users`(`ID`, `Username`, `First_Name`, `Last_Name`, `Password`) VALUES (?,?,?,?,?)")
			if err == nil {
				//eksekusi statment dengan nilai Username, HPassword, FName dan LName(ingat harus dengan pointer)
				_, err := stmt.Exec(nil, &Username, &FName, &LName, &HPassword)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//jika berhasil redirect ke login
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		}

	} else {
		//jika gagal balikkan ke register
		http.Redirect(w, r, "/register", 302)
	}
}
