package handle

import (
	"Go-Simple-Auth/db"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kataras/go-sessions/v3"
	"golang.org/x/crypto/bcrypt"
)

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		fmt.Println(r.Host + r.URL.Path)

		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}

	return true
}

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
		if len(HPassword) != 0 && checkErr(w, r, err) {
			//siapkan statement untuk nantinya dimasukkan ke db
			stmt, err := database.Prepare("INSERT INTO users SET username=?, password=?, first_name=?, last_name=?")
			if err == nil {
				//eksekusi statment dengan nilai Username, HPassword, FName dan LName(ingat harus dengan pointer)
				_, err := stmt.Exec(&Username, &HPassword, &FName, &LName)
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

//PostLogin adalah fungsi untuk Login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	//buat variabel untuk menampung sessions
	session := sessions.Start(w, r)
	//baut kondisi dimana panjang dari username tidak kosong
	if len(session.GetString("username")) != 0 {
		http.Redirect(w, r, "/", 302)
	}
	if r.Method != "POST" {
		http.ServeFile(w, r, "asset/login.html")
		return
	}

	//masukkan nilai dari Login
	Username := r.FormValue("username")
	Password := r.FormValue("password")
	//panggil fungsu query dan tampung nilainya
	//users := db.QueryUser(Username)
	//storage := &db.User{}
	selectData := db.QueryUser(Username)
	//dekrib passwordnya dan bandingkan
	//var PwdTamp = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(Password))
	password_test := bcrypt.CompareHashAndPassword([]byte(selectData.Password), []byte(Password))
	//cek apakah pwd kosong
	if password_test == nil {
		//jika kosong maka login berhasil
		//beri sessionnya
		session := sessions.Start(w, r)
		//ingat jangan pointer
		session.Set("username", selectData.Username)
		session.Set("name", selectData.Firsname)
		http.Redirect(w, r, "/", 302)
	} else {
		//gagal login
		http.Redirect(w, r, "/login", 302)
		fmt.Println(password_test.Error())
	}

}

//PostHome berfungsi untuk menampilkan user pada home
func PostHome(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/login", 301)
	}

	var data = map[string]string{
		"username": session.GetString("username"),
		"message":  "Welcome to the Go !",
	}
	var t, err = template.ParseFiles("asset/home.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, data)
	return

}
func LogOut(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}
